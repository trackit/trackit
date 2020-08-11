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

package plugins_account_core

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/accountPlugins"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// pluginsQueryParams will store the parsed query params
type pluginsQueryParams struct {
	accountList []string
	indexList   []string
}

// pluginsQueryArgs allows to get required queryArgs params
var pluginsQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getPluginsResults).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(pluginsQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get the latests plugins results",
				Description: "Responds with the latests plugins results for the account(s) specified in the request",
			},
		),
	}.H().Register("/plugins/results")
}

// makeElasticSearchPluginsRequest prepares and run the request to retrieve the latest plugins results
// based on the pluginsQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empy data
func makeElasticSearchPluginsRequest(ctx context.Context, parsedParams pluginsQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchPluginsParams(
		parsedParams.accountList,
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists : "+index, err)
			return nil, http.StatusOK, err
		}
		l.Error("Query execution failed : "+err.Error(), nil)
		return nil, http.StatusInternalServerError, fmt.Errorf("could not execute the ElasticSearch query")
	}
	return res, http.StatusOK, nil
}

// getPluginsResults returns the list of plugins results based on the query params, in JSON format.
func getPluginsResults(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := pluginsQueryParams{
		accountList: []string{},
	}
	if a[pluginsQueryArgs[0]] != nil {
		parsedParams.accountList = a[pluginsQueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, accountPlugins.Model.IndexSuffix)
	if err != nil {
		return returnCode, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	pluginsResult, returnCode, err := makeElasticSearchPluginsRequest(request.Context(), parsedParams)
	if err != nil {
		return returnCode, err
	}
	res, err := prepareResponse(request.Context(), pluginsResult)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}
