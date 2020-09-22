//   Copyright 2020 MSolution.IO
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
	"net/http"
	"database/sql"
	"errors"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type RetrieveSubscriptionRequestBody struct {
	SubscriptionID string `json:"subscriptionId" req:"nonzero"`
}

func routeRetrieveSubscription(request *http.Request, a routes.Arguments) (int, interface{}) {
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
	stripe.Key = config.StripeKey
	var body RetrieveSubscriptionRequestBody
	routes.MustRequestBody(a, &body)

	s, err := sub.Get(body.SubscriptionID, nil)
	if err != nil {
		l.Error("Stripe failed to retrieve customer subscription", err)
		return http.StatusInternalServerError, err
	}

	if (s.Status == "active") {
		dbUser.StripeSubscriptionIdentifier = s.ID
		dbUser.StripeCustomerEntitlement = true
	}
	err = dbUser.Update(tx)
	if err != nil {
		l.Error("Failed to update Tagbot subscription ID and Stripe Customer Entitlement", err)
		return http.StatusInternalServerError, errors.New("Failed to update Tagbot subscription ID and stripe payment method ID")
	}
	return http.StatusOK, s
}
