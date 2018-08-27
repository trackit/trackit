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
	"fmt"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/satori/go.uuid"

	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/mail"
)

// inviteUserRequest is the expected request body for the invite user route handler.
type inviteUserRequest struct {
	Email              string `json:"email" req:"nonzero"`
	AccountId          int `json:"accountId"`
	PermissionLevel    int `json:"permissionLevel"`
}

func init() {
	routes.MethodMuxer{
		http.MethodPost: routes.H(inviteUser).With(
			routes.RequestContentType{"application/json"},
			db.RequestTransaction{db.Db},
			RequireAuthenticatedUser{ViewerAsParent},
			routes.RequestBody{inviteUserRequest{"example@example.com", 1234, 0}},
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
	user := a[AuthenticatedUser].(User)
	return inviteUserWithValidBody(request, body, tx, user)
}

// checkuserWithEmail checks if user already exist.
// true is returned if invited user already exist.
func checkUserWithEmail(ctx context.Context, db models.XODB, userEmail string) (bool, int, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByEmail(db, userEmail)
	if err == sql.ErrNoRows {
		return false, 0 , nil
	} else if err != nil {
		logger.Error("Error getting user from database.", err.Error())
		return false, 0, err
	} else {
		return true, dbUser.ID,nil
	}
}

// checkSharedAccount checks if an account is already shared with a user.
// true is returned if invited user already have an access to this account.
func checkSharedAccount(ctx context.Context, db models.XODB, accountId int, userId int) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbSharedAccounts, err := models.SharedAccountsByUserID(db, userId)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		logger.Error("Error getting shared account from database.", err.Error())
		return false, err
	} else {
		for _, key := range dbSharedAccounts {
			if key.AccountID == accountId {
				return true, nil
			}
		}
	}
	return false,nil
}

// addAccountToGuest adds an entry in shared_account table allowing a user
// to share an access to all or part of his account
func addAccountToGuest(ctx context.Context, db *sql.Tx, accountId int, permissionLevel int, guestId int, ownerId int) (error) {
	const sqlstr = `INSERT INTO shared_account(
			account_id, owner_id, user_id, user_permission, account_status
		) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(sqlstr, accountId, ownerId, guestId, permissionLevel, 0)
	return err
}

// createAccountForGuest creates an account for invited user who do not already own an account
func createAccountForGuest(ctx context.Context, db *sql.Tx, userMail string, accountId int, permissionLevel int, user User) (int, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	tempPassword := uuid.NewV1().String()
	usr, err := CreateUserWithPassword(ctx, db, userMail, tempPassword, "")
	if err == nil {
		err = addAccountToGuest(ctx, db, accountId, permissionLevel, usr.Id, user.Id)
		if err != nil {
			logger.Error("Error occured while adding account to an newly created user.", err.Error())
			return 0, err
		}
	} else {
		logger.Error("Error occured while creating an automatic new account.", err.Error())
		return 0, err
	}
	return usr.Id,nil
}

// resetPasswordGenerator returns a reset password token. It is used in order to
// create an account and let the user choose his own password
func resetPasswordGenerator(ctx context.Context, tx *sql.Tx, newUserId int) (models.ForgottenPassword, string, error) {
	var dbForgottenPassword models.ForgottenPassword
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	token := uuid.NewV1().String()
	tokenHash, err := getPasswordHash(token)
	if err != nil {
		logger.Error("Failed to create token hash.", err.Error())
		return dbForgottenPassword, "", err
	}
	dbForgottenPassword = models.ForgottenPassword{
		UserID:  newUserId,
		Token:   tokenHash,
		Created: time.Now(),
	}
	err = dbForgottenPassword.Insert(tx)
	if err == nil {
		return dbForgottenPassword, token, err
	} else {
		logger.Error("Failed to insert forgotten password", err.Error())
		return dbForgottenPassword, "", err
	}
}

// sendMailNotification sends an email to user how has been invited to access a AWS account on trackit.io
func sendMailNotification(ctx context.Context, tx *sql.Tx, userMail string, userNew bool, newUserId int) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	if userNew {
		mailSubject := "An AWS account has been added to your Trackit account"
		mailBody := fmt.Sprintf("%s", "Hi, a new AWS account has been added to your Trackit Account. " +
			"You can connect to your account to manage it : https://re.trackit.io/")
		err := mail.SendMail(userMail, mailSubject, mailBody, ctx)
		if err != nil {
			logger.Error("Failed to send email.", err.Error())
			return err
		}
	} else {
		dbForgottenPassword, token, err := resetPasswordGenerator(ctx, tx, newUserId)
		mailSubject := "You are invited to join Trackit"
		mailBody := fmt.Sprintf("Hi, you have been invited to join trackit. Please follow this link to create" +
			" your account: https://re.trackit.io/reset/%d/%s.", dbForgottenPassword.ID, token)
		err = mail.SendMail(userMail, mailSubject, mailBody, ctx)
		if err != nil {
			logger.Error("Failed to send viewer password email.", err.Error())
			return err
		}
	}
	return nil
}

// inviteUserWithValidBody tries to share an account with a specific user
func inviteUserWithValidBody(request *http.Request, body inviteUserRequest, tx *sql.Tx, user User) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	result, guestId, err := checkUserWithEmail(request.Context(), tx, body.Email)
	if err == nil {
		if result {
			isAlreadyShared, err := checkSharedAccount(request.Context(), tx, body.AccountId, guestId)
			if err != nil {
				return 403, errors.New("An error occured while inviting a user. Please try again.")
			} else if isAlreadyShared {
				return 200, errors.New("You are already sharing this account with this user.")
			}
			err = addAccountToGuest(request.Context(), tx, body.AccountId, body.PermissionLevel, guestId, user.Id)
			if err == nil {
				err = sendMailNotification(request.Context(), tx, body.Email,true, 0)
				if err != nil {
					logger.Error("Error occured while sending an email to an existing user.", err.Error())
					return 403, errors.New("An error occured while inviting a user. Please, try again.")
				}
				return 200, nil
			} else {
				logger.Error("Error occured while adding account to an existing user.", err.Error())
				return 403, errors.New("An error occured while inviting a user. Please, try again.")
			}
		} else {
			newUserId, err := createAccountForGuest(request.Context(), tx, body.Email, body.AccountId, body.PermissionLevel, user)
			if err == nil {
				err = sendMailNotification(request.Context(), tx, body.Email,false, newUserId)
				if err != nil {
					logger.Error("Error occured while sending an email to a new user.", err.Error())
					return 403, errors.New("An error occured while inviting a new user. Please, try again.")
				}
				return 200, nil
			} else {
				logger.Error("Error occured while creating new account for a guest.", err.Error())
				return 403, errors.New("An error occured while inviting a new user. Please, try again.")
			}
		}
	} else {
		logger.Error("Error occured while checking if user already exist.", err.Error())
		return 403, errors.New("An error occured while inviting a new user. Please, try again.")
	}
}
