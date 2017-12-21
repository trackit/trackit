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

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

// patchAwsAccountRequestBody is all the possible bodies for the
// patchAwsAccount request handler.
type patchAwsAccountRequestBody struct {
	Pretty string `json:"pretty"`
}

var (
	errFailUpdateAccount = errors.New("Failed to update AWS account.")
)

// patchAwsAccount is a route handler which lets the user update AwsAccounts from
// their account.
func patchAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body patchAwsAccountRequestBody
	err := decodeRequestBody(r, &body)
	if err == nil {
		tx := a[db.Transaction].(*sql.Tx)
		u := a[users.AuthenticatedUser].(users.User)
		id := a[AwsAccountQueryArg].(int)
		return patchAwsAccountWithValidBody(r, tx, u, body, int(id))
	} else {
		return 400, errors.New("Body is invalid.")
	}
}

// patchAwsAccountWithValidBody handles the logic of patchAwsAccount assuming the
// request body is valid.
func patchAwsAccountWithValidBody(r *http.Request, tx *sql.Tx, user users.User, body patchAwsAccountRequestBody, id int) (int, interface{}) {
	ctx := r.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	awsAccount, err := GetAwsAccountWithIdFromUser(user, id, tx)
	if err == nil {
		awsAccount.Pretty = body.Pretty
		if err := awsAccount.UpdatePrettyAwsAccount(ctx, tx); err != nil {
			logger.Error("Failed to update AWS Account.", err)
			return 500, errFailUpdateAccount
		}
	} else {
		logger.Error("Failed to get user's AWS accounts.", err.Error())
		return 500, errors.New("Failed to retrieve AWS accounts.")
	}
	return 200, awsAccount
}
