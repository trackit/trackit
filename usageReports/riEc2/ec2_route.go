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

// Package riEc2 implements the generation of EC2 reserved instances usage reports and of the corresponding /ri/ec2 route
package riEc2

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type (
	// ReservedInstancesQueryParams will store the parsed query params
	ReservedInstancesQueryParams struct {
		AccountList []string
		IndexList   []string
		Date        time.Time
	}
)

var (
	// reservedInstancesQueryArgs allows to get required queryArgs params
	reservedInstancesQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.DateQueryArg,
	}
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getReservedInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(reservedInstancesQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get the list of Reserved Instances",
				Description: "Responds with the list of Reserved Instances based on the queryparams passed to it",
			},
		),
	}.H().Register("/ri/ec2")
}

// getReservedInstances returns the list of Reserved Instances reports based on the query params, in JSON format.
func getReservedInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := ReservedInstancesQueryParams{
		AccountList: []string{},
		Date:        a[routes.DateQueryArg].(time.Time),
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	returnCode, report, err := GetReservedInstancesData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}
