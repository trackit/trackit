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

	"github.com/trackit/jsonlog"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/sub"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
    "github.com/trackit/trackit/routes"
    "github.com/trackit/trackit/users"
)

type CreateSubscriptionRequestBody struct {
	PaymentMethodID string `json:"paymentMethodId" req:"nonzero"`
	CustomerID      string `json:"customerId"      req:"nonzero"`
	PriceID         string `json:"priceId"         req:"nonzero"`
}

func routeCreateSubscription(request *http.Request, a routes.Arguments) (int, interface{}) {
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
	stripe.Key = "sk_test_51HAGBpHPvmk5HTchHxZJ0h9RGC1M8DAVuoSQeQJc3fbLFXII39hjAEMMxp0A7sAm5leXcZe8qsKbxHUXfBefjySO00ZJ4TSYeP"
	var body CreateSubscriptionRequestBody
	routes.MustRequestBody(a, &body)

	// Attach PaymentMethod
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(body.CustomerID),
	}
	pm, err := paymentmethod.Attach(
		body.PaymentMethodID,
		params,
	)
	if err != nil {
		l.Error("Stripe failed to attach payment method", err)
		return http.StatusInternalServerError, err
	}
	dbUser.StripePaymentMethodIdentifier = pm.ID
	err = dbUser.Update(tx)
	if err != nil {
        l.Error("Failed to update Tagbot stripe payment method ID.", err)
        return http.StatusInternalServerError, errors.New("Failed to update Tagbot stripe payment method ID.")
    }

	// Update invoice settings default
	customerParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
		DefaultPaymentMethod: stripe.String(pm.ID),
		},
	}
	c, err := customer.Update(
		body.CustomerID,
		customerParams,
	)
    if err != nil {
		l.Error("Stripe failed to update customer invoice settings", c.ID)
        return http.StatusInternalServerError, err
	}

	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(body.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
		{
			Plan: stripe.String("price_1HIaP7HPvmk5HTchKgWfjBwe"),
		},
		},
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	s, err := sub.New(subscriptionParams)
    if err != nil {
		l.Error("Stripe failed to create subscritpion", err)
        return http.StatusInternalServerError, err
	}
	if (s.Status == "active") {
		dbUser.StripeCustomerEntitlement = true
	}
	dbUser.StripeSubscriptionIdentifier = s.ID
	err = dbUser.Update(tx)
	if err != nil {
		l.Error("Failed to update Tagbot subscription ID and Tagbot customer entitlement", err)
		return http.StatusInternalServerError, errors.New("Failed to update Tagbot subscription ID and Tagbot customer entitlement")
	}
	return http.StatusOK, s
}
