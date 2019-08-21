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

package lambda

import (
	"database/sql"
	"github.com/trackit/trackit/pagination"
	"net/http"
	"time"

	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type (
	// LambdaQueryParams will store the parsed query params
	LambdaQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
		Pagination  pagination.Pagination
	}
)

var (
	// lambdaQueryArgs allows to get required queryArgs params
	lambdaQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
		routes.PaginationPageQueryArg,
		routes.PaginationNumberElementsQueryArg,
	}
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getLambdaFunctions).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(lambdaQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get the list of Lambda functions",
				Description: "Responds with the list of Lambda functions based on the queryparams passed to it",
			},
		),
	}.H().Register("/lambda")
}

// getLambdaFunctions returns the list of Lambda reports based on the query params, in JSON format.
func getLambdaFunctions(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := LambdaQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
		Pagination:  pagination.NewPagination(a),
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	returnCode, report, parsedParams, err := GetLambdaData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, pagination.WrapPagination(parsedParams.Pagination, report)
	}
}
