//   Copyright 2019 MSolution.IO
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
	"strings"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// patchAwsSubaccountRequestBody is the body expected by patchAwsSubaccount
type patchAwsSubaccountRequestBody struct {
	RoleArn  string `json:"roleArn"  req:"nonzero"`
	External string `json:"external" req:"nonzero"`
}

// patchAwsSubaccount is a route handler which lets the user update a subaccount
func patchAwsSubaccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body patchAwsSubaccountRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	u := a[users.AuthenticatedUser].(users.User)
	id := a[routes.AwsAccountIdQueryArg].(int)
	return patchAwsSubaccountWithValidBody(r, tx, u, body, id)
}

// patchAwsSubaccountWithValidBody handles the logic of patchAwsSubaccount assuming the
// request body is valid.
func patchAwsSubaccountWithValidBody(r *http.Request, tx *sql.Tx, user users.User, body patchAwsSubaccountRequestBody, id int) (int, interface{}) {
	ctx := r.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	awsAccount, err := aws.GetAwsAccountWithIdFromUser(user, id, tx)
	if err != nil {
		logger.Warning("failed to get user's AWS accounts", err.Error())
		return http.StatusBadRequest, errors.New("failed to retrieve AWS accounts.")
	}
	awsAccount.RoleArn = body.RoleArn
	awsAccount.External = body.External
	if !awsAccount.ParentId.Valid {
		logger.Info("tried to edit an AWS account as a sub-account", awsAccount)
		return http.StatusBadRequest, errors.New("not a sub-account.")
	}
	if awsAccount.External != user.NextExternal {
		logger.Info("tried to edit AWS account with bad external", awsAccount)
		return http.StatusBadRequest, errors.New("incorrect external.")
	}
	if testRoleIdentityMatch(awsAccount) == false {
		logger.Info("role account id does not match aws identity", awsAccount)
		return http.StatusBadRequest, errors.New("role account id does not match aws identity.")
	}
	if err := testAndUpdateSubaccount(ctx, tx, awsAccount, user); err != nil {
		switch err {
		case errInvalidAccount:
			return http.StatusBadRequest, err
		default:
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}

// testRoleIdentityMatch checks that the account id in the role matches the aws identity
func testRoleIdentityMatch(awsAccount aws.AwsAccount) bool {
	arnElems := strings.Split(awsAccount.RoleArn, ":")
	if len(arnElems) != 6 {
		return false
	}
	return arnElems[4] == awsAccount.AwsIdentity
}

// testAndUpdateSubaccount tests an AwsAccount can be assumed-role and then
// updates it in the database.
// It also sets the NextExternal to "".
func testAndUpdateSubaccount(ctx context.Context, tx *sql.Tx, awsAccount aws.AwsAccount, user users.User) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	if _, err := aws.GetTemporaryCredentials(awsAccount, "validityTest"); err != nil {
		logger.Info("failed to generate temporary credentials", newTestAndCreateAwsAccountError(err, awsAccount, user))
		return errInvalidAccount
	}
	if err := awsAccount.UpdateRoleAndExternalAwsAccount(ctx, tx); err != nil {
		logger.Error("failed to update AWS account", newTestAndCreateAwsAccountError(err, awsAccount, user))
		return errFailUpdateAccount
	}
	user.NextExternal = ""
	if err := user.UpdateNextExternal(ctx, tx); err != nil {
		logger.Error("failed to update external", newTestAndCreateAwsAccountError(err, awsAccount, user))
		return errFailUpdateExternal
	}
	return nil
}
