package unusedaccounts

import (
	"context"
	"fmt"
	"time"

	"github.com/trackit/trackit/mail"
	"github.com/trackit/trackit/models"
)

func sendRemainder(ctx context.Context, user models.User, timeBeforeDeletion time.Duration) error {
	body := fmt.Sprintf("Your TrackIt account is not used anymore. Please login again or your data will be deleted in %s.", timeBeforeDeletion.String())
	return mail.SendMail(user.Email, "Your TrackIt account is not used anymore.", body, ctx)
}
