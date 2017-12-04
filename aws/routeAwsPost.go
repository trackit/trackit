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
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

// postAwsAccountRequestBody is the expected request body for the
// postAwsAccount request handler.
type postAwsAccountRequestBody struct {
	RoleArn  string `json:"roleArn"  req:"nonzero"`
	External string `json:"external" req:"nonzero"`
	Pretty   string `json:"pretty"`
}

var (
	errInvalidAccount     = errors.New("Could not validate role and external ID.")
	errFailCreateAccount  = errors.New("Failed to create AWS account.")
	errFailUpdateExternal = errors.New("Failed to update external.")
)

// postAwsAccount is a route handler which lets the user add AwsAccounts to
// their account.
func postAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body postAwsAccountRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	u := a[users.AuthenticatedUser].(users.User)
	return postAwsAccountWithValidBody(r, tx, u, body)
}

// postAwsAccountWithValidBody handles the logic of postAwsAccount assuming the
// request body is valid.
func postAwsAccountWithValidBody(r *http.Request, tx *sql.Tx, user users.User, body postAwsAccountRequestBody) (int, interface{}) {
	ctx := r.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	account := AwsAccount{
		RoleArn:  body.RoleArn,
		External: body.External,
		UserId:   user.Id,
		Pretty:   body.Pretty,
	}
	if account.External != user.NextExternal {
		logger.Warning("Tried to add AWS account with bad external.", account)
		return 400, errors.New("Incorrect external. Use /aws/next to get expected external.")
	} else if err := testAndCreateAwsAccount(ctx, tx, &account, &user); err == nil {
		return 200, account
	} else {
		switch err {
		case errInvalidAccount:
			return 400, err
		default:
			return 500, err
		}
	}
}

// testAndCreateAwsAccount tests an AwsAccount can be assumed-role and then
// saves it to the database.
func testAndCreateAwsAccount(ctx context.Context, tx *sql.Tx, account *AwsAccount, user *users.User) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	if _, err := GetTemporaryCredentials(*account, "validityTest"); err != nil {
		return errInvalidAccount
	}
	if err := account.CreateAwsAccount(ctx, tx); err != nil {
		logger.Error("Failed to insert AWS account.", newTestAndCreateAwsAccountError(err, *account, *user))
		return errFailCreateAccount
	}
	user.NextExternal = ""
	if err := user.UpdateNextExternal(ctx, tx); err != nil {
		logger.Error("Failed to update external.", newTestAndCreateAwsAccountError(err, *account, *user))
		return errFailUpdateExternal
	}
	return nil
}

// testAndCreateAwsAccountError is used to log errors in
// testAndCreateAwsAccount.
type testAndCreateAwsAccountError struct {
	err     string     `json:"error"`
	account AwsAccount `json:"account"`
	user    users.User `json:"user"`
}

// newTestAndCreateAwsAccountError is used to log errors in
// testAndCreateAwsAccount.
func newTestAndCreateAwsAccountError(e error, a AwsAccount, u users.User) testAndCreateAwsAccountError {
	return testAndCreateAwsAccountError{
		err:     e.Error(),
		account: a,
		user:    u,
	}
}
