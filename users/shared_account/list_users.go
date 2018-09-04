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
	"github.com/trackit/trackit-server/models"
	"context"
	"github.com/trackit/jsonlog"
	"fmt"
)

type SharedResults struct {
	Mail string
	Level int
	Uid int
	SharingStatus bool
}

func GetSharingList(ctx context.Context, db models.XODB, accountId int) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var UserStruct = make(map[int]SharedResults)
	dbSharedAccounts, err := models.SharedAccountsByAccountID(db, accountId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		logger.Error("Error getting shared account from database.", err.Error())
		return err
	} else {
		for _, key := range dbSharedAccounts {
			dbUser, err := models.UserByID(db, key.UserID)
			if err != nil {
				return err
			}
			var s = SharedResults{Mail: dbUser.Email, Level: key.UserPermission, Uid: key.UserID, SharingStatus: key.SharingAccepted}
			UserStruct[key.UserID] = s
		}
		fmt.Print("--- TABLE : ")
		fmt.Print(UserStruct)
		return nil
	}
}