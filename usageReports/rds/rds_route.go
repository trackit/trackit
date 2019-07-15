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

package rds

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/trackit/trackit-server/cache"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

type (
	// RdsQueryParams will store the parsed query params
	RdsQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
	}

	// RdsUnusedQueryParams will store the parsed query params
	RdsUnusedQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
		Count       int
	}
)

var (
	// rdsQueryArgs allows to get required queryArgs params
	rdsQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
	}

	// rdsUnusedQueryArgs allows to get required queryArgs params
	rdsUnusedQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
		routes.QueryArg{
			Name:        "count",
			Type:        routes.QueryArgInt{},
			Description: "Number of element in the response, all if not precised or negative",
			Optional:    true,
		},
	}
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getRdsReport).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(rdsQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get a RDS report of a month",
				Description: "Responds with the a RDS report for the account and date specified in the request",
			},
		),
	}.H().Register("/rds")
	routes.MethodMuxer{
		http.MethodGet: routes.H(getRdsUnusedInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(rdsUnusedQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get the list of the most unused RDS instances of a month",
				Description: "Responds with the list of the most unused RDS instances of a month based on the queryparams passed to it",
			},
		),
	}.H().Register("/rds/unused")
}

// getRdsReport returns the list of RDS reports based on the query params, in JSON format.
func getRdsReport(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := RdsQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	returnCode, report, err := GetRdsData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}

// getRdsUnusedInstances returns the list of the most unused RDS instances based on the query params, in JSON format.
func getRdsUnusedInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := RdsUnusedQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
		Count:       -1,
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	if a[rdsUnusedQueryArgs[2]] != nil {
		parsedParams.Count = a[rdsUnusedQueryArgs[2]].(int)
	}
	returnCode, report, err := GetRdsUnusedData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
