//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package routes

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/s3"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getBillRepositoryUpdates).With(
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get user's bill repositories and info about their update status",
				Description: "Gets the list of the user's bill repositories and info about when they have updated or will update.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
	).Register("/aws/billrepositoryupdates")
}

type dbAccessor interface {
	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

type BillRepositoryUpdateInfo struct {
	BillRepositoryId int        `json:"billRepositoryId"`
	AwsAccountPretty string     `json:"awsAccountPretty"`
	AwsAccountId     int        `json:"awsAccountId"`
	Bucket           string     `json:"bucket"`
	Prefix           string     `json:"prefix"`
	NextStarted      *time.Time `json:"nextStarted"`
	NextPending      *bool      `json:"nextPending"`
	LastStarted      *time.Time `json:"lastStarted"`
	LastFinished     *time.Time `json:"lastFinished"`
	LastError        *string    `json:"lastError"`
}

func getBillRepositoryUpdates(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	logger := jsonlog.LoggerFromContextOrDefault(r.Context())
	updateInfo, err := BillRepositoryUpdates(tx, u.Id)
	if err != nil {
		logger.Error("Failed to get bill repository update jobs.", map[string]interface{}{
			"user":  u,
			"error": err.Error(),
		})
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, updateInfo
}

func BillRepositoryUpdates(db dbAccessor, userId int) ([]BillRepositoryUpdateInfo, error) {
	// for each seleced bill repository, find data about the last and
	// next/current update and join it all
	var sqlstr = `
		SELECT
		  aws_bill_repository.id             AS id,
		  aws_account.pretty                 AS aws_account_pretty,
		  aws_bill_repository.aws_account_id AS aws_account_id,
		  aws_bill_repository.bucket         AS bucket,
		  aws_bill_repository.prefix         AS prefix,
		  aws_bill_repository.next_update    AS next_update,
		  (last_pending.id IS NOT NULL)      AS next_pending,
		  last_completed.created             AS last_started,
		  last_completed.completed           AS last_finished,
		  last_completed.error               AS last_error
		FROM aws_bill_repository
		INNER JOIN aws_account ON
		  aws_bill_repository.aws_account_id = aws_account.id
		LEFT OUTER JOIN (
		  SELECT * FROM (
		    SELECT
		      *,
		      ROW_NUMBER() OVER(PARTITION BY aws_bill_repository_id
		                        ORDER BY     created DESC) AS rn
		    FROM aws_bill_update_job
		    WHERE completed > 0
		  ) AS completed
		  WHERE completed.rn = 1
		) AS last_completed ON
		  aws_bill_repository.id = last_completed.aws_bill_repository_id
		LEFT OUTER JOIN (
		  SELECT * FROM (
		    SELECT
		      *,
		      ROW_NUMBER() OVER(PARTITION BY aws_bill_repository_id
		                        ORDER BY     completed DESC) AS rn
		    FROM aws_bill_update_job
		    WHERE completed = 0 AND expired >= NOW()
		  ) AS pending
		  WHERE pending.rn = 1
		) AS last_pending ON
		  aws_bill_repository.id = last_pending.aws_bill_repository_id
		WHERE aws_account.user_id = ?
	`
	q, err := db.Query(sqlstr, userId)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var res []BillRepositoryUpdateInfo
	var i int
	for i = 0; q.Next(); i++ {
		res = append(res, BillRepositoryUpdateInfo{})
		err = q.Scan(
			&res[i].BillRepositoryId,
			&res[i].AwsAccountPretty,
			&res[i].AwsAccountId,
			&res[i].Bucket,
			&res[i].Prefix,
			&res[i].NextStarted,
			&res[i].NextPending,
			&res[i].LastStarted,
			&res[i].LastFinished,
			&res[i].LastError,
		)
		if err != nil {
			return nil, err
		}
	}
	return res[:i], nil
}

func billRepositoryWithPendingFromAwsAccount(db dbAccessor, awsAccountId int) ([]s3.BillRepositoryWithPending, error) {
	var sqlstr = `
		SELECT
		  aws_bill_repository.id                     AS id,
		  aws_bill_repository.aws_account_id         AS aws_account_id,
		  aws_bill_repository.bucket                 AS bucket,
		  aws_bill_repository.prefix                 AS prefix,
		  aws_bill_repository.error                  AS error,
		  aws_bill_repository.last_imported_manifest AS last_imported_manifest,
		  aws_bill_repository.next_update            AS next_update,
		  (last_pending.id IS NOT NULL)              AS next_pending
		FROM aws_bill_repository
		LEFT OUTER JOIN (
		  SELECT *
		  FROM aws_bill_update_job
		  WHERE completed = 0 AND expired >= NOW()
		  LIMIT 1
		) AS last_pending ON
		  aws_bill_repository.id = last_pending.aws_bill_repository_id
		WHERE aws_bill_repository.aws_account_id = ?
	`
	q, err := db.Query(sqlstr, awsAccountId)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var res []s3.BillRepositoryWithPending
	var i int
	for i = 0; q.Next(); i++ {
		res = append(res, s3.BillRepositoryWithPending{})
		err = q.Scan(
			&res[i].Id,
			&res[i].AwsAccountId,
			&res[i].Bucket,
			&res[i].Prefix,
			&res[i].Error,
			&res[i].LastImportedManifest,
			&res[i].NextUpdate,
			&res[i].NextPending,
		)
		if err != nil {
			return nil, err
		}
	}
	return res[:i], nil
}

func intArrayToStringArray(integers []int) (strings []string) {
	strings = make([]string, len(integers))
	for i := range integers {
		strings[i] = strconv.FormatInt(int64(integers[i]), 10)
	}
	return
}

func intArrayToSqlSet(integers []int) string {
	ss := intArrayToStringArray(integers)
	return "(" + strings.Join(ss, ",") + ")"
}

// AwsAccountsFromUserIDByAccountID retrieves rows from 'trackit.aws_account' as AwsAccount.
// The result is filtered by a slice of accountID
func AwsAccountsFromUserIDByAccountID(db models.XODB, userID int, accountIDs []int) ([]aws.AwsAccount, error) {
	var err error
	var stringAccountIDs []string

	// gen account_id
	stringAccountIDs = intArrayToStringArray(accountIDs)
	accountID := "(" + strings.Join(stringAccountIDs, ",") + ")"

	// sql query
	var sqlstr = `SELECT ` +
		`id, user_id, pretty, role_arn, external, aws_identity ` +
		`FROM trackit.aws_account ` +
		`WHERE user_id = ? ` +
		`AND id IN ` + accountID

	// run query
	models.XOLog(sqlstr, userID)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	// load results
	res := []aws.AwsAccount{}
	for q.Next() {
		aa := aws.AwsAccount{}

		// scan
		err = q.Scan(&aa.Id, &aa.UserId, &aa.Pretty, &aa.RoleArn, &aa.External, &aa.AwsIdentity)
		aa.AccountOwner = true
		if err != nil {
			return nil, err
		}

		res = append(res, aa)
	}

	return res, nil
}

type AwsAccountWithBillRepositories struct {
	aws.AwsAccount
	BillRepositories []s3.BillRepositoryWithPending   `json:"billRepositories"`
	SubAccounts      []AwsAccountWithBillRepositories `json:"subAccounts,omitempty"`
}

// getAwsAccount is a route handler which returns the caller's list of
// AwsAccounts.
func getAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	var awsErr error
	var awsAccounts []aws.AwsAccount
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	if accountIds, ok := a[routes.AwsAccountIdsOptionalQueryArg]; ok {
		awsAccounts, awsErr = AwsAccountsFromUserIDByAccountID(tx, u.Id, accountIds.([]int))
	} else {
		awsAccounts, awsErr = aws.GetAwsAccountsFromUser(u, tx)
	}
	if awsErr == nil {
		return 200, awsAccounts
	} else {
		l.Error("failed to get user's AWS accounts", awsErr.Error())
		return 500, errors.New("failed to retrieve AWS accounts")
	}
}

func sortSubAccounts(awsAccountsWithBillRepositories []AwsAccountWithBillRepositories) ([]AwsAccountWithBillRepositories, error) {
	accounts := make([]AwsAccountWithBillRepositories, 0)
	for _, aa := range awsAccountsWithBillRepositories {
		if aa.ParentId.Valid == false {
			accounts = append(accounts, aa)
		}
	}
	for _, aa := range awsAccountsWithBillRepositories {
		if aa.ParentId.Valid == false {
			continue
		}
		foundMatch := false
	AccountsLoop:
		for i, account := range accounts {
			if aa.ParentId.Int64 == int64(account.Id) {
				if accounts[i].SubAccounts == nil {
					accounts[i].SubAccounts = make([]AwsAccountWithBillRepositories, 0)
				}
				accounts[i].SubAccounts = append(accounts[i].SubAccounts, aa)
				foundMatch = true
				continue AccountsLoop
			}
		}
		if foundMatch == false {
			accounts = append(accounts, aa)
		}
	}
	return accounts, nil
}

func buildAwsAccountsWithBillRepositoriesFromAwsAccounts(awsAccounts []aws.AwsAccount, tx *sql.Tx) (awsAccountsWithBillRepositories []AwsAccountWithBillRepositories, err error) {
	timeLimit := time.Now().AddDate(0, 0, -7)
	for _, aa := range awsAccounts {
		aawbr := AwsAccountWithBillRepositories{
			aa,
			[]s3.BillRepositoryWithPending{},
			nil,
		}
		if aawbr.BillRepositories, err = billRepositoryWithPendingFromAwsAccount(tx, aa.Id); err != nil {
			return
		}
		for id, br := range aawbr.BillRepositories {
			if br.LastImportedManifest.Before(timeLimit) {
				aawbr.BillRepositories[id].NextPending = true
			}
		}
		awsAccountsWithBillRepositories = append(awsAccountsWithBillRepositories, aawbr)
	}
	return sortSubAccounts(awsAccountsWithBillRepositories)
}
