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

package shared_account

import (
	"database/sql"
	"net/http"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/users"
)

// inviteUserRequest is the expected request body for the invite user route handler.
type InviteUserRequest struct {
	Email           string `json:"email" req:"nonzero"`
	AccountId       int    `json:"accountId"`
	PermissionLevel int    `json:"permissionLevel"`
}

type listUserSharedAccountRequest struct {
	AccountId       int    `json:"accountId" req:"nonzero"`
}

type updateUsersSharedAccountRequest struct {
	ShareId         int `json:"shareId" req:"nonzero"`
	PermissionLevel int `json:"permissionLevel"`
}

type deleteUsersSharedAccountRequest struct {
	ShareId         int `json:"shareId" req:"nonzero"`
}
//
//type sharedAccount struct {
//	AccountId       int
//	userId          int
//	UserPermission  int
//	SharingAccepted int
//}

func init() {
	routes.MethodMuxer{
		http.MethodPost: routes.H(inviteUser).With(
			routes.RequestContentType{"application/json"},
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestBody{InviteUserRequest{"example@example.com", 1234, 0}},
			routes.Documentation{
				Summary:     "Creates an invite",
				Description: "Creates an invite for account team sharing",
			},
		),
	}.H().Register("/user/share/add")
	routes.MethodMuxer{
		http.MethodPost: routes.H(listUsers).With(
			routes.RequestContentType{"application/json"},
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestBody{listUserSharedAccountRequest{1}},
			routes.Documentation{
				Summary:     "List shared users",
				Description: "Return a list of user who have an access to an AWS account on Trackit",
			},
		),
	}.H().Register("/user/share/list")
	routes.MethodMuxer{
		http.MethodPost: routes.H(updateUsers).With(
			routes.RequestContentType{"application/json"},
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestBody{updateUsersSharedAccountRequest{1, 2}},
			routes.Documentation{
				Summary:     "Update shared users",
				Description: "Update shared users associated with a specific AWS account",
			},
		),
	}.H().Register("/user/share/update")
	routes.MethodMuxer{
		http.MethodPost: routes.H(deleteUsers).With(
			routes.RequestContentType{"application/json"},
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestBody{deleteUsersSharedAccountRequest{1}},
			routes.Documentation{
				Summary:     "Delete shared users",
				Description: "Delete shared users associated with a specific AWS account",
			},
		),
	}.H().Register("/user/share/delete")
}

// inviteUser handles users invite for team sharing.
func inviteUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body InviteUserRequest
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	return InviteUserWithValidBody(request, body, tx, user)
}

// listUsers return the list of user who have an access to a specific AWS account
func listUsers(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body listUserSharedAccountRequest
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	return listUserWithValidBody(request, body, tx, user)
}

// listUsers return the list of user who have an access to a specific AWS account
func updateUsers(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body updateUsersSharedAccountRequest
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	return updateUserWithValidBody(request, body, tx)
}

// listUsers return the list of user who have an access to a specific AWS account
func deleteUsers(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body deleteUsersSharedAccountRequest
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	return deleteUserWithValidBody(request, body, tx)
}

// listUsersWithValidBody return the list of user who have an access to a specific AWS account
func listUserWithValidBody(request *http.Request, body listUserSharedAccountRequest, tx *sql.Tx, user users.User) (int, interface{}) {
	check, err := checkUserOwnAccount(request.Context(), db.Db, body.AccountId, user.Id)
	if err != nil {
		return 403, "Error retrieving shared users list"
	}
	if check {
		res, err := GetSharingList(request.Context(), db.Db, body.AccountId)
		if err != nil {
			return 403, "Error retrieving shared users list"
		} else {
			return 200, res
		}
	} else {
		return 403, "You do not own this account"
	}
}

// listUsersWithValidBody return the list of user who have an access to a specific AWS account
func updateUserWithValidBody(request *http.Request, body updateUsersSharedAccountRequest, tx *sql.Tx) (int, interface{}) {
	err := UpdateSharedUser(request.Context(), db.Db, body.ShareId, body.PermissionLevel)
	if err != nil {
		return 403, "Error updating shared user list"
	}
	return 200, "ok"
}

// listUsersWithValidBody return the list of user who have an access to a specific AWS account
func deleteUserWithValidBody(request *http.Request, body deleteUsersSharedAccountRequest, tx *sql.Tx) (int, interface{}) {
	err := DeleteSharedUser(request.Context(), db.Db, body.ShareId)
	if err != nil {
		return 403, "Error deleting shared user"
	}
	return 200, "ok"
}
