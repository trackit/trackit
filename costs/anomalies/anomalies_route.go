//   Copyright 2018 MSolution.IO
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

package anomalies

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/anomaliesDetection"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
	"github.com/trackit/trackit-server/aws/s3"
)

// AnomalyEsQueryParams will store the parsed query params
type AnomalyEsQueryParams struct {
	DateBegin   time.Time
	DateEnd     time.Time
	AccountList []string
	IndexList   []string
}

// elasticSearchSearchParamsGetter represents the function used to get search params
// used to request ElasticSearch.
type elasticSearchSearchParamsGetter func(accountList []string, durationBegin time.Time,
	durationEnd time.Time, client *elastic.Client, index string) *elastic.SearchService

// anomalyQueryArgs allows to get required queryArgs params
var anomalyQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAnomaliesData).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(anomalyQueryArgs),
			routes.Documentation{
				Summary:     "get the cost anomalies",
				Description: "Responds with the cost anomalies based on the query args passed to it",
			},
		),
	}.H().Register("/costs/anomalies")
}

// makeElasticSearchRequest prepares and run the request to retrieve the cost anomalies.
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams AnomalyEsQueryParams, elasticSearchSearchParamsGetter elasticSearchSearchParamsGetter) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.IndexList, ",")
	searchService := elasticSearchSearchParamsGetter(
		parsedParams.AccountList,
		parsedParams.DateBegin,
		parsedParams.DateEnd,
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, errors.GetErrorMessage(ctx, err)
		} else if err.(*elastic.Error).Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}

// getAnomaliesData checks the request and returns AnomaliesData.
func getAnomaliesData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := AnomalyEsQueryParams{
		AccountList: []string{},
		DateBegin:   a[anomalyQueryArgs[1]].(time.Time),
		DateEnd:     a[anomalyQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
	}
	if a[anomalyQueryArgs[0]] != nil {
		parsedParams.AccountList = a[anomalyQueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, anomalies.IndexPrefixAnomaliesDetection)
	if err != nil {
		return returnCode, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	_, returnCode, err = makeElasticSearchRequest(request.Context(), parsedParams, getProductAnomaliesElasticSearchParams)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, nil
		} else {
			return http.StatusInternalServerError, err
		}
	}
	accountsAndIndexes, returnCode, err = es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, s3.IndexPrefixLineItem)
	if err != nil {
		return returnCode, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	products, returnCode, err := makeElasticSearchRequest(request.Context(), parsedParams, getElasticSearchParams)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, nil
		} else {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, products
}
