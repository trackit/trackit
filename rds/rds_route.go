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

package rds

import (
	"database/sql"
	"net/http"
	"errors"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
	"time"
)

type (
	// rdsQueryParams will store the parsed query params
	rdsQueryParams struct {
		accountList []string
		indexList   []string
		date        time.Time
	}

	// rdsUnusedQueryParams will store the parsed query params
	rdsUnusedQueryParams struct {
		accountList []string
		indexList   []string
		date        time.Time
		count       int
		by          string
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
		routes.QueryArg{
			Name:        "by",
			Type:        routes.QueryArgString{},
			Description: "Element choose to sort unused data, (cpu, freespace), default cpu",
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
	parsedParams := rdsQueryParams{
		accountList: []string{},
		date:        a[rdsQueryArgs[1]].(time.Time),
	}
	if a[rdsQueryArgs[0]] != nil {
		parsedParams.accountList = a[rdsQueryArgs[0]].([]string)
	}
	returnCode, report, err := getRdsData(request, parsedParams, user, tx)
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
	parsedParams := rdsUnusedQueryParams{
		accountList: []string{},
		date:        a[rdsUnusedQueryArgs[1]].(time.Time),
		count:       -1,
		by:          "cpu",
	}
	if a[rdsUnusedQueryArgs[0]] != nil {
		parsedParams.accountList = a[rdsUnusedQueryArgs[0]].([]string)
	}
	if a[rdsUnusedQueryArgs[2]] != nil {
		parsedParams.count = a[rdsUnusedQueryArgs[2]].(int)
	}
	if a[rdsUnusedQueryArgs[3]] != nil {
		parsedParams.by = a[rdsUnusedQueryArgs[3]].(string)
		if parsedParams.by != "cpu" && parsedParams.by != "freespace" {
			return http.StatusBadRequest, errors.New("bad argument for the query arg 'by'")
		}
	}
	returnCode, report, err := getRdsUnusedData(request, parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
