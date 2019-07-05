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

package ebs

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

type (
	// EbsQueryParams will store the parsed query params
	EbsQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
	}

	// EbsUnusedQueryParams will store the parsed query params
	EbsUnusedQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
		Count       int
	}
)

var (
	// ebsQueryArgs allows to get required queryArgs params
	ebsQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
	}

	// ebsUnusedQueryArgs allows to get required queryArgs params
	ebsUnusedQueryArgs = []routes.QueryArg{
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
		http.MethodGet: routes.H(getEbsInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(ebsQueryArgs),
			routes.Documentation{
				Summary:     "get the list of EBS instances",
				Description: "Responds with the list of EBS instances based on the queryparams passed to it",
			},
		),
	}.H().Register("/ebs")
	routes.MethodMuxer{
		http.MethodGet: routes.H(getEbsUnusedInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(ebsUnusedQueryArgs),
			routes.Documentation{
				Summary:     "get the list of the most unused EBS instances of a month",
				Description: "Responds with the list of the most unused EBS instances of a month based on the queryparams passed to it",
			},
		),
	}.H().Register("/ebs/unused")
}

// getEbsInstances returns the list of EBS reports based on the query params, in JSON format.
func getEbsInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := EbsQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	returnCode, report, err := GetEbsData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}

// getEbsUnusedInstances returns the list of EBS reports based on the query params, in JSON format.
func getEbsUnusedInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := EbsUnusedQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
		Count:       -1,
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	if a[ebsUnusedQueryArgs[2]] != nil {
		parsedParams.Count = a[ebsUnusedQueryArgs[2]].(int)
	}
	returnCode, report, err := GetEbsUnusedData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
