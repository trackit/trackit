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
				Summary:     "get value for a tags's key",
				Description: "take in query a time range, aws accounts and a tags's key to get tags and costs of the key",
			},
		),
	}.H().Register("/tags/values")
	routes.MethodMuxer{
		http.MethodGet: routes.H(getTagsKeys).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(tagsKeysQueryArgs),
			routes.Documentation{
				Summary:     "get tags's keys list",
				Description: "return the list of the tags's key for one, many are all aws accounts list",
			},
		),
	}.H().Register("/tags/keys")
}

var tagsValuesQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
	routes.QueryArg{
		Name:        "key",
		Description: "key of the tags to search",
		Type:        routes.QueryArgString{},
		Optional:    false,
	},
}

type tagsValuesQueryParams struct {
	DateBegin   time.Time `json:"begin"`
	DateEnd     time.Time `json:"end"`
	AccountList []string  `json:"awsAccounts"`
	TagsKey     string    `json:"key"`
}

func getTagsValues(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := tagsValuesQueryParams{
		AccountList: []string{},
		DateBegin:   a[tagsValuesQueryArgs[1]].(time.Time),
		DateEnd:     a[tagsValuesQueryArgs[2]].(time.Time),
		TagsKey:     a[tagsValuesQueryArgs[3]].(string),
	}
	if a[tagsValuesQueryArgs[0]] != nil {
		parsedParams.AccountList = a[tagsValuesQueryArgs[0]].([]string)
	}
	if err := aws.ValidateAwsAccounts(parsedParams.AccountList); err != nil {
		return http.StatusBadRequest, err
	}
	return getTagsValuesWithParsedParams(request.Context(), parsedParams, user)
}

var tagsKeysQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
}

type tagsKeysQueryParams struct {
	AccountList []string `json:"awsAccounts"`
}

func getTagsKeys(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := tagsKeysQueryParams{
		AccountList: []string{},
	}
	if a[tagsKeysQueryArgs[0]] != nil {
		parsedParams.AccountList = a[tagsKeysQueryArgs[0]].([]string)
	}
	if err := aws.ValidateAwsAccounts(parsedParams.AccountList); err != nil {
		return http.StatusBadRequest, err
	}
	return getTagsKeysWithParsedParams(request.Context(), parsedParams, user)
}
