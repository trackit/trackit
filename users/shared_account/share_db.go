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
	"context"
	"database/sql"
	"errors"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/models"
)

type SharedResults struct {
	ShareId       int    `json:"sharedId"      req:"nonzero"`
	Mail          string `json:"email"         req:"nonzero"`
	Level         int    `json:"level"`
	UserId        int    `json:"userId"        req:"nonzero"`
	SharingStatus bool   `json:"sharingStatus"`
}

// GetSharingList returns a list of users who have access to a specific AWS account
func GetSharingList(ctx context.Context, db models.DB, accountId int) ([]SharedResults, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	response := []SharedResults{}
	dbSharedAccounts, err := models.SharedAccountByAccountID(db, accountId)
	if err == sql.ErrNoRows {
		return response, nil
	} else if err != nil {
		logger.Error("Error getting shared account from database.", err.Error())
		return nil, errors.New("Error while getting data from database")
	} else {
		for _, key := range dbSharedAccounts {
			dbUser, err := models.UserByID(db, key.UserID)
			if err != nil {
				logger.Error("Error getting users from database.", err.Error())
				return nil, errors.New("Error while getting data from database")
			}
			response = append(response, SharedResults{key.ID, dbUser.Email, key.UserPermission, key.UserID, key.SharingAccepted})
		}
		return response, nil
	}
}

// UpdateSharedUser updates user permission level
func UpdateSharedUser(ctx context.Context, db models.DB, shareId int, permissionLevel int) (interface{}, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbSharedAccount, err := models.SharedAccountByID(db, shareId)
	if err != nil {
		logger.Error("Error while getting shared user information", err)
		return nil, err
	}
	dbSharedAccount.UserPermission = permissionLevel
	err = dbSharedAccount.Update(db)
	if err != nil {
		logger.Error("Error while updating user permission", err)
		return nil, err
	}
	return dbSharedAccount, nil
}

// DeleteSharedUser deletes a user access to an AWS account by removing entry in shared_account database table.
func DeleteSharedUser(ctx context.Context, db models.DB, shareId int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbSharedAccount, err := models.SharedAccountByID(db, shareId)
	if err != nil {
		logger.Error("Error while getting shared user information", err)
		return err
	}
	err = dbSharedAccount.Delete(db)
	if err != nil {
		logger.Error("Error while deleting shared user", err)
		return err
	}
	return nil
}
