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
	"github.com/trackit/trackit-server/config"
)

// inviteUserRequest is the expected request body for the invite user route handler.
type inviteUserRequest struct {
	Email              string `json:"email" req:"nonzero"`
	AccountId          int `json:"accountId"`
	PermissionLevel    int `json:"level"`
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

// checkuserWithEmail checks if user already exist. if user exists, user Id is return
func checkUserWithEmail(ctx context.Context, db models.XODB, userEmail string) (res bool, userId int, err error) {
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

func checkSharedAccount(ctx context.Context, db models.XODB, accountId int, userId int) (res bool, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	const sqlstr = `SELECT account_id, user_id FROM shared_account WHERE user_id = ?`
	res, err := db.Query(sqlstr, userId)
	defer res.Close()
	if err != nil {
		return "", err
	}
	var token string
	for res.Next() {
		err := res.Scan(&token)
		if err != nil {
			return "", err
		}
	}
	return token, nil
}

// addAccountToGuest adds an entry in shared_account table with element that enable a user
// to share an access to all or part of his account
func addAccountToGuest(ctx context.Context, db *sql.Tx, accountId int, permissionLevel int, guestId int, ownerId int) (err error) {
	//TODO : Check if user already have an access to the account if so, abort and return 200, already sharing with this user
	fmt.Print("--- GUEST ID : ")
	fmt.Print(guestId)
	isAlreadyShared, err := checkSharedAccount(ctx, db, guestId, accountId)
	_ = isAlreadyShared
	if err != nil {
		return err
	}
	const sqlstr = `INSERT INTO shared_account(
			account_id, owner_id, user_id, user_permission, account_status
		) VALUES (?, ?, ?, ?, ?)`
	res, err := db.Exec(sqlstr, accountId, ownerId, guestId, permissionLevel, 0)
	_ = res
	if err != nil {
		return err
	}
	return nil
}

// createAccountForGuest creates an account for invited user who do not already own an account
func createAccountForGuest(ctx context.Context, db *sql.Tx, userMail string, accountId int, permissionLevel int, user User) (newUserId int, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	tempPassword := uuid.NewV1().String()
	usr, err := CreateUserWithPassword(ctx, db, userMail, tempPassword, "")
	if err == nil {
		err = addAccountToGuest(ctx, db, accountId, permissionLevel, usr.Id, user.Id)
		if err != nil {
			logger.Warning("Error occured while adding account to an newly created user.", err)
			return 0, err
		}
	} else {
		logger.Warning("Error occured while creating an automatic new account.", err)
		return 0, err
	}
	return usr.Id,nil
}

// resetPasswordGenerator returns a reset password token. It is used in order to
// create an account and let the user choose his own password
func resetPasswordGenerator(ctx context.Context, tx *sql.Tx, newUserId int) (dbForgottenPassword models.ForgottenPassword, token string, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	token = uuid.NewV1().String()
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
		return dbForgottenPassword, token, nil
	} else {
		logger.Error("Failed to insert forgotten password", err.Error())
		return dbForgottenPassword, "", err
	}
}

// sendMailNotification sends an email to user how has been invited to access a AWS account on trackit.io
func sendMailNotification(ctx context.Context, tx *sql.Tx, userMail string, userNew bool, newUserId int) (err error) {
	//TODO : Needs to be removed before merging. This is for tests purpose ONLY ----
	config.SmtpAddress = "email-smtp.us-west-2.amazonaws.com"
	config.SmtpPort = "587"
	config.SmtpUser = "AKIAJUT3EB3EH2V6SX5A"
	config.SmtpPassword = "Ahk+AhVhVnjb/gtKVAugFjJZQt9qvWQTGAvJy19kwmyU"
	config.SmtpSender = "team@trackit.io"
	//TODO : ---- Ends of TODO
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbForgottenPassword, token, err := resetPasswordGenerator(ctx, tx, newUserId)
	if userNew {
		mailSubject := "An AWS account has been added to your Trackit account"
		mailBody := fmt.Sprintf("%s", "Hi, a new AWS account has been added to your Trackit Account. " +
			"You can connect to your account to manage it : https://re.trackit.io/")
		err = mail.SendMail(userMail, mailSubject, mailBody, ctx)
		if err != nil {
			logger.Error("Failed to send email.", err.Error())
			return err
		}
	} else {
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

// logInWithValidBody tries to authenticate and log a user in using a
// validated login request.
func inviteUserWithValidBody(request *http.Request, body inviteUserRequest, tx *sql.Tx, user User) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	result, guestId, err := checkUserWithEmail(request.Context(), tx, body.Email)
	if err == nil {
		if result {
			err = addAccountToGuest(request.Context(), tx, body.AccountId, body.PermissionLevel, guestId, user.Id)
			if err == nil {
				err = sendMailNotification(request.Context(), tx, body.Email,true, 0)
				if err != nil {
					logger.Warning("Error occured while sending an email to an existing user.", err)
					return 403, errors.New("An error occured while inviting a user. Please, try again.")
				}
				return 200, "account shared"
			} else {
				logger.Warning("Error occured while adding account to an existing user.", err)
				return 403, errors.New("An error occured while inviting a user. Please, try again.")
			}
		} else {
			newUserId, err := createAccountForGuest(request.Context(), tx, body.Email, body.AccountId, body.PermissionLevel, user)
			if err == nil {
				err = sendMailNotification(request.Context(), tx, body.Email,false, newUserId)
				if err != nil {
					logger.Warning("Error occured while sending an email to a new user.", err)
					return 403, errors.New("An error occured while inviting a new user. Please, try again.")
				}
				return 200, "account created and shared"
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
