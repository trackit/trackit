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
	"fmt"
	"time"

	"github.com/trackit/trackit/mail"
	"github.com/trackit/trackit/models"
)

func sendReminder(ctx context.Context, user models.User, timeBeforeDeletion time.Duration) error {
	daysBeforeDeletion := timeBeforeDeletion / (time.Hour * 24)
	if daysBeforeDeletion < 1 {
		daysBeforeDeletion = 1
	}

	body := fmt.Sprintf("Your TrackIt account is not used anymore. Please login again or your data will be deleted in %d days.", daysBeforeDeletion)
	return mail.SendMail(user.Email, "Your TrackIt account is not used anymore.", body, ctx)
}
