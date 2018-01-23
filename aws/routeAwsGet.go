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

package aws

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

// AwsAccountsFromUserIDByAccountID retrieves rows from 'trackit.aws_account' as AwsAccount.
// The result is filtered by a slice of accountID
func AwsAccountsFromUserIDByAccountID(db models.XODB, userID int, accountIDs []int) ([]AwsAccount, error) {
	var err error
	var stringAccountIDs []string

	// gen account_id
	for _, id := range accountIDs {
		stringAccountIDs = append(stringAccountIDs, strconv.FormatInt(int64(id), 10))
	}
	accountID := "(" + strings.Join(stringAccountIDs, ",") + ")"

	// sql query
	var sqlstr = `SELECT ` +
		`id, user_id, pretty, role_arn, external ` +
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
	res := []AwsAccount{}
	for q.Next() {
		aa := AwsAccount{}

		// scan
		err = q.Scan(&aa.Id, &aa.UserId, &aa.Pretty, &aa.RoleArn, &aa.External)
		if err != nil {
			return nil, err
		}

		res = append(res, aa)
	}

	return res, nil
}

// getAwsAccount is a route handler which returns the caller's list of
// AwsAccounts.
func getAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	var err error
	var awsAccounts []AwsAccount
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	if accountIds, ok := a[AwsAccountsOptionalQueryArg]; ok {
		awsAccounts, err = AwsAccountsFromUserIDByAccountID(tx, u.Id, accountIds.([]int))
	} else {
		awsAccounts, err = GetAwsAccountsFromUser(u, tx)
	}
	if err == nil {
		return 200, awsAccounts
	} else {
		l.Error("failed to get user's AWS accounts", err.Error())
		return 500, errors.New("failed to retrieve AWS accounts")
	}
}
