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

package ec2

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

type (
	// ec2QueryParams will store the parsed query params
	ec2QueryParams struct {
		accountList []string
		indexList   []string
		date        time.Time
	}

	// ec2UnusedQueryParams will store the parsed query params
	ec2UnusedQueryParams struct {
		accountList []string
		indexList   []string
		date        time.Time
		count       int
	}
)

var (
	// ec2QueryArgs allows to get required queryArgs params
	ec2QueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
	}

	// ec2UnusedQueryArgs allows to get required queryArgs params
	ec2UnusedQueryArgs = []routes.QueryArg{
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
		http.MethodGet: routes.H(getEc2Instances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(ec2QueryArgs),
			routes.Documentation{
				Summary:     "get the list of EC2 instances",
				Description: "Responds with the list of EC2 instances based on the queryparams passed to it",
			},
		),
	}.H().Register("/ec2")
	routes.MethodMuxer{
		http.MethodGet: routes.H(getEc2UnusedInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(ec2UnusedQueryArgs),
			routes.Documentation{
				Summary:     "get the list of the most unused EC2 instances of a month",
				Description: "Responds with the list of the most unused EC2 instances of a month based on the queryparams passed to it",
			},
		),
	}.H().Register("/ec2/unused")
}

// getEc2Instances returns the list of EC2 reports based on the query params, in JSON format.
func getEc2Instances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := ec2QueryParams{
		accountList: []string{},
		date:        a[routes.DateQueryArg].(time.Time),
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.accountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	returnCode, report, err := getEc2Data(request, parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}

// getEc2UnusedInstances returns the list of EC2 reports based on the query params, in JSON format.
func getEc2UnusedInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := ec2UnusedQueryParams{
		accountList: []string{},
		date:        a[routes.DateQueryArg].(time.Time),
		count:       -1,
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.accountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	if a[ec2UnusedQueryArgs[2]] != nil {
		parsedParams.count = a[ec2UnusedQueryArgs[2]].(int)
	}
	returnCode, report, err := getEc2UnusedData(request, parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
