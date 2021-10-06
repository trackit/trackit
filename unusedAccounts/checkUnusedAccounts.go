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
package unusedaccounts

import (
	"context"
	"time"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
)

const day = time.Hour * 24
const month = day * 30

const unusedThreshold = month
const deleteThreshold = unusedThreshold + month

var remindersThresholds = []time.Duration{unusedThreshold, unusedThreshold + deleteThreshold - day*7, unusedThreshold + deleteThreshold - day*3, unusedThreshold + deleteThreshold - day*1}

// CheckUnusedAccounts checks for unused accounts, sends reminders and delete unused data
func CheckUnusedAccounts(ctx context.Context) error {
	users, err := models.GetUnusedAccounts(db.Db, unusedThreshold)

	if err != nil {
		return err
	}

	for _, user := range users {
		if user == nil || len((*user).AwsCustomerIdentifier) > 0 {
			continue
		}
		if checkErr := checkUnusedAccount(ctx, *user); err == nil {
			err = checkErr
		}
	}

	return nil
}

func checkUnusedAccount(ctx context.Context, user models.User) error {
	unusedTime := time.Since(user.LastSeen)

	thresholdStage := 0
	for i, reminderThreshold := range remindersThresholds {
		if unusedTime > reminderThreshold {
			thresholdStage = i
		} else {
			break
		}
	}

	var err error = nil

	if user.LastUnusedReminder.Sub(user.LastSeen) < remindersThresholds[thresholdStage] {
		timeBeforeDeletion := time.Until(user.LastSeen.Add(deleteThreshold))
		err = sendReminder(ctx, user, timeBeforeDeletion)

		user.LastUnusedReminder = time.Now()
		if updateErr := user.Update(db.Db); err == nil {
			err = updateErr
		}
	}

	return err
}
