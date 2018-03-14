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
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
)

const (
	passwordMaxLength = 12
)

var (
	ErrPasswordTooShort = errors.New(fmt.Sprintf("Password must be at least %u characters.", passwordMaxLength))
)

func init() {
	routes.MethodMuxer{
		http.MethodPost: routes.H(createUser).With(
			routes.RequestContentType{"application/json"},
			routes.RequestBody{createUserRequestBody{"example@example.com", "pa55w0rd"}},
			routes.Documentation{
				Summary:     "register a new user",
				Description: "Registers a new user using an e-mail and password, and responds with the user's data.",
			},
		),
		http.MethodPatch: routes.H(patchUser).With(
			RequireAuthenticatedUser{},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{createUserRequestBody{"example@example.com", "pa55w0rd"}},
			routes.Documentation{
				Summary:     "edit the current user",
				Description: "Edit the current user, and responds with the user's data.",
			},
		),
		http.MethodGet: routes.H(me).With(
			RequireAuthenticatedUser{},
			routes.Documentation{
				Summary:     "get the current user",
				Description: "Responds with the currently authenticated user's data.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
		routes.Documentation{
			Summary: "register or get the user",
		},
	).Register("/user")
}

type createUserRequestBody struct {
	Email    string `json:"email"    req:"nonzero"`
	Password string `json:"password" req:"nonzero"`
}

func createUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body createUserRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	return createUserWithValidBody(request, body, tx)
}

func createUserWithValidBody(request *http.Request, body createUserRequestBody, tx *sql.Tx) (int, interface{}) {
	ctx := request.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	user, err := CreateUserWithPassword(ctx, tx, body.Email, body.Password)
	if err == nil {
		logger.Info("User created.", user)
		return 200, user
	} else {
		logger.Error(err.Error(), nil)
		return 500, errors.New("Failed to create user.")
	}
}
