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
	"github.com/trackit/trackit-server/routes"
	"net/http"
	"github.com/trackit/trackit-server/db"
	"database/sql"
)

// TokenQueryArg allows to get the Trackit authentication token in the URL.
// We can use it to check a token validity.
var TokenQueryArg = routes.QueryArg{
	Name:        "token",
	Type:        routes.QueryArgString{},
	Description: "Token to check.",
	Optional:    false,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(checkTokenValidity).With(
			db.RequestTransaction{db.Db},
			routes.Documentation{
				Summary:     "Check token validity",
				Description: "Check the validity of an authentication token",
			},
			routes.QueryArgs{TokenQueryArg},
		),
	}.H().Register("/token/check")
}

type awsTokenCheck struct {
	awsToken string
}

// checkTokenValidity is a route handler which returns the user if the
// Trackit authentication token is valid.
func checkTokenValidity(r *http.Request, a routes.Arguments) (int, interface{}) {
	tx := a[db.Transaction].(*sql.Tx)
	token := a[TokenQueryArg].(string)
	if user, err := testToken(tx, token); err == nil {
		return 200, user
	} else {
		return http.StatusUnauthorized, err
	}
}

