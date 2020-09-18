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

	"github.com/trackit/jsonlog"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentmethod"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// routeGetStripeCustomerInformation returns the stripe customer information.
func routeGetStripeCustomerInformation(request *http.Request, a routes.Arguments) (int, interface{}) {
	isSubscribed := false
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

	if (dbUser.AwsCustomerEntitlement || dbUser.StripeCustomerEntitlement) {
		isSubscribed = true
	}
	if (dbUser.StripePaymentMethodIdentifier != "") {
		stripe.Key = "sk_test_51HAGBpHPvmk5HTchHxZJ0h9RGC1M8DAVuoSQeQJc3fbLFXII39hjAEMMxp0A7sAm5leXcZe8qsKbxHUXfBefjySO00ZJ4TSYeP"
		pm, err := paymentmethod.Get(dbUser.StripePaymentMethodIdentifier, nil)
		if err != nil {
			l.Error("Failed to get stripe customer payment method", err)
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, map[string]interface{}{
			"customerId": dbUser.StripeCustomerIdentifier,
			"subscriptionId": dbUser.StripeSubscriptionIdentifier,
			"paymentMethod": pm,
			"isSubscribed": isSubscribed,
		}
	}

	return http.StatusOK, map[string]interface{}{
		"customerId": dbUser.StripeCustomerIdentifier,
		"subscriptionId": dbUser.StripeSubscriptionIdentifier,
		"paymentMethod": dbUser.StripePaymentMethodIdentifier,
		"isSubscribed": isSubscribed,
	}
}
