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
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
)

type RequireAuthenticatedUser struct{}

type authenticatedUserArgumentKey uint

const (
	AuthenticatedUser            = authenticatedUserArgumentKey(iota)
	TagRequireUserAuthentication = "require:userauth"
)

func (d RequireAuthenticatedUser) Decorate(h routes.Handler) routes.Handler {
	h.Func = d.getFunc(h.Func)
	h.Documentation = d.getDocumentation(h.Documentation)
	return h
}

func (_ RequireAuthenticatedUser) getFunc(hf routes.HandlerFunc) routes.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		logger := jsonlog.LoggerFromContextOrDefault(r.Context())
		auth := r.Header["Authorization"]
		tx := a[db.Transaction].(*sql.Tx)
		if auth != nil && len(auth) == 1 {
			tokenString := auth[0]
			if user, err := testToken(tx, tokenString); err == nil {
				a[AuthenticatedUser] = user
				return hf(w, r, a)
			} else if err != ErrCannotReadToken && err != ErrInvalidClaims {
				logger.Error("Abnormal authentication failure.", err.Error())
				return http.StatusInternalServerError, ErrFailedToValidateToken
			} else {
				return http.StatusUnauthorized, err
			}
		} else {
			return http.StatusUnauthorized, ErrMissingToken
		}
	}
}

func (_ RequireAuthenticatedUser) getDocumentation(hd routes.HandlerDocumentation) routes.HandlerDocumentation {
	if hd.Tags == nil {
		hd.Tags = make(routes.Tags)
	}
	hd.Tags[TagRequireUserAuthentication] = []string{"authenticated"}
	return hd
}
