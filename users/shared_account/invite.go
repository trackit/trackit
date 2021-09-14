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

package shared_account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"github.com/trackit/jsonlog"
	"golang.org/x/crypto/bcrypt"

	"github.com/trackit/trackit/mail"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/users"
)

var (
	ErrorInviteNewUser = errors.New("An error occurred while inviting a new user. Please, try again.")
	ErrorInviteUser    = errors.New("An error occurred while inviting a user. Please, try again.")
	ErrorAlreadyShared = errors.New("You are already sharing this account with this user.")
)

const bCryptCost = 12

// getPasswordHash generates a hash string for a given password.
func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
	return string(hash), err
}

// checkUserWithEmailAndAccountType checks if user already exist.
// true is returned if invited user already exist.
func checkUserWithEmailAndAccountType(ctx context.Context, db models.XODB, userEmail string, accountType string, user users.User) (bool, int, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByEmailAccountType(db, userEmail, accountType)
	if err == sql.ErrNoRows {
		return false, 0, nil
	} else if err != nil {
		logger.Error("Error getting user from database.", err.Error())
		return false, 0, err
	} else {
		if user.Id != dbUser.ID {
			return true, dbUser.ID, nil
		} else {
			logger.Warning("User tries to share an account with himself", nil)
			return false, dbUser.ID, errors.New("You can't share an account with yourself")
		}
	}
}

// checkSharedAccount checks if an account is already shared with a user.
// true is returned if invited user already have an access to this account.
func checkSharedAccount(ctx context.Context, db models.XODB, accountId int, guestId int) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount, err := models.AwsAccountByID(db, accountId)
	if err != nil {
		logger.Error("Error while retrieving AWS account from DB", err)
		return false, err
	}
	if dbAwsAccount.UserID == guestId {
		logger.Warning("User tries to share an account with the owner of the account", nil)
		return false, errors.New("You can't share an account with this user")
	}
	dbSharedAccounts, err := models.SharedAccountsByUserID(db, guestId)
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
	return false, nil
}

// addAccountToGuest adds an entry in shared_account table allowing a user
// to share an access to all or part of his account
func addAccountToGuest(ctx context.Context, db *sql.Tx, accountId int, permissionLevel int, guestId int) (models.SharedAccount, error) {
	dbSharedAccount := models.SharedAccount{
		AccountID:      accountId,
		UserID:         guestId,
		UserPermission: permissionLevel,
	}
	err := dbSharedAccount.Insert(db)
	return dbSharedAccount, err
}

// createAccountForGuest creates an account for invited user who do not already own an account
func createAccountForGuest(ctx context.Context, db *sql.Tx, body InviteUserRequest, accountId int) (int, models.SharedAccount, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var sharedAccount models.SharedAccount
	tempPassword := uuid.NewV1().String()
	usr, err := users.CreateUserWithPassword(ctx, db, body.Email, tempPassword, "", "trackit")
	if err == nil {
		sharedAccount, err = addAccountToGuest(ctx, db, accountId, body.PermissionLevel, usr.Id)
		if err != nil {
			logger.Error("Error occurred while adding account to an newly created user.", err.Error())
			return 0, models.SharedAccount{}, err
		}
	} else {
		logger.Error("Error occurred while creating an automatic new account.", err.Error())
		return 0, models.SharedAccount{}, err
	}
	return usr.Id, sharedAccount, nil
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
func sendMailNotification(ctx context.Context, tx *sql.Tx, userMail string, userNew bool, newUserId int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	if userNew {
		mailSubject := "An AWS account has been added to your Trackit account"
		mailBody := "Hi, a new AWS account has been added to your Trackit Account. " +
			"You can connect to your account to manage it : https://re.trackit.io/"
		err := mail.SendMail(userMail, mailSubject, mailBody, ctx)
		if err != nil {
			logger.Error("Failed to send email.", err.Error())
			return err
		}
	} else {
		dbForgottenPassword, token, err := resetPasswordGenerator(ctx, tx, newUserId)
		if err != nil {
			logger.Error("Failed to create reset password token.", err.Error())
			return err
		}
		mailSubject := "You are invited to join Trackit"
		mailBody := fmt.Sprintf("Hi, you have been invited to join trackit. Please follow this link to create"+
			" your account: https://re.trackit.io/reset/%d/%s.", dbForgottenPassword.ID, token)
		err = mail.SendMail(userMail, mailSubject, mailBody, ctx)
		if err != nil {
			logger.Error("Failed to send viewer password email.", err.Error())
			return err
		}
	}
	return nil
}

// inviteUserAlreadyExist handles sharing for user that already exists
func inviteUserAlreadyExist(ctx context.Context, tx *sql.Tx, body InviteUserRequest, accountId int, guestId int) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	isAlreadyShared, err := checkSharedAccount(ctx, tx, accountId, guestId)
	if err != nil {
		return http.StatusForbidden, ErrorInviteUser
	} else if isAlreadyShared {
		return http.StatusBadRequest, ErrorAlreadyShared
	}
	sharedAccount, err := addAccountToGuest(ctx, tx, accountId, body.PermissionLevel, guestId)
	if err == nil {
		err = sendMailNotification(ctx, tx, body.Email, true, 0)
		if err != nil {
			logger.Error("Error occurred while sending an email to an existing user.", err.Error())
			return http.StatusForbidden, ErrorInviteUser
		}
		return http.StatusOK, sharedAccount
	} else {
		logger.Error("Error occurred while adding account to an existing user.", err.Error())
		return http.StatusForbidden, ErrorInviteUser
	}
}

// inviteNewUser handles sharing for user that do not already exist
func inviteNewUser(ctx context.Context, tx *sql.Tx, body InviteUserRequest, accountId int) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	newUserId, newUser, err := createAccountForGuest(ctx, tx, body, accountId)
	if err == nil {
		err = sendMailNotification(ctx, tx, body.Email, false, newUserId)
		if err != nil {
			logger.Error("Error occurred while sending an email to a new user.", err.Error())
			return http.StatusForbidden, ErrorInviteNewUser
		}
		return http.StatusOK, newUser
	} else {
		logger.Error("Error occurred while creating new account for a guest.", err.Error())
		return http.StatusForbidden, ErrorInviteNewUser
	}
}

// InviteUserWithValidBody tries to share an account with a specific user
func InviteUserWithValidBody(request *http.Request, body InviteUserRequest, accountId int, tx *sql.Tx, user users.User) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	security, err := safetyCheckByAccountIdAndPermissionLevel(request.Context(), tx, accountId, body, user)
	if err != nil {
		return http.StatusBadRequest, err
	} else if !security {
		return http.StatusForbidden, errors.New("You do not have permission to edit this sharing")
	}
	if !checkPermissionLevel(body.PermissionLevel) {
		logger.Info("Non existing user permission", nil)
		return http.StatusBadRequest, ErrorInviteUser
	}
	result, guestId, err := checkUserWithEmailAndAccountType(request.Context(), tx, body.Email, body.Origin, user)
	if err == nil {
		if result {
			code, res := inviteUserAlreadyExist(request.Context(), tx, body, accountId, guestId)
			return code, res
		} else {
			code, res := inviteNewUser(request.Context(), tx, body, accountId)
			return code, res
		}
	} else {
		logger.Error("Error occurred while checking body elements.", err.Error())
		return http.StatusForbidden, ErrorInviteNewUser
	}
}
