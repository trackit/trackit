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

package elasticache

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
	// ElastiCacheQueryParams will store the parsed query params
	ElastiCacheQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
	}

	// ElastiCacheUnusedQueryParams will store the parsed query params
	ElastiCacheUnusedQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
		Count       int
	}
)

var (
	// elastiCacheQueryArgs allows to get required queryArgs params
	elasticacheQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
	}

	// elastiCacheUnusedQueryArgs allows to get required queryArgs params
	elasticacheUnusedQueryArgs = []routes.QueryArg{
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
		http.MethodGet: routes.H(getElastiCacheInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(elasticacheQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get the list of ElastiCache instances",
				Description: "Responds with the list of ElastiCache instances based on the queryparams passed to it",
			},
		),
	}.H().Register("/elasticache")
	routes.MethodMuxer{
		http.MethodGet: routes.H(getElastiCacheUnusedInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(elasticacheUnusedQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get the list of the most unused ElastiCache instances of a month",
				Description: "Responds with the list of the most unused ElastiCache instances of a month based on the queryparams passed to it",
			},
		),
	}.H().Register("/elasticache/unused")
}

// getElastiCacheInstances returns the list of ElastiCache reports based on the query params, in JSON format.
func getElastiCacheInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := ElastiCacheQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	returnCode, report, err := GetElastiCacheData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}

// getElastiCacheUnusedInstances returns the list of ElastiCache reports based on the query params, in JSON format.
func getElastiCacheUnusedInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := ElastiCacheUnusedQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
		Count:       -1,
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	if a[elasticacheUnusedQueryArgs[2]] != nil {
		parsedParams.Count = a[elasticacheUnusedQueryArgs[2]].(int)
	}
	returnCode, report, err := GetElastiCacheUnusedData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
