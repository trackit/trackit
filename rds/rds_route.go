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

package rds

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/rds"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

// esQueryParams will store the parsed query params
type esQueryParams struct {
	accountList []string
	indexList   []string
}

// rdsQueryArgs allows to get required queryArgs params
var rdsQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getRDSReport).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(rdsQueryArgs),
			routes.Documentation{
				Summary:     "get the latest RDS report",
				Description: "Responds with the latest RDS report for the account specified in the request",
			},
		),
	}.H().Register("/rds")
}

// makeElasticSearchRequest prepares and run the request to retrieve the latest reports
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empy data
func makeElasticSearchRequest(ctx context.Context, parsedParams esQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchParams(
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

// getRDSReport returns the list of rds reports based on the query params, in JSON format.
func getRDSReport(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := esQueryParams{
		accountList: []string{},
	}
	if a[rdsQueryArgs[0]] != nil {
		parsedParams.accountList = a[rdsQueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, rds.IndexPrefixRDSReport)
	if err != nil {
		return returnCode, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	searchResult, returnCode, err := makeElasticSearchRequest(request.Context(), parsedParams)
	if err != nil {
		return returnCode, err
	}
	res, err := prepareResponse(request.Context(), searchResult)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}
