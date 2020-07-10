package unusedaccounts

import (
	"time"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
)

const unusedDelay = time.Hour

// CheckUnusedAccounts checks for unused accounts, sends reminders and delete unused data
func CheckUnusedAccounts() error {
	users, err := models.GetUnusedAccounts(db.Db, unusedDelay)

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
	// Process unused account
	return nil
}
