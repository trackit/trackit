//   Copyright 2018 MSolution.IO
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

package users

import (
	"database/sql"
	"net/http"
	"errors"
	"context"

	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"

	"github.com/trackit/jsonlog"
	"go/ast"
)

// inviteUserRequest is the expected request body for the invite user route handler.
type inviteUserRequest struct {
	Email              string `json:"email"    req:"nonzero"`
	AccountId          int `json:"accountId"   req:"nonzero"`
	PermissionLevel    int `json:"level"       req:"nonzero"`
}

func init() {
	routes.MethodMuxer{
		http.MethodPost: routes.H(inviteUser).With(
			routes.RequestContentType{"application/json"},
			routes.RequestBody{inviteUserRequest{"example@example.com", 1234, 0}},
			db.RequestTransaction{db.Db},
			routes.Documentation{
				Summary:     "Creates an invite",
				Description: "Creates an invite for account team sharing",
			},
		),
	}.H().Register("/user/invite")
}

// inviteUser handles users invite for team sharing.
func inviteUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body inviteUserRequest
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	return inviteUserWithValidBody(request, body, tx)
}

// checkuserWithEmail checks if user already exist
func CheckUserWithEmail(ctx context.Context, db models.XODB, userEmail string) (res bool, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByEmail(db, userEmail)
	_ = dbUser
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		logger.Error("Error getting user from database.", err.Error())
		return false, err
	} else {
		return true, nil
	}
}

// addAccountToGuest adds an entry in shared_account table with element that enable a user
// to share an access to all or part of his account
func addAccountToGuest(ctx context.Context, db models.XODB, userMail string, accountId int, permissionLevel int) (err error) {
	
}

// logInWithValidBody tries to authenticate and log a user in using a
// validated login request.
func inviteUserWithValidBody(request *http.Request, body inviteUserRequest, tx *sql.Tx) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	result, err := CheckUserWithEmail(request.Context(), tx, body.Email)
	if err == nil {
		if result {
			//Add account in database shared_accounts
			err = addAccountToGuest(request.Context(), tx, body.Email, body.AccountId, body.PermissionLevel)
			if err == nil {
				//Send mail notification to user
				sendMailNotification(body.Email, body.PermissionLevel)
			} else {
				logger.Warning("Error occured while adding account to an existing user.", err)
				return 403, errors.New("An error occured while inviting a new user. Please, try again.")
			}
		} else {
			//Create an account with temp password
			//Set value that set the account as "temp password so the user is guided to change it after first connection
			tempPassword, err := createAccountForGuest(request.Context(), tx, body.Email, body.AccountId, body.PermissionLevel)
			if err == nil {
				//Send mail notification to user
				sendMailNotification(body.Email, body.PermissionLevel, tempPassword)
			} else {
				logger.Warning("Error occured while creating new account for a guest.", err)
				return 403, errors.New("An error occured while inviting a new user. Please, try again.")
			}
		}
	} else {
		logger.Warning("Error occured while checking if user already exist.", err)
		return 403, errors.New("An error occured while inviting a new user. Please, try again.")
	}
}
