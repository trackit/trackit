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
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/routes"
)

func patchUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body createUserRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	user := a[AuthenticatedUser].(User)
	return patchUserWithValidBody(request, user, body, tx)
}

func patchUserWithValidBody(request *http.Request, user User, body createUserRequestBody, tx *sql.Tx) (int, interface{}) {
	ctx := request.Context()
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByID(tx, user.Id)
	if err != nil {
		l.Error("Failed to find bill repository to update.", err.Error())
		return http.StatusInternalServerError, errors.New("failed to find user in database")
	}
	user, err = UpdateUserWithPassword(ctx, tx, dbUser, body.Email, body.Password)
	if err == nil {
		l.Info("User updated.", user)
		return http.StatusOK, user
	} else {
		l.Error(err.Error(), nil)
		return http.StatusInternalServerError, errors.New("Failed to update user.")
	}
}
