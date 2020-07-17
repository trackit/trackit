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

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetMostUsedTags).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get most used tags",
				Description: "Responds with most used tags for an AWS account.",
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
