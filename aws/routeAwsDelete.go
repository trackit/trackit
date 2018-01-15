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
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

// DeleteAwsAccountFromAccountID delete an AWS account based on the
// accountID passed to it. It does not perform any check, especially
// on authorizations, which needs to be done by the caller
func DeleteAwsAccountFromAccountID(db models.XODB, accountID int) error {
	var sqlstr = `DELETE FROM trackit.aws_account WHERE id = ?`

	models.XOLog(sqlstr, accountID)
	buff, err := db.Exec(sqlstr, accountID)
	fmt.Print(buff)
	if err != nil {
		return err
	}
	return nil
}

// checkIfAccountInSlice will check for the presence of the accountID
// in the slice of AWS accounts passed to it
func checkIfAccountInSlice(aa int, saa []AwsAccount) bool {
	for _, account := range saa {
		if aa == account.Id {
			return true
		}
	}
	return false
}

// deleteAwsAccount is a route handler which delete the
// AWS account passed in query args.
func deleteAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	accountToDeleteID := a[AwsAccountQueryArg].(int)
	userAccounts, err := GetAwsAccountsFromUser(u, tx)
	if err != nil {
		l.Error("failed to retrieve user's AWS accounts", err.Error())
		return http.StatusInternalServerError, fmt.Sprintf("failed to retrieve user's aws accounts")
	}
	if !checkIfAccountInSlice(accountToDeleteID, userAccounts) {
		l.Info("aws account not in user's accounts", nil)
		return http.StatusForbidden, fmt.Sprintf("Specified AWS account is not in user's accounts")
	}
	err = DeleteAwsAccountFromAccountID(tx, accountToDeleteID)
	if err != nil {
		l.Error("error while deleting user's AWS account", err.Error())
		return http.StatusInternalServerError, fmt.Sprintf("error while deleting user's account")
	}
	return http.StatusOK, nil
}
