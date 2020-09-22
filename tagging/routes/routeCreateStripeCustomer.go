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
	"database/sql"
	"errors"
	"net/http"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type CreateCustomerRequestBody struct {
	Email string `json:"email" req:"nonzero"`
}

func routeCreateCustomer(request *http.Request, a routes.Arguments) (int, interface{}) {
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
	var body CreateCustomerRequestBody
	routes.MustRequestBody(a, &body)
	params := &stripe.CustomerParams{
		Email: stripe.String(body.Email),
	}
	c, err := customer.New(params)
	if err != nil {
		l.Error("Failed to create stripe customer", err)
		return http.StatusInternalServerError, err
	}

	dbUser.StripeCustomerIdentifier = c.ID
	err = dbUser.Update(tx)
	if err != nil {
		l.Error("Failed to update Tagbot Customer ID", err)
		return http.StatusInternalServerError, errors.New("Failed to update Tagbot customer ID")
	}

	res := struct {
		Customer *stripe.Customer `json:"customer"`
	}{
		Customer: c,
	}
	return http.StatusOK, res
}
