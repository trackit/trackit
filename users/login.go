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
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
)

// loginRequestBody is the expected request body for the LogIn route handler.
type loginRequestBody struct {
	Email    string `json:"email"    req:"nonzero"`
	Password string `json:"password" req:"nonzero"`
}

// loginResponseBody is the response body in case LogIn succeeds.
type loginResponseBody struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func init() {
	routes.MethodMuxer{
		http.MethodPost: routes.H(logIn).With(
			routes.RequestContentType{"application/json"},
			routes.RequestBody{loginRequestBody{"example@example.com", "pA55w0rd"}},
			db.RequestTransaction{db.Db},
			routes.Documentation{
				Summary:     "log in as a user",
				Description: "Logs a user in based on an e-mail/password couple and returns a JWT token and the user's data.",
			},
		),
	}.H().Register("/user/login")
}

// LogIn handles users attempting to log in. It shall return a valid token the
// caller can then use to call other routes.
func logIn(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body loginRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	return logInWithValidBody(request, body, tx)
}

// logInWithValidBody tries to authenticate and log a user in using a
// validated login request.
func logInWithValidBody(request *http.Request, body loginRequestBody, tx *sql.Tx) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	user, err := GetUserWithEmailAndPassword(request.Context(), tx, body.Email, body.Password)
	if err == nil {
		return logAuthenticatedUserIn(request, user)
	} else {
		logger.Warning("Authentication failure.", struct {
			Email string `json:"user"`
		}{user.Email})
		return 403, errors.New("Authentication failure.")
	}
}

// logAuthenticatedUserIn generates a token for a user that's already been
// authenticated.
func logAuthenticatedUserIn(request *http.Request, user User) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	token, err := generateToken(user)
	if err == nil {
		logger.Info("User logged in.", user)
		return 200, loginResponseBody{
			User:  user,
			Token: token,
		}
	} else {
		logger.Error("Failed to generate token.", err.Error())
		return 500, errors.New("Failed to generate token.")
	}
}

// TestToken tests a token's validity. For a valid token, it returns the user
// the token belongs to.
func me(request *http.Request, a routes.Arguments) (int, interface{}) {
	return 200, a[AuthenticatedUser].(User)
}
