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

// Package routes adds a bunch of routes that interact with a user's AWS accounts/sub-accounts
package routes

import (
	"encoding/json"
	"net/http"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

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
		http.MethodPatch: routes.H(patchAwsSubaccount).With(
			users.RequireAuthenticatedUser{users.ViewerCannot},
			routes.RequestContentType{"application/json"},
			routes.QueryArgs{routes.AwsAccountIdQueryArg},
			routes.RequestBody{patchAwsSubaccountRequestBody{
				RoleArn:  "arn:aws:iam::123456789012:role/example",
				External: "LlzrwHeiM-SGKRLPgaGbeucx_CJC@QBl,_vOEF@o",
			}},
			routes.Documentation{
				Summary:     "link a role to a subaccount",
				Description: "Edits an AWS subaccount from the user's list of accounts.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
		routes.Documentation{
			Summary: "interact with user's aws subaccounts",
		},
	).Register("/aws/subaccount")
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

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAwsAccountsStatus).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get status of aws accounts",
				Description: "Gets status of AWS Accounts and their bill repositories.",
			},
		),
	}.H().Register("/aws/status")
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}
