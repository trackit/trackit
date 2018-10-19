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

package es

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

// esQueryParams will store the parsed query params
type (
	esQueryParams struct {
		accountList []string
		indexList   []string
		date        time.Time
	}

	esUnusedQueryParams struct {
		accountList []string
		indexList   []string
		date        time.Time
		count       int
	}
)

var (
	esQueryArgs = []routes.QueryArg{
		routes.AwsAccountIdsOptionalQueryArg,
		routes.DateQueryArg,
	}

	esUnusedQueryArgs = []routes.QueryArg{
		routes.AwsAccountIdsOptionalQueryArg,
		routes.DateQueryArg,
		routes.QueryArg{
			Name:        "count",
			Type:        routes.QueryArgInt{},
			Description: "Number of element in the response, all if not specified or negative",
			Optional:    true,
		},
	}
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getESReport).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(esQueryArgs),
			routes.Documentation{
				Summary:     "get the latest ES report",
				Description: "Responds with the latest ES report for the account specified in the request",
			},
		),
	}.H().Register("/es")
	routes.MethodMuxer{
		http.MethodGet: routes.H(getESUnusedInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(esUnusedQueryArgs),
			routes.Documentation{
				Summary:     "get the list of the most unused ES domains of a month",
				Description: "Responds with the list of the most unused ES domains of a month based on the queryparams passed to it",
			},
		),
	}.H().Register("/es/unused")
}

// getRdsReport returns the list of RDS reports based on the query params, in JSON format.
func getESReport(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := esQueryParams{
		accountList: []string{},
		date:        a[esQueryArgs[1]].(time.Time),
	}
	if a[esQueryArgs[0]] != nil {
		parsedParams.accountList = a[esQueryArgs[0]].([]string)
	}
	returnCode, report, err := getEsData(request, parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}

// getRdsUnusedInstances returns the list of the most unused RDS instances based on the query params, in JSON format.
func getESUnusedInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := esUnusedQueryParams{
		accountList: []string{},
		date:        a[esUnusedQueryArgs[1]].(time.Time),
		count:       -1,
	}
	if a[esUnusedQueryArgs[0]] != nil {
		parsedParams.accountList = a[esUnusedQueryArgs[0]].([]string)
	}
	if a[esUnusedQueryArgs[2]] != nil {
		parsedParams.count = a[esUnusedQueryArgs[2]].(int)
	}
	returnCode, report, err := getEsUnusedData(request, parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
