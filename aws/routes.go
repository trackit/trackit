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

package aws

import (
	"encoding/json"
	"net/http"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAwsAccount).With(
			routes.Documentation{
				Summary:     "get aws accounts' data",
				Description: "Gets the data for all of the user's AWS accounts.",
			},
		),
		http.MethodPost: routes.H(postAwsAccount).With(
			routes.RequestContentType{"application/json"},
			routes.Documentation{
				Summary:     "add an aws account",
				Description: "Adds an AWS account to the user's list of accounts, validating it before succeeding.",
			},
		),
	}.H().With(
		users.RequireAuthenticatedUser{},
		db.RequestTransaction{},
		routes.Documentation{
			Summary: "interact with user's aws accounts",
		},
	).Register("/aws")
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(nextExternal).With(
			db.RequestTransaction{},
			users.RequireAuthenticatedUser{},
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
