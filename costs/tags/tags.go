//   Copyright 2017-2018 MSolution.IO
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

package tags

import (
	"errors"
	"time"
	"net/http"

	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/users"
	"github.com/trackit/trackit-server/aws"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getTagsValues).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(tagsValuesQueryArgs),
			routes.Documentation{
				Summary:     "get the tag values and their cost with a filter",
				Description: "get the tag values and their cost with filter for a specified time range, aws accounts and keys",
			},
		),
	}.H().Register("/costs/tags/values")
	routes.MethodMuxer{
		http.MethodGet: routes.H(getTagsKeys).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(tagsKeysQueryArgs),
			routes.Documentation{
				Summary:     "get every tag keys",
				Description: "get every tag keys for a specified time range and aws accounts",
			},
		),
	}.H().Register("/costs/tags/keys")
}

// tagsValuesQueryArgs allows to get required queryArgs params for /tags/values endpoint
var tagsValuesQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
	routes.QueryArg{
		Name:        "keys",
		Description: "keys of the tags to search",
		Type:        routes.QueryArgStringSlice{},
		Optional:    true,
	},
	routes.QueryArg{
		Name:        "by",
		Description: "Criteria for the ES aggregation: product, availabilityzone, region or account.",
		Type:        routes.QueryArgString{},
		Optional:    false,
	},
}

// tagsValuesQueryParams will store the parsed query params for /tags/values endpoint
type tagsValuesQueryParams struct {
	AccountList []string  `json:"awsAccounts"`
	DateBegin   time.Time `json:"begin"`
	DateEnd     time.Time `json:"end"`
	TagsKeys    []string  `json:"keys"`
	By          string    `json:"by"`
}

// getTagsValues returns tags and their values (cost) based on the query params, in JSON format.
func getTagsValues(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := tagsValuesQueryParams{
		AccountList: []string{},
		DateBegin:   a[tagsValuesQueryArgs[1]].(time.Time),
		DateEnd:     a[tagsValuesQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
		TagsKeys:    []string{},
		By:          a[tagsValuesQueryArgs[4]].(string),
	}
	if a[tagsValuesQueryArgs[0]] != nil {
		parsedParams.AccountList = a[tagsValuesQueryArgs[0]].([]string)
	}
	if err := aws.ValidateAwsAccounts(parsedParams.AccountList); err != nil {
		return http.StatusBadRequest, err
	}
	if a[tagsValuesQueryArgs[3]] != nil {
		parsedParams.TagsKeys = a[tagsValuesQueryArgs[3]].([]string)
	}
	if getTagsValuesFilter(parsedParams.By) == "error" {
		return http.StatusBadRequest, errors.New("Invalid filter: " + parsedParams.By)
	}
	return getTagsValuesWithParsedParams(request.Context(), parsedParams, user)
}

// tagsKeysQueryArgs allows to get required queryArgs params for /tags/keys endpoint
var tagsKeysQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
}

// tagsKeysQueryParams will store the parsed query params for /tags/keys endpoint
type tagsKeysQueryParams struct {
	AccountList []string  `json:"awsAccounts"`
	DateBegin   time.Time `json:"begin"`
	DateEnd     time.Time `json:"end"`
}

// getTagsKeys returns the list of the tag keys based on the query params, in JSON format.
func getTagsKeys(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := tagsKeysQueryParams{
		AccountList: []string{},
		DateBegin:   a[tagsValuesQueryArgs[1]].(time.Time),
		DateEnd:     a[tagsValuesQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
	}
	if a[tagsKeysQueryArgs[0]] != nil {
		parsedParams.AccountList = a[tagsKeysQueryArgs[0]].([]string)
	}
	if err := aws.ValidateAwsAccounts(parsedParams.AccountList); err != nil {
		return http.StatusBadRequest, err
	}
	return getTagsKeysWithParsedParams(request.Context(), parsedParams, user)
}
