package unusedaccounts

import (
	"time"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
)

const unusedThreshold = time.Second * 15
const deleteThreshold = unusedThreshold + time.Second*45

var remaindersThresholds = []time.Duration{unusedThreshold, unusedThreshold + time.Second*15, unusedThreshold + time.Second*30}

// CheckUnusedAccounts checks for unused accounts, sends reminders and delete unused data
func CheckUnusedAccounts() error {
	users, err := models.GetUnusedAccounts(db.Db, unusedThreshold)

	if err != nil {
		return err
	}

	for _, user := range users {
		if user != nil {
			checkUnusedAccount(*user)
		}
	}

	return nil
}

func checkUnusedAccount(user models.User) error {
	unusedTime := time.Now().Sub(user.LastSeen)

	if unusedTime > deleteThreshold {
		return deleteData(user)
	}

	thresholdStage := 0
	for i, remainderThreshold := range remaindersThresholds {
		if unusedTime > remainderThreshold {
			thresholdStage = i
		} else {
			break
		}
	}

	var err error = nil

	if user.LastUnusedReminder.Sub(user.LastSeen) < remaindersThresholds[thresholdStage] {
		err = sendRemainder(user)

		user.LastUnusedReminder = time.Now()
		user.Update(db.Db)
	}

	return err
}
