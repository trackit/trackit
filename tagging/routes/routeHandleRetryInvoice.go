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
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type RetryInvoiceRequestBody struct {
	CustomerID      string `json:"customerId"      req:"nonzero"`
	PaymentMethodID string `json:"paymentMethodId" req:"nonzero"`
	InvoiceID       string `json:"invoiceId"       req:"nonzero"`
}

func routeHandleRetryInvoice(request *http.Request, a routes.Arguments) (int, interface{}) {
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
	var body RetryInvoiceRequestBody
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

	// Retrieve Invoice
	invoiceParams := &stripe.InvoiceParams{}
	invoiceParams.AddExpand("payment_intent")
	in, err := invoice.Get(
		body.InvoiceID,
		invoiceParams,
	)

	if err != nil {
		l.Error("Stripe failed to retrieve invoice", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, in
}
