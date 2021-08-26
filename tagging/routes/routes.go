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

	"github.com/stripe/stripe-go/v72"

	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// taggingComplianceQueryArgs allows to get required queryArgs params
var taggingComplianceQueryArgs = []routes.QueryArg{
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
}

// suggestionsQueryArgs allows to get required queryArgs params
var suggestionsQueryArgs = []routes.QueryArg{
	routes.QueryArg{
		Name:        "tagkey",
		Type:        routes.QueryArgString{},
		Description: "Tag key for suggestions",
		Optional:    false,
	},
}

func init() {

	stripe.Key = config.StripeKey

	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetMostUsedTags).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get most used tags",
				Description: "Responds with most used tags for a user.",
			},
		),
	}.H().Register("/tagging/mostusedtags")

	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetTaggingCompliance).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(taggingComplianceQueryArgs),
			routes.Documentation{
				Summary:     "get tagging compliance",
				Description: "Responds with tagging compliance data in a specified range",
			},
		),
	}.H().Register("/tagging/compliance")

	routes.MethodMuxer{
		http.MethodPost: routes.H(routeGetResources).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{ResourcesRequestBody{[]string{"394125495069"}, []string{"us-west-2"}, []string{"lambda"}, []Tag{{"project", "trackit"}}, []Tag{{"Product", "msol"}}}},
			routes.Documentation{
				Summary:     "get list of resources",
				Description: "Responds with the list of resources based on the request body passed to it",
			},
		),
	}.H().Register("/tagging/resources")

	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetTaggingSuggestions).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(suggestionsQueryArgs),
			routes.Documentation{
				Summary:     "get suggestions for a tag's value",
				Description: "Responds with suggestions for a tag's value for a user.",
			},
		),
	}.H().Register("/tagging/suggestions/tag-value")

	routes.MethodMuxer{
		http.MethodGet: routes.H(shouldPopup).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get Tagbot access",
				Description: "Returns whether or not to display subscription popup",
			},
		),
	}.H().Register("/tagging/should-popup")
	routes.MethodMuxer{
		http.MethodPost: routes.H(routeCreateCustomer).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{CreateCustomerRequestBody{"example@example.com"}},
			routes.Documentation{
				Summary:     "Create a stripe customer",
				Description: "Responds with customer information",
			},
		),
	}.H().Register("/tagging/create-customer")

	routes.MethodMuxer{
		http.MethodPost: routes.H(routeCreateSubscription).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{CreateSubscriptionRequestBody{"pm_1HL9NKHPvmk5HTchutPljt1d", "cus_HuzN2Ie7ZFLvHC", "tagbot"}},
			routes.Documentation{
				Summary:     "Create stripe payment method",
				Description: "Responds with payment method information",
			},
		),
	}.H().Register("/tagging/create-subscription")

	routes.MethodMuxer{
		http.MethodPost: routes.H(routeHandleRetryInvoice).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{RetryInvoiceRequestBody{"cus_HuzN2Ie7ZFLvHC", "pm_1HL9NKHPvmk5HTchutPljt1d", "in_1HOj7bHPvmk5HTchOPvXsNry"}},
			routes.Documentation{
				Summary:     "Handle retry invoice",
				Description: "Updates stripe customer with new payment method",
			},
		),
	}.H().Register("/tagging/retry-invoice")

	routes.MethodMuxer{
		http.MethodPost: routes.H(routeCancelSubscription).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{CancelSubscriptionRequestBody{"sub_HyiV5QLQrY1YIc"}},
			routes.Documentation{
				Summary:     "Cancel subscription",
				Description: "Cancels customer subscription",
			},
		),
	}.H().Register("/tagging/cancel-subscription")

	routes.MethodMuxer{
		http.MethodPost: routes.H(routeRetrieveSubscription).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{RetrieveSubscriptionRequestBody{"sub_HyiV5QLQrY1YIc"}},
			routes.Documentation{
				Summary:     "Retrieve subscription",
				Description: "Retrieves customer subscription information",
			},
		),
	}.H().Register("/tagging/retrieve-subscription")

	routes.MethodMuxer{
		http.MethodPost: routes.H(routeChangePaymentMethod).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{ChangePaymentMehtodRequestBody{"pm_1HL9NKHPvmk5HTchutPljt1d"}},
			routes.Documentation{
				Summary:     "Change payment method",
				Description: "Changes customer payment method",
			},
		),
	}.H().Register("/tagging/change-payment-method")

	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetStripeCustomerInformation).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get stripe customer information",
				Description: "Returns stripe customer information",
			},
		),
	}.H().Register("/tagging/stripe-customer-information")
}
