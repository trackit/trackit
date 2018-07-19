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

package routes

import (
	"encoding/json"
	"net/http"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
	"github.com/trackit/trackit2/aws"
)

// AwsAccountIdsOptionalQueryArg allows to get the AWS Account IDs in the URL
// Parameters with routes.RequiredQueryArgs. This AWS Account IDs will be a
// slice of Uint stored in the routes.Arguments map with itself for key.
// AwsAccountIdsOptionalQueryArg is optional and will not panic if no query
// argument is found.
var DetailedOptionalQueryArg = routes.QueryArg{
	Name:        "detailed",
	Type:        routes.QueryArgBool{},
	Description: "Detailed view.",
	Optional:    true,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAwsAccount).With(
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get aws accounts' data",
				Description: "Gets the data for all of the user's AWS accounts.",
			},
			routes.QueryArgs{
				routes.AwsAccountIdsOptionalQueryArg,
				DetailedOptionalQueryArg,
			},
	),
		http.MethodPost: routes.H(postAwsAccount).With(
		users.RequireAuthenticatedUser{users.ViewerCannot},
		routes.RequestContentType{"application/json"},
		routes.RequestBody{postAwsAccountRequestBody{
		RoleArn:  "arn:aws:iam::123456789012:role/example",
		External: "LlzrwHeiM-SGKRLPgaGbeucx_CJC@QBl,_vOEF@o",
		Pretty:   "My AWS account",
	}},
		routes.Documentation{
		Summary:     "add an aws account",
		Description: "Adds an AWS account to the user's list of accounts, validating it before succeeding.",
	},
	),
		http.MethodPatch: routes.H(patchAwsAccount).With(
		users.RequireAuthenticatedUser{users.ViewerCannot},
		routes.RequestContentType{"application/json"},
		routes.QueryArgs{routes.AwsAccountIdQueryArg},
		routes.Documentation{
		Summary:     "edit an aws account",
		Description: "Edits an AWS account from the user's list of accounts.",
	},
	),
		http.MethodDelete: routes.H(deleteAwsAccount).With(
		users.RequireAuthenticatedUser{users.ViewerCannot},
		routes.QueryArgs{routes.AwsAccountIdQueryArg},
		aws.RequireAwsAccountId{},
		routes.Documentation{
		Summary:     "delete an aws account",
		Description: "Delete the aws account passed in the query args.",
	},
	),
	}.H().With(
		db.RequestTransaction{db.Db},
		routes.Documentation{
			Summary: "interact with user's aws accounts",
		},
	).Register("/aws")
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(aws.NextExternal).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerCannot},
			routes.Documentation{
				Summary:     "get data to add next aws account",
				Description: "Gets data the user must have in order to successfully set up their account with the product.",
			},
		),
	}.H().Register("/aws/next")
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}
