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

package instanceCount

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type (
	// InstanceCountQueryParams will store the parsed query params
	InstanceCountQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
	}

	// InstanceCountUnusedQueryParams will store the parsed query params
	InstanceCountUnusedQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
		Count       int
	}
)

var (
	// instanceCountQueryArgs allows to get required queryArgs params
	instanceCountQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
	}
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getInstanceCount).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(instanceCountQueryArgs),
			routes.Documentation{
				Summary:     "get the list of InstanceCount",
				Description: "Responds with the list of InstanceCount based on the queryparams passed to it",
			},
		),
	}.H().Register("/instanceCount")
}

// getInstanceCount returns the list of InstanceCount reports based on the query params, in JSON format.
func getInstanceCount(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := InstanceCountQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	returnCode, report, err := GetInstanceCountData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
