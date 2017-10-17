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
	routes.Register(
		"/user",
		createUser,
		routes.RequireMethod{"POST"},
		routes.WithErrorBody{},
		db.WithTransaction{db.Db},
	)
}

type createUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body createUserRequestBody
	err := decodeRequestBody(request, &body)
	tx := a[db.Transaction].(*sql.Tx)
	if err == nil && body.valid() {
		return createUserWithValidBody(request, body, tx)
	} else {
		return 400, errors.New("Body is invalid.")
	}
}

func (body createUserRequestBody) valid() bool {
	return body.Email != "" && body.Password != ""
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
