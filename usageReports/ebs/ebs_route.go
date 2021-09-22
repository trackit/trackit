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

// Package ebs implements the generation of EBS usage reports and of the corresponding /ebs route
package ebs

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
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
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getEbsSnapshots).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(ebsQueryArgs),
			routes.Documentation{
				Summary:     "get the list of EBS snapshots",
				Description: "Responds with the list of EBS snapshots based on the queryparams passed to it",
			},
		),
	}.H().Register("/ebs")
}

// getEbsSnapshots returns the list of EBS reports based on the query params, in JSON format.
func getEbsSnapshots(request *http.Request, a routes.Arguments) (int, interface{}) {
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
