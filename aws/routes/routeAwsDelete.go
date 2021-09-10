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

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// DeleteAwsAccountFromAccountID delete an AWS account based on the
// accountID passed to it. It does not perform any check, especially
// on authorizations, which needs to be done by the caller
func DeleteAwsAccountFromAccountID(db models.DB, userID int, accountID int) (int, error) {
	var sqlstr = `DELETE FROM trackit.aws_account WHERE id = ? and user_id = ?`

	models.Logf(sqlstr, accountID, userID)
	buff, err := db.Exec(sqlstr, accountID, userID)
	if err != nil {
		return -1, err
	}
	res, err := buff.RowsAffected()
	return int(res), err
}

// deleteAwsAccount is a route handler which delete the
// AWS account passed in query args.
func deleteAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	aa := a[aws.AwsAccountSelection].(aws.AwsAccount)
	accountToDeleteID := aa.Id
	dbAwsBillRepositories, err := models.AwsBillRepositoryByAwsAccountID(tx, aa.Id)
	if err != nil {
		l.Error("unable to retrieve bill repositories", err.Error())
	}
	res, err := DeleteAwsAccountFromAccountID(tx, u.Id, accountToDeleteID)
	if err != nil {
		l.Error("error while deleting user's AWS account", err.Error())
		if res == -1 {
			l.Error("failed to retrieve user's AWS accounts", err.Error())
			return http.StatusInternalServerError, errors.New("failed to retrieve user's AWS accounts")
		} else {
			l.Error("specified AWS account is not in user's accounts", err.Error())
			return http.StatusInternalServerError, errors.New("specified AWS account is not in user's accounts")
		}
	}
	go func() {
		for _, br := range dbAwsBillRepositories {
			err = es.CleanByBillRepositoryId(context.Background(), aa.UserId, br.ID)
			if err != nil {
				l.Error("Failed to clean ES data for bill repository", err.Error())
			}
		}
	}()
	return http.StatusOK, nil
}
