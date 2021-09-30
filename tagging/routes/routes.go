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

package routes

import (
	"net/http"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// taggingComplianceQueryArgs allows to get required queryArgs params
var taggingComplianceQueryArgs = []routes.QueryArg{
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
}

// resourcesQueryArgs allows to get required queryArgs params
var resourcesQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.QueryArg{
		Name:        "regions",
		Type:        routes.QueryArgStringSlice{},
		Description: "Regions of the resources",
		Optional:    true,
	},
	routes.QueryArg{
		Name:        "resourceTypes",
		Type:        routes.QueryArgStringSlice{},
		Description: "Types of the resources",
		Optional:    true,
	},
}

// suggestionsQueryArgs allows to get required queryArgs params
var suggestionsQueryArgs = []routes.QueryArg{
	routes.QueryArg{
		Name:        "tagkey",
		Type:        routes.QueryArgString{},
		Description: "Tag key for suggestions",
		Optional:    false,
	},
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetMostUsedTags).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get most used tags",
				Description: "Responds with most used tags for a user.",
			},
		),
	}.H().Register("/tagging/mostusedtags")
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetTaggingCompliance).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(taggingComplianceQueryArgs),
			routes.Documentation{
				Summary:     "get tagging compliance",
				Description: "Responds with tagging compliance data in a specified range",
			},
		),
	}.H().Register("/tagging/compliance")
}

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

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetTaggingSuggestions).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(suggestionsQueryArgs),
			routes.Documentation{
				Summary:     "get suggestions for a tag's value",
				Description: "Responds with suggestions for a tag's value for a user.",
			},
		),
	}.H().Register("/tagging/suggestions/tag-value")
}
