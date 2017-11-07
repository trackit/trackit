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

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
)

type WithAuthenticatedUser struct{}

type withAuthenticatedUserArgumentKey uint

const (
	AuthenticatedUser = withAuthenticatedUserArgumentKey(iota)
)

func (d WithAuthenticatedUser) Decorate(h routes.IntermediateHandler) routes.IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		auth := r.Header["Authorization"]
		tx := a[db.Transaction].(*sql.Tx)
		if auth != nil && len(auth) == 1 {
			tokenString := auth[0]
			if user, err := testToken(tx, tokenString); err == nil {
				a[AuthenticatedUser] = user
				return h(w, r, a)
			}
		}
		return 401, errors.New("Invalid or missing token.")
	}
}
