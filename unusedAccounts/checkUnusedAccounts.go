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
	unusedTime := time.Now().Sub(user.LastSeen)

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
		timeBeforeDeletion := user.LastSeen.Add(deleteThreshold).Sub(time.Now())
		err = sendReminder(ctx, user, timeBeforeDeletion)

		user.LastUnusedReminder = time.Now()
		if updateErr := user.Update(db.Db); err == nil {
			err = updateErr
		}
	}

	return err
}
