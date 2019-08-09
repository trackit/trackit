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
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// patchAwsAccountRequestBody is all the possible bodies for the
// patchAwsAccount request handler.
type patchAwsAccountRequestBody struct {
	Pretty string `json:"pretty"`
	Payer  bool   `json:"payer"`
}

var (
	errFailUpdateAccount = errors.New("failed to update AWS account")
)

// patchAwsAccount is a route handler which lets the user update AwsAccounts from
// their account.
func patchAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body patchAwsAccountRequestBody
	err := decodeRequestBody(r, &body)
	if err == nil {
		tx := a[db.Transaction].(*sql.Tx)
		u := a[users.AuthenticatedUser].(users.User)
		id := a[routes.AwsAccountIdQueryArg].(int)
		return patchAwsAccountWithValidBody(r, tx, u, body, int(id))
	} else {
		return 400, errors.New("body is invalid")
	}
}

// patchAwsAccountWithValidBody handles the logic of patchAwsAccount assuming the
// request body is valid.
func patchAwsAccountWithValidBody(r *http.Request, tx *sql.Tx, user users.User, body patchAwsAccountRequestBody, id int) (int, interface{}) {
	ctx := r.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	awsAccount, err := aws.GetAwsAccountWithIdFromUser(user, id, tx)
	if err == nil {
		awsAccount.Pretty = body.Pretty
		awsAccount.Payer = body.Payer
		if err := awsAccount.UpdatePrettyAwsAccount(ctx, tx); err != nil {
			logger.Error("failed to update AWS Account", err)
			return 500, errFailUpdateAccount
		}
	} else {
		logger.Error("failed to get user's AWS accounts", err.Error())
		return 500, errors.New("failed to retrieve AWS accounts")
	}
	return 200, awsAccount
}
