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

package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// PopupInfoResponseBody is the response body in case /tagging/should-popup is called.
type PopupInfoResponseBody struct {
	Popup bool `json:"popup"`
}

var (
	tagbotFreeTrialDuration = time.Hour * 24 * 14
)

// shouldPopup verify if the user has access to Tagbot.
func shouldPopup(request *http.Request, a routes.Arguments) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(request.Context())
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	dbUser, err := models.TagbotUserByUserID(tx, user.Id)
	if err != nil {
		l.Error("Failed to get tagbot user with id", map[string]interface{}{
			"userId": user.Id,
			"error":  err.Error(),
		})
		return http.StatusInternalServerError, errors.New("Failed to get Tagbot user with id")
	}
	customer, err := models.UserByID(db.Db, user.Id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error while getting customer infos")
	}

	return checkPopup(dbUser, customer)
}

func checkPopup(dbUser *models.TagbotUser, customer *models.User) (int, interface{}) {
	if dbUser.AwsCustomerEntitlement {
		return http.StatusOK, PopupInfoResponseBody{
			false,
		}
	}
	if dbUser.StripeCustomerEntitlement {
		return http.StatusOK, PopupInfoResponseBody{
			false,
		}
	}
	if checkUserTagbotFreeTrial(customer.Created) {
		return http.StatusOK, PopupInfoResponseBody{
			false,
		}
	}
	return http.StatusOK, PopupInfoResponseBody{
		true,
	}
}

// checkUserTagbotFreeTrial returns whether the time since the creation date exceeds the duration of a TagBot free trial
func checkUserTagbotFreeTrial(creationDate time.Time) bool {
	return time.Since(creationDate) <= tagbotFreeTrialDuration
}
