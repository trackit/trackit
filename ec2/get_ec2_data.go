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

package ec2

import (
	"fmt"
	"time"
	"context"
	"strings"
	"net/http"
	"database/sql"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/users"
	"github.com/trackit/trackit-server/aws/ec2"
)

// makeElasticSearchCostRequests prepares and run the request to retrieve the cost per instance
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchCostRequest(ctx context.Context, params ec2QueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(params.indexList, ",")
	searchService := GetElasticSearchCostParams(
		params.accountList,
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

// makeElasticSearchEc2Request prepares and run the request to retrieve the latest reports
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchEc2Request(ctx context.Context, parsedParams ec2QueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchEc2Params(
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

// makeElasticSearchEc2HistoryRequest prepares and run the request to retrieve a month report
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchEc2HistoryRequest(ctx context.Context, parsedParams ec2HistoryQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchEc2HistoryParams(
		parsedParams.accountList,
		parsedParams.date,
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

func getEc2Data(request *http.Request, parsedParams ec2QueryParams, user users.User, tx *sql.Tx) (int, []Report, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, ec2.IndexPrefixEC2Report)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	searchResult, returnCode, err := makeElasticSearchEc2Request(request.Context(), parsedParams)
	if err != nil {
		return returnCode, nil, err
	}
	accountsAndIndexes, returnCode, err = es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, es.IndexPrefixLineItems)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	costResult, _, _ := makeElasticSearchCostRequest(request.Context(), parsedParams)
	res, err := prepareResponseEc2(request.Context(), searchResult, costResult)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}

func getEc2HistoryData(request *http.Request, parsedParams ec2HistoryQueryParams, user users.User, tx *sql.Tx) (int, []Report, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, ec2.IndexPrefixEC2Report)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	searchResult, returnCode, err := makeElasticSearchEc2HistoryRequest(request.Context(), parsedParams)
	if err != nil {
		return returnCode, nil, err
	}
	res, err := prepareResponseEc2History(request.Context(), searchResult)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}

func getEc2UnusedData(request *http.Request, params ec2UnusedQueryParams, user users.User, tx *sql.Tx) (int, []Instance, error) {
	var returnCode int
	var reports []Report
	var err error
	now := time.Now().UTC()
	if now.Year() < params.date.Year() || (now.Year() == params.date.Year() && now.Month() <= params.date.Month()) {
		returnCode, reports, err = getEc2Data(request,
			ec2QueryParams{params.accountList, nil}, user, tx)
	} else {
		returnCode, reports, err = getEc2HistoryData(request,
			ec2HistoryQueryParams{params.accountList, nil, params.date}, user, tx)
	}
	if err != nil {
		return returnCode, nil, err
	}
	return prepareResponseEc2Unused(params, reports)
}
