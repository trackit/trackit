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

package shared_account

import (
	"database/sql"
	"net/http"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

// inviteUser handles users invite for team sharing.
func inviteUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body InviteUserRequest
	accountId := a[routes.AwsAccountIdQueryArg].(int)
	err := decodeRequestBody(request, &body)
	if err != nil {
		return http.StatusBadRequest, errors.GetErrorMessage(request.Context(), err)
	}
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	return InviteUserWithValidBody(request, body, accountId, tx, user)
}

// listSharedUsers handles listing of users who have an access to an AWS account.
func listSharedUsers(request *http.Request, a routes.Arguments) (int, interface{}) {
	accountId := a[routes.AwsAccountIdQueryArg].(int)
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	return listSharedUserAccessWithValidBody(request, accountId, tx, user)
}

// updateSharedUsers handles updates of user permission level for team sharing.
func updateSharedUsers(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body updateUsersSharedAccountRequest
	shareId := a[routes.ShareIdQueryArg].(int)
	err := decodeRequestBody(request, &body)
	if err != nil {
		return http.StatusBadRequest, errors.GetErrorMessage(request.Context(), err)
	}
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	return updateSharedUserAccessWithValidBody(request, body, shareId, tx, user)
}

// deleteSharedUsers handles user access deletion for team sharing
func deleteSharedUsers(request *http.Request, a routes.Arguments) (int, interface{}) {
	shareId := a[routes.ShareIdQueryArg].(int)
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	return deleteSharedUserAccessWithValidBody(request, shareId, tx, user)
}

// listSharedUserAccessWithValidBody tries to list users who have an access to an AWS account
func listSharedUserAccessWithValidBody(request *http.Request, accountId int, tx *sql.Tx, user users.User) (int, interface{}) {
	ctx := request.Context()
	security, err := safetyCheckByAccountId(ctx, tx, accountId, user)
	if err != nil {
		return http.StatusBadRequest, err
	} else if !security {
		return http.StatusForbidden, errors.GetErrorMessage(ctx, &errors.SharedAccountError{errors.SharedAccountNoPermission, "You do not have permission to view users of this account"})
	}
	res, err := GetSharingList(request.Context(), db.Db, accountId)
	if err != nil {
		return http.StatusForbidden, errors.GetErrorMessage(ctx, &errors.SharedAccountError{errors.SharedAccountRequestError, "Error retrieving shared users list"})
	} else {
		return http.StatusOK, res
	}
}

// updateSharedUserAccessWithValidBody tries to update users permission level for team sharing
func updateSharedUserAccessWithValidBody(request *http.Request, body updateUsersSharedAccountRequest, shareId int, tx *sql.Tx, user users.User) (int, interface{}) {
	ctx := request.Context()
	security, err := safetyCheckByShareIdAndPermissionLevel(ctx, tx, shareId, body.PermissionLevel, user)
	if err != nil {
		return http.StatusBadRequest, err
	} else if !security {
		return http.StatusForbidden, errors.GetErrorMessage(ctx, &errors.SharedAccountError{errors.SharedAccountNoPermission, "You do not have permission to edit this sharing"})
	}
	if !checkPermissionLevel(body.PermissionLevel) {
		return http.StatusBadRequest, errors.GetErrorMessage(ctx, &errors.SharedAccountError{errors.SharedAccountBadPermission, "Bad permission level"})
	}
	res, err := UpdateSharedUser(request.Context(), db.Db, shareId, body.PermissionLevel)
	if err != nil {
		return http.StatusForbidden, errors.GetErrorMessage(ctx, &errors.SharedAccountError{errors.SharedAccountRequestError, "Error retrieving shared user list"})
	}
	return http.StatusOK, res
}

// deleteSharedUserAccessWithValidBody tries to delete users from accessing specific shared aws account
func deleteSharedUserAccessWithValidBody(request *http.Request, shareId int, tx *sql.Tx, user users.User) (int, interface{}) {
	ctx := request.Context()
	security, err := safetyCheckByShareId(ctx, tx, shareId, user)
	if err != nil {
		return http.StatusBadRequest, err
	} else if !security {
		return http.StatusForbidden, errors.GetErrorMessage(ctx, &errors.SharedAccountError{errors.SharedAccountNoPermission, "You do not have permission to delete this sharing"})
	}
	err = DeleteSharedUser(request.Context(), db.Db, shareId)
	if err != nil {
		return http.StatusBadRequest, errors.GetErrorMessage(ctx, &errors.SharedAccountError{errors.SharedAccountRequestError, "Error deleting shared user"})
	}
	return http.StatusOK, nil
}
