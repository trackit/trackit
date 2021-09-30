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

// Package tags implements the /costs/tags/ routes
package tags

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getTagsValues).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(tagsValuesQueryArgs),
			cache.UsersCache{},
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
			cache.UsersCache{},
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
	{
		Name:        "keys",
		Description: "keys of the tags to search",
		Type:        routes.QueryArgStringSlice{},
		Optional:    true,
	},
	{
		Name:        "by",
		Description: "Criteria for the ES aggregation: product, availabilityzone, region or account.",
		Type:        routes.QueryArgString{},
		Optional:    false,
	},
	routes.DetailedQueryArg,
	{
		Name:        "Detailed",
		Description: "Specify if the report will be detailed or not",
		Type:        routes.QueryArgBool{},
		Optional:    true,
	},
}

// TagsValuesQueryParams will store the parsed query params for /tags/values endpoint
type TagsValuesQueryParams struct {
	AccountList []string  `json:"awsAccounts"`
	IndexList   []string  `json:"indexes"`
	DateBegin   time.Time `json:"begin"`
	DateEnd     time.Time `json:"end"`
	TagsKeys    []string  `json:"keys"`
	By          string    `json:"by"`
	Detailed    bool      `json:"detailed"`
}

// getTagsValues returns tags and their values (cost) based on the query params, in JSON format.
func getTagsValues(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := TagsValuesQueryParams{
		AccountList: []string{},
		IndexList:   []string{},
		DateBegin:   a[tagsValuesQueryArgs[1]].(time.Time),
		DateEnd:     a[tagsValuesQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
		TagsKeys:    []string{},
		By:          a[tagsValuesQueryArgs[4]].(string),
		Detailed:    false,
	}
	if a[tagsValuesQueryArgs[0]] != nil {
		parsedParams.AccountList = a[tagsValuesQueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, s3.IndexPrefixLineItem)
	if err != nil {
		return returnCode, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	if a[tagsValuesQueryArgs[3]] != nil {
		parsedParams.TagsKeys = a[tagsValuesQueryArgs[3]].([]string)
	}
	if a[tagsValuesQueryArgs[5]] != nil {
		parsedParams.Detailed = a[tagsValuesQueryArgs[5]].(bool)
	}
	if getTagsValuesFilter(parsedParams.By).Filter == "error" {
		return http.StatusBadRequest, errors.New("Invalid filter: " + parsedParams.By)
	}
	returnCode, res, err := GetTagsValuesWithParsedParams(request.Context(), parsedParams)
	if returnCode == http.StatusOK {
		return returnCode, res
	}
	return returnCode, err
}

// tagsKeysQueryArgs allows to get required queryArgs params for /tags/keys endpoint
var tagsKeysQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
}

// TagsKeysQueryParams will store the parsed query params for /tags/keys endpoint
type TagsKeysQueryParams struct {
	AccountList []string  `json:"awsAccounts"`
	IndexList   []string  `json:"indexes"`
	DateBegin   time.Time `json:"begin"`
	DateEnd     time.Time `json:"end"`
}

// getTagsKeys returns the list of the tag keys based on the query params, in JSON format.
func getTagsKeys(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := TagsKeysQueryParams{
		AccountList: []string{},
		IndexList:   []string{},
		DateBegin:   a[tagsKeysQueryArgs[1]].(time.Time),
		DateEnd:     a[tagsKeysQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
	}
	if a[tagsKeysQueryArgs[0]] != nil {
		parsedParams.AccountList = a[tagsKeysQueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, s3.IndexPrefixLineItem)
	if err != nil {
		return returnCode, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	returnCode, res, err := GetTagsKeysWithParsedParams(request.Context(), parsedParams)
	if returnCode == http.StatusOK {
		return returnCode, res
	}
	return returnCode, err
}
