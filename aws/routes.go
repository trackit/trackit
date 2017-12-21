//   Copyright 2017 MSolution.IO
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

package aws

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAwsAccount).With(
			routes.Documentation{
				Summary:     "get aws accounts' data",
				Description: "Gets the data for all of the user's AWS accounts.",
			},
			routes.QueryArgs{AwsAccountsOptionalQueryArg},
		),
		http.MethodPost: routes.H(postAwsAccount).With(
			routes.RequestContentType{"application/json"},
			routes.RequestBody{postAwsAccountRequestBody{
				RoleArn:  "arn:aws:iam::123456789012:role/example",
				External: "LlzrwHeiM-SGKRLPgaGbeucx_CJC@QBl,_vOEF@o",
				Pretty:   "My AWS account",
			}},
			routes.Documentation{
				Summary:     "add an aws account",
				Description: "Adds an AWS account to the user's list of accounts, validating it before succeeding.",
			},
		),
		http.MethodPatch: routes.H(patchAwsAccount).With(
			routes.RequestContentType{"application/json"},
			routes.QueryArgs{AwsAccountQueryArg},
			routes.Documentation{
				Summary:     "edit an aws account",
				Description: "Edits an AWS account from the user's list of accounts.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
		users.RequireAuthenticatedUser{},
		routes.Documentation{
			Summary: "interact with user's aws accounts",
		},
	).Register("/aws")
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(nextExternal).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{},
			routes.Documentation{
				Summary:     "get data to add next aws account",
				Description: "Gets data the user must have in order to successfully set up their account with the product.",
			},
		),
	}.H().Register("/aws/next")
}

var (
	// AwsAccountsOptionalQueryArg allows to get the AWS Account IDs in the URL
	// Parameters with routes.RequiredQueryArgs. This AWS Account IDs will be a
	// slice of Uint stored in the routes.Arguments map with itself for key.
	// AwsAccountsOptionalQueryArg is optional and will not panic if no query
	// argument is found.
	AwsAccountsOptionalQueryArg = routes.QueryArg{
		Name:        "aa",
		Type:        routes.QueryArgIntSlice{},
		Description: "The IDs for many AWS account.",
		Optional:    true,
	}

	// AwsAccountQueryArg allows to get the AWS Account ID in the URL Parameters
	// with routes.RequiredQueryArgs. This AWS Account ID will be an Uint stored
	// in the routes.Arguments map with itself for key.
	AwsAccountQueryArg = routes.QueryArg{
		Name:        "aa",
		Type:        routes.QueryArgInt{},
		Description: "The ID for an AWS account.",
	}
)

// RequireAwsAccount decorates handler to require that an AwsAccount be
// selected using RequiredQueryArgs{AwsAccountQueryArg}. The decorator will
// panic if no AwsAccountQueryArg query argument is found.
type RequireAwsAccount struct{}
type routeArgKey uint

const (
	AwsAccountSelection = routeArgKey(iota)
)

func (d RequireAwsAccount) Decorate(h routes.Handler) routes.Handler {
	h.Func = d.getFunc(h.Func)
	return h
}

func (_ RequireAwsAccount) getFunc(hf routes.HandlerFunc) routes.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		user, tx, err := getUserAndTransactionFromArguments(a)
		if err != nil {
			l.Error("Missing transaction or user for handler with AWS account.", err.Error())
			return http.StatusInternalServerError, nil
		}
		aaid := a[AwsAccountQueryArg].(int)
		aa, err := GetAwsAccountWithIdFromUser(user, aaid, tx)
		if err != nil {
			return http.StatusNotFound, errors.New("AWS account not found.")
		} else {
			a[AwsAccountSelection] = aa
			return hf(w, r, a)
		}

	}
}

func getUserAndTransactionFromArguments(a routes.Arguments) (users.User, *sql.Tx, error) {
	u := a[users.AuthenticatedUser]
	t := a[db.Transaction]
	if ut, ok := u.(users.User); ok {
		if tt, ok := t.(*sql.Tx); ok && t != nil {
			return ut, tt, nil
		} else {
			return users.User{}, nil, errors.New("found no transaction")
		}
	} else {
		return users.User{}, nil, errors.New("found no user")
	}
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}
