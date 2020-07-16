//   Copyright 2020 MSolution.IO
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

package resourcesTagging

import (
	"net/http"
	"database/sql"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type (
	// ResourcesQueryParams will store the parsed query params
	ResourcesQueryParams struct {
		AccountList  []string
		IndexList    []string
		Region       []string
		ResourceType []string
	}
)

var (
	// resourcesQueryArgs allows to get required queryArgs params
	resourcesQueryArgs = []routes.QueryArg{
		routes.AwsAccountsOptionalQueryArg,
		routes.QueryArg{
			Name:        "region",
			Type:        routes.QueryArgStringSlice{},
			Description: "The region of the resource",
			Optional:    true,
		},
		routes.QueryArg{
			Name:        "resourceType",
			Type:        routes.QueryArgStringSlice{},
			Description: "The type of the resource",
			Optional:    true,
		},
	}
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetResources).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(resourcesQueryArgs),
			routes.Documentation{
				Summary:     "get list of resources",
				Description: "Responds with the list of resources based on the queryparams passed to it",
			},
		),
	}.H().Register("/tagging/resources")
}

// routeGetResources returns the list of resources based on the query params, in JSON format.
func routeGetResources(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := ResourcesQueryParams{
		AccountList:  []string{},
		Region:       []string{},
		ResourceType: []string{},
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	if a[resourcesQueryArgs[1]] != nil {
		parsedParams.Region = a[resourcesQueryArgs[1]].([]string)
	}
	if a[resourcesQueryArgs[2]] != nil {
		parsedParams.ResourceType = a[resourcesQueryArgs[2]].([]string)
	}
	returnCode, report, err := GetResourcesData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}