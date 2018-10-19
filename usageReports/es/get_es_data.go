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

package es

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
	tes "github.com/trackit/trackit-server/aws/usageReports/es"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/users"
)

// makeElasticSearchEsRequest prepares and run the request to retrieve the latest reports
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empy data
// func makeElasticSearchEsRequest(ctx context.Context, parsedParams esQueryParams) (*elastic.SearchResult, int, error) {
// 	l := jsonlog.LoggerFromContextOrDefault(ctx)
// 	index := strings.Join(parsedParams.indexList, ",")
// 	searchService := GetElasticSearchEsParams(
// 		parsedParams.accountList,
// 		es.Client,
// 		index,
// 	)
// 	res, err := searchService.Do(ctx)
// 	if err != nil {
// 		if elastic.IsNotFound(err) {
// 			l.Warning("Query execution failed, ES index does not exists : "+index, err)
// 			return nil, http.StatusOK, err
// 		}
// 		l.Error("Query execution failed : "+err.Error(), nil)
// 		return nil, http.StatusInternalServerError, fmt.Errorf("could not execute the ElasticSearch query")
// 	}
// 	return res, http.StatusOK, nil
// }

// makeElasticSearchEsRequest prepares and run the request to retrieve the latest reports
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empy data
func makeElasticSearchCostRequest(ctx context.Context, parsedParams esQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchCostParams(
		parsedParams,
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

// makeElasticSearchEsDailyRequest prepares and run the request to retrieve the latest reports
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchEsDailyRequest(ctx context.Context, parsedParams esQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchEsDailyParams(
		parsedParams,
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

// makeElasticSearchEsMonthlyRequest prepares and run the request to retrieve a month report
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchEsMonthlyRequest(ctx context.Context, parsedParams esQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchEsMonthlyParams(
		parsedParams,
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

// getEsDailyData gets RDS daily reports and parse them based on query params
func getEsDailyData(ctx context.Context, params esQueryParams, user users.User, tx *sql.Tx) (int, []tes.Report, error) {
	searchResult, returnCode, err := makeElasticSearchEsDailyRequest(ctx, params)
	if err != nil {
		return returnCode, nil, err
	}
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(params.accountList, user, tx, es.IndexPrefixLineItems)
	if err != nil {
		return returnCode, nil, err
	}
	params.accountList = accountsAndIndexes.Accounts
	params.indexList = accountsAndIndexes.Indexes
	costResult, _, _ := makeElasticSearchCostRequest(ctx, params)
	res, err := prepareResponseEsDaily(ctx, searchResult, costResult)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}

// getEsDailyData gets RDS daily reports and parse them based on query params
func getEsData(request *http.Request, parsedParams esQueryParams, user users.User, tx *sql.Tx) (int, []tes.Report, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, tes.IndexPrefixESReport)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	searchResult, returnCode, err := makeElasticSearchEsMonthlyRequest(request.Context(), parsedParams)
	if err != nil {
		return returnCode, nil, err
	}
	if searchResult.Hits.TotalHits > 0 {
		res, err := prepareResponseEsMonthly(request.Context(), searchResult)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusOK, res, nil
	} else {
		return getEsDailyData(request.Context(), parsedParams, user, tx)
	}
}

func getEsUnusedData(request *http.Request, params esUnusedQueryParams, user users.User, tx *sql.Tx) (int, []tes.Domain, error) {
	returnCode, reports, err := getEsData(request, esQueryParams{params.accountList, nil, params.date}, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return prepareResponseEsUnused(params, reports)
}
