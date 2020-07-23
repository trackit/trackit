package unusedaccounts

import (
	"context"
	"fmt"
	"time"

	"github.com/trackit/trackit/mail"
	"github.com/trackit/trackit/models"
)

func sendRemainder(ctx context.Context, user models.User, timeBeforeDeletion time.Duration) error {
	daysBeforeDeletion := timeBeforeDeletion / (time.Hour * 24)
	if daysBeforeDeletion < 1 {
		daysBeforeDeletion = 1
	}

	body := fmt.Sprintf("Your TrackIt account is not used anymore. Please login again or your data will be deleted in %d days.", daysBeforeDeletion)
	return mail.SendMail(user.Email, "Your TrackIt account is not used anymore.", body, ctx)
}
