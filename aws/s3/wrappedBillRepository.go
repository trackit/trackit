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

package s3

import (
	"database/sql"
	"time"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/models"
)

// BillRepositoryWithPending is a BillRepository
// wrapped with NextPending
type BillRepositoryWithPending struct {
	BillRepository
	NextPending bool `json:"nextPending"`
}

// Status contains a value and a human readable detail.
type Status struct {
	Value  string `json:"value"`
	Detail string `json:"detail"`
}

// BillRepositoryWithStatus is a BillRepository
// wrapped with Status
type BillRepositoryWithStatus struct {
	BillRepositoryWithPending
	Status Status `json:"status"`
}

// AwsAccountWithBillRepositoriesWithPending represents a
// client's AWS account with its bill repositories.
type AwsAccountWithBillRepositoriesWithPending struct {
	aws.AwsAccount
	BillRepositories []BillRepositoryWithPending                 `json:"billRepositories"`
	SubAccounts      []AwsAccountWithBillRepositoriesWithPending `json:"subAccounts,omitempty"`
}

// AwsAccountWithBillRepositoriesWithStatus represents a
// client's AWS account with its bill repositories.
type AwsAccountWithBillRepositoriesWithStatus struct {
	aws.AwsAccount
	BillRepositories []BillRepositoryWithStatus                 `json:"billRepositories"`
	SubAccounts      []AwsAccountWithBillRepositoriesWithStatus `json:"subAccounts,omitempty"`
}

// GetBillRepositoryWithPendingForAwsAccount gets a BillRepositoryWithPending by Aws Account id
func GetBillRepositoryWithPendingForAwsAccount(tx *sql.Tx, awsAccountId int) ([]BillRepositoryWithPending, error) {
	timeLimit := time.Now().AddDate(0, 0, -7)
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
	q, err := tx.Query(sqlstr, awsAccountId)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var res []BillRepositoryWithPending
	var i int
	for i = 0; q.Next(); i++ {
		res = append(res, BillRepositoryWithPending{})
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
	for id, br := range res {
		if br.LastImportedManifest.Before(timeLimit) {
			res[id].NextPending = true
		}
	}
	return res[:i], nil
}

// WrapAwsAccountsWithBillRepositories wraps AwsAccounts with BillRepositories.
// Value returned contains all sub accounts.
func WrapAwsAccountsWithBillRepositories(awsAccounts []aws.AwsAccount, tx *sql.Tx) (awsAccountsWithBillRepositories []AwsAccountWithBillRepositoriesWithPending, err error) {
	for _, aa := range awsAccounts {
		aawbr := AwsAccountWithBillRepositoriesWithPending{
			aa,
			[]BillRepositoryWithPending{},
			nil,
		}
		if aawbr.BillRepositories, err = GetBillRepositoryWithPendingForAwsAccount(tx, aa.Id); err != nil {
			return
		}
		awsAccountsWithBillRepositories = append(awsAccountsWithBillRepositories, aawbr)
	}
	return sortSubAccounts(awsAccountsWithBillRepositories)
}

func sortSubAccounts(awsAccountsWithBillRepositories []AwsAccountWithBillRepositoriesWithPending) ([]AwsAccountWithBillRepositoriesWithPending, error) {
	accounts := make([]AwsAccountWithBillRepositoriesWithPending, 0)
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
					accounts[i].SubAccounts = make([]AwsAccountWithBillRepositoriesWithPending, 0)
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

// WrapBillRepositoriesWithPendingWithStatus wraps BillRepositoryWithPending with a status.
func WrapBillRepositoriesWithPendingWithStatus(tx *sql.Tx, billRepositoryWithPending BillRepositoryWithPending) (BillRepositoryWithStatus, error) {
	abuj, err := models.LastAwsBillUpdateJobsByAwsBillRepositoryID(tx, billRepositoryWithPending.Id)
	return BillRepositoryWithStatus{
		billRepositoryWithPending,
		getStatusMessage(billRepositoryWithPending, abuj),
	}, err
}

// WrapAwsAccountsWithBillRepositoriesWithPendingWithStatus wraps AwsAccountWithBillRepositoriesWithPending with a status.
func WrapAwsAccountsWithBillRepositoriesWithPendingWithStatus(awsAccountsWithBillRepositories []AwsAccountWithBillRepositoriesWithPending, tx *sql.Tx) []AwsAccountWithBillRepositoriesWithStatus {
	result := make([]AwsAccountWithBillRepositoriesWithStatus, 0)
	for _, awsAccount := range awsAccountsWithBillRepositories {
		var billRepositories []BillRepositoryWithStatus
		var subAccounts []AwsAccountWithBillRepositoriesWithStatus
		for _, billRepository := range awsAccount.BillRepositories {
			brws, _ := WrapBillRepositoriesWithPendingWithStatus(tx, billRepository)
			billRepositories = append(billRepositories, brws)
		}
		if awsAccount.SubAccounts != nil && len(awsAccount.SubAccounts) > 0 {
			subAccounts = WrapAwsAccountsWithBillRepositoriesWithPendingWithStatus(awsAccount.SubAccounts, tx)
		}
		account := AwsAccountWithBillRepositoriesWithStatus{
			awsAccount.AwsAccount,
			billRepositories,
			subAccounts,
		}
		result = append(result, account)
	}
	return result
}

func getStatusMessage(br BillRepositoryWithPending, item *models.AwsBillUpdateJob) Status {
	if item == nil {
		return Status{}
	} else if len(br.Error) > 0 {
		return Status{
			Value:  "error",
			Detail: br.Error,
		}
	} else if item == nil {
		return Status{
			Value:  "not_started",
			Detail: "",
		}
	} else if len(item.Error) > 0 {
		return Status{
			Value:  "error",
			Detail: item.Error,
		}
	} else if item.Completed.IsZero() {
		return Status{
			Value:  "in_progress",
			Detail: "",
		}
	} else {
		return Status{
			Value:  "ok",
			Detail: "",
		}
	}
}
