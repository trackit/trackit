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

package shared_account

import (
	"time"
	"fmt"
	"errors"
	"database/sql"
	"context"
	"net/http"

	"github.com/trackit/jsonlog"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/trackit/trackit-server/mail"
	"github.com/trackit/trackit-server/users"
	"github.com/trackit/trackit-server/models"
)

var (
	ErrorInviteNewUser = errors.New("An error occured while inviting a new user. Please, try again.")
	ErrorInviteUser = errors.New("An error occured while inviting a user. Please, try again.")
	ErrorAlreadyShared = errors.New("You are already sharing this account with this user.")
)

const bCryptCost = 12

// getPasswordHash generates a hash string for a given password.
func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
	return string(hash), err
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
func addAccountToGuest(ctx context.Context, db *sql.Tx, accountId int, permissionLevel int, guestId int) (error) {
	dbSharedAccount := models.SharedAccount{
		AccountID:  accountId,
		UserID:   guestId,
		UserPermission: permissionLevel,
	}
	err := dbSharedAccount.Insert(db)
	return err
}

// createAccountForGuest creates an account for invited user who do not already own an account
func createAccountForGuest(ctx context.Context, db *sql.Tx, userMail string, accountId int, permissionLevel int) (int, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	tempPassword := uuid.NewV1().String()
	usr, err := users.CreateUserWithPassword(ctx, db, userMail, tempPassword, "")
	if err == nil {
		err = addAccountToGuest(ctx, db, accountId, permissionLevel, usr.Id)
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

func inviteUserAlreadyExist(ctx context.Context, tx *sql.Tx, body InviteUserRequest, guestId int) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	isAlreadyShared, err := checkSharedAccount(ctx, tx, body.AccountId, guestId)
	if err != nil {
		return 403, ErrorInviteUser
	} else if isAlreadyShared {
		return 200, ErrorAlreadyShared
	}
	err = addAccountToGuest(ctx, tx, body.AccountId, body.PermissionLevel, guestId)
	if err == nil {
		err = sendMailNotification(ctx, tx, body.Email,true, 0)
		if err != nil {
			logger.Error("Error occured while sending an email to an existing user.", err.Error())
			return 403, ErrorInviteUser
		}
		return 200, nil
	} else {
		logger.Error("Error occured while adding account to an existing user.", err.Error())
		return 403, ErrorInviteUser
	}
}

func inviteNewUser(ctx context.Context, tx *sql.Tx, body InviteUserRequest) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	newUserId, err := createAccountForGuest(ctx, tx, body.Email, body.AccountId, body.PermissionLevel)
	if err == nil {
		err = sendMailNotification(ctx, tx, body.Email,false, newUserId)
		if err != nil {
			logger.Error("Error occured while sending an email to a new user.", err.Error())
			return 403, ErrorInviteNewUser
		}
		return 200, nil
	} else {
		logger.Error("Error occured while creating new account for a guest.", err.Error())
		return 403, ErrorInviteNewUser
	}
}

// inviteUserWithValidBody tries to share an account with a specific user
func InviteUserWithValidBody(request *http.Request, body InviteUserRequest, tx *sql.Tx, user users.User) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	result, guestId, err := checkUserWithEmail(request.Context(), tx, body.Email)
	if err == nil {
		if result {
			code, res := inviteUserAlreadyExist(request.Context(), tx, body, guestId)
			return code, res
		} else {
			code, res := inviteNewUser(request.Context(), tx, body)
			return code, res
		}
	} else {
		logger.Error("Error occured while checking if user already exist.", err.Error())
		return 403, ErrorInviteNewUser
	}
}
