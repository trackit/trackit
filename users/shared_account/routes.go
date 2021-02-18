package shared_account

import (
	"encoding/json"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
	"net/http"
)

// inviteUserRequest is the expected request body for the invite user route handler.
type InviteUserRequest struct {
	Email           string `json:"email" req:"nonzero"`
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
