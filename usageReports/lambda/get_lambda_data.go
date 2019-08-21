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

package lambda

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/trackit/trackit/pagination"
	"net/http"
	"strings"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws/usageReports/lambda"
	terrors "github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/users"
)

// makeElasticSearchRequest prepares and run an ES request
// based on the lambdaQueryParams and search params
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams LambdaQueryParams,
	esSearchParams func(LambdaQueryParams, *elastic.Client, string) *elastic.SearchService) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.IndexList, ",")
	searchService := esSearchParams(
		parsedParams,
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if res != nil && res.Hits != nil {
		fmt.Printf("Total hits from es requestion: %v\n", res.TotalHits())
	}
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, terrors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, terrors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}

// GetLambdaDailyFunctions does an elastic request and returns an array of functions daily report based on query params
func GetLambdaDailyFunctions(ctx context.Context, params LambdaQueryParams, user users.User, tx *sql.Tx) (int, []FunctionReport, LambdaQueryParams, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSearchLambdaDailyParams)
	//search.Aggregation("accounts", elastic.NewCompositeAggregation().Field("account"))
	//val := elastic.NewCompositeAggregation().
	//	SubAggregation("compose", elastic.NewCompositeAggregation().Size(5).Sources(elastic.NewCompositeAggregationTermsValuesSource("by_dates").Field("reportDate"))).
	//	SubAggregation("accounts", elastic.NewTermsAggregation().Field("account").
	//		SubAggregation("dates", elastic.NewTermsAggregation().Field("reportDate").
	//			SubAggregation("functions", elastic.NewTopHitsAggregation().Sort("reportDate", false))))
	//valSource, _ := val.Source()
	//data, _ := json.Marshal(valSource)
	//searchComposite, err := es.Client.Search().Source(valSource).Do(ctx)
	//if searchComposite != nil && searchComposite.Hits != nil {
	//	fmt.Printf("Source for composite aggregation: '%+v' & total hits: %v\n", string(data), searchComposite.Hits.TotalHits)
	//} else {
	//	fmt.Print("Search composite is nil\n")
	//}
	if err != nil {
		return returnCode, nil, params, err
	} else if res == nil {
		return http.StatusInternalServerError, nil, params, errors.New("Error while getting data. Please check again in few hours.")
	}
	functions, err := prepareResponseLambdaDaily(ctx, res)
	if err != nil {
		return http.StatusInternalServerError, nil, params, err
	}
	pagination.StoreTotalHits(&params.Pagination, res.Hits.TotalHits, len(functions))
	return http.StatusOK, functions, params, nil
}

// GetLambdaData gets Lambda monthly reports based on query params, if there isn't a monthly report, it gets daily reports
func GetLambdaData(ctx context.Context, parsedParams LambdaQueryParams, user users.User, tx *sql.Tx) (int, []FunctionReport, LambdaQueryParams, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, lambda.IndexPrefixLambdaReport)
	if err != nil {
		return returnCode, nil, parsedParams, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	returnCode, dailyFunctions, parsedParams, err := GetLambdaDailyFunctions(ctx, parsedParams, user, tx)
	fmt.Printf("Daily func from parsed params, hits: %v\n", parsedParams.Pagination.TotalElements)
	if err != nil {
		return returnCode, nil, parsedParams, err
	}
	return returnCode, dailyFunctions, parsedParams, nil
}
