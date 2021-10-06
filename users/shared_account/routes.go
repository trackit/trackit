//   Copyright 2021 MSolution.IO
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
package shared_account

import (
	"encoding/json"
	"net/http"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// InviteUserRequest is the expected request body for the invite user route handler.
type InviteUserRequest struct {
	Email           string `json:"email"            req:"nonzero"`
	Origin          string `json:"origin"           req:"nonzero"`
	PermissionLevel int    `json:"permissionLevel"`
}

type updateUsersSharedAccountRequest struct {
	PermissionLevel int `json:"permissionLevel"`
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(listSharedUsers).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "List shared users",
				Description: "Return a list of user who have an access to an AWS account on Trackit",
			},
			routes.QueryArgs{
				routes.AwsAccountIdQueryArg,
			},
		),
		http.MethodPost: routes.H(inviteUser).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.Documentation{
				Summary:     "Creates an invite",
				Description: "Creates an invite for account team sharing. Permission level can be 0 for admin, 1 for standard and 2 for read-only.",
			},
			routes.QueryArgs{
				routes.AwsAccountIdQueryArg,
			},
		),
		http.MethodPatch: routes.H(updateSharedUsers).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs{
				routes.ShareIdQueryArg,
			},
			routes.RequestContentType{"application/json"},
			routes.Documentation{
				Summary:     "Update shared users",
				Description: "Update shared users associated with a specific AWS account. Permission level can be 0 for admin, 1 for standard and 2 for read-only.",
			},
		),
		http.MethodDelete: routes.H(deleteSharedUsers).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "Delete shared users",
				Description: "Delete shared users associated with a specific AWS account",
			},
			routes.QueryArgs{
				routes.ShareIdQueryArg,
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
		routes.Documentation{
			Summary: "interact with shared accounts",
		},
	).Register("/user/share")
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}
