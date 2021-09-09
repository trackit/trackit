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

package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/mail"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
)

var (
	nbHoursValidityForgottenToken = 1.0
)

// forgottenPasswordRequestBody is the expected request body for the forgottenPassword route handler.
type forgottenPasswordRequestBody struct {
	Email  string `json:"email"   req:"nonzero"`
	Origin string `json:"origin"  req:"nonzero"`
}

// resetPasswordRequestBody is the expected request body for the resetPassword route handler.
type resetPasswordRequestBody struct {
	Id       int    `json:"id"       req:"nonzero"`
	Token    string `json:"token"    req:"nonzero"`
	Password string `json:"password" req:"nonzero"`
}

func init() {
	routes.MethodMuxer{
		http.MethodPost: routes.H(forgottenPassword).With(
			routes.RequestContentType{"application/json"},
			routes.RequestBody{forgottenPasswordRequestBody{
				Email:  "example@example.com",
				Origin: "trackit",
			}},
			db.RequestTransaction{db.Db},
			routes.Documentation{
				Summary:     "request a forgotten password reset",
				Description: "Sends an email to a user with a link to reset his forgotten password",
			},
		),
	}.H().Register("/user/password/forgotten")
	routes.MethodMuxer{
		http.MethodPost: routes.H(resetPassword).With(
			routes.RequestContentType{"application/json"},
			routes.RequestBody{resetPasswordRequestBody{2, "abcdefg", "password"}},
			db.RequestTransaction{db.Db},
			routes.Documentation{
				Summary:     "reset a forgotten password",
				Description: "Allows a user to reset a forgotten password using a temporary token",
			},
		),
	}.H().Register("/user/password/reset")
}

func cleanExpiredTokens() {
	transaction, err := db.Db.BeginTx(context.Background(), nil)
	if err == nil {
		now := time.Now()
		expire := now.Add(time.Hour * time.Duration(-1*int64(nbHoursValidityForgottenToken)))
		models.DeleteExpiredForgottenPassword(transaction, expire)
		defer func() {
			rec := recover()
			if rec != nil {
				transaction.Rollback()
			} else {
				transaction.Commit()
			}
		}()
	}
}

// forgottenPassword handles users requesting to reset a forgotten password
func forgottenPassword(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body forgottenPasswordRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	go cleanExpiredTokens()
	return forgottenPasswordWithValidBody(request, body, tx)
}

func forgottenPasswordWithValidBody(request *http.Request, body forgottenPasswordRequestBody, tx *sql.Tx) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	user, err := GetUserWithEmailAndOrigin(request.Context(), tx, body.Email, body.Origin)
	if err != nil {
		logger.Warning("Forgotten password request failure", struct {
			Email string `json:"user"`
			Error string `json:"error"`
		}{body.Email, err.Error()})
		return 404, errors.New("User not found")
	}
	return createForgottenPasswordEntry(request, body, tx, user)
}

func createForgottenPasswordEntry(request *http.Request, body forgottenPasswordRequestBody, tx *sql.Tx, user User) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	token := uuid.NewV1().String()
	tokenHash, err := getPasswordHash(token)
	if err != nil {
		logger.Error("Failed to create token hash.", err.Error())
		return 500, errors.New("Failed to create token hash")
	}
	dbForgottenPassword := models.ForgottenPassword{
		UserID:  user.Id,
		Token:   tokenHash,
		Created: time.Now(),
	}
	err = dbForgottenPassword.Insert(tx)
	if err != nil {
		logger.Error("Failed to insert forgotten password token in database.", err.Error())
		return 500, errors.New("Failed to create forgotten password token")
	}
	mailSubject := "Reset your Trackit password"
	mailBody := fmt.Sprintf("Please follow this link to recover your password: https://re.trackit.io/reset/%d/%s. This link is valid for an hour.", dbForgottenPassword.ID, token)
	err = mail.SendMail(user.Email, mailSubject, mailBody, request.Context())
	if err != nil {
		logger.Error("Failed to send password recovery email.", err.Error())
		return 500, errors.New("Failed to send password recovery email")
	}
	return 200, nil
}

// forgottenPassword handles users attempting to reset a forgotten password
func resetPassword(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body resetPasswordRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	return resetPasswordWithValidBody(request, body, tx)
}

func resetPasswordWithValidBody(request *http.Request, body resetPasswordRequestBody, tx *sql.Tx) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	forgottenPassword, err := models.ForgottenPasswordByID(tx, body.Id)
	if err != nil {
		logger.Warning("Forgotten password token not found", struct {
			Token string `json:"token"`
			Error string `json:"error"`
		}{body.Token, err.Error()})
		return 404, errors.New("Reset token not found")
	}
	err = passwordMatchesHash(body.Token, forgottenPassword.Token)
	delta := time.Now().Sub(forgottenPassword.Created)
	expired := delta.Hours() > nbHoursValidityForgottenToken
	if err != nil || expired == true {
		logger.Warning("Invalid token", struct {
			Token string `json:"token"`
		}{body.Token})
		return 400, errors.New("Invalid token")
	}
	return updatePasswordFromForgottenToken(request, body, tx, forgottenPassword)
}

func updatePasswordFromForgottenToken(request *http.Request, body resetPasswordRequestBody, tx *sql.Tx, forgottenPassword *models.ForgottenPassword) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	user, err := GetUserWithId(tx, forgottenPassword.UserID)
	if err != nil {
		logger.Warning("Unable to retrieve user associated to forgotten password token", struct {
			Token string `json:"token"`
			Error string `json:"error"`
		}{forgottenPassword.Token, err.Error()})
		return 404, errors.New("Unable to retrieve user")
	}
	err = user.UpdatePassword(tx, body.Password)
	if err != nil {
		logger.Warning("Unable to update user password", err.Error())
		return 500, errors.New("Unable to update user")
	}
	err = forgottenPassword.Delete(tx)
	if err != nil {
		logger.Warning("Unable to delete forgotten password token", err.Error())
	}
	return 200, nil
}
