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

package rds

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/usageReports/rds"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/users"
)

// makeElasticSearchCostRequests prepares and run the request to retrieve the cost per instance
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchCostRequest(ctx context.Context, params RdsQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(params.indexList, ",")
	searchService := GetElasticSearchCostParams(
		params,
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
				"type": fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}

// makeElasticSearchRdsDailyRequest prepares and run the request to retrieve the latest reports
// based on the RdsQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRdsDailyRequest(ctx context.Context, parsedParams RdsQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchRdsDailyParams(
		parsedParams,
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
				"type": fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}

// makeElasticSearchRdsMonthlyRequest prepares and run the request to retrieve a month report
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRdsMonthlyRequest(ctx context.Context, parsedParams RdsQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchRdsMonthlyParams(
		parsedParams,
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
				"type": fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}

// GetRdsMonthlyInstances does an elastic request and returns an array of instances monthly report based on query params
func GetRdsMonthlyInstances(ctx context.Context, params RdsQueryParams) (int, []rds.InstanceReport, error) {
	res, returnCode, err := makeElasticSearchRdsMonthlyRequest(ctx, params)
	if err != nil {
		return returnCode, nil, err
	}
	instances, err := prepareResponseRdsMonthly(ctx, res)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, instances, nil
}

// GetRdsDailyInstances does an elastic request and returns an array of instances daily report based on query params
func GetRdsDailyInstances(ctx context.Context, params RdsQueryParams, user users.User, tx *sql.Tx) (int, []rds.InstanceReport, error) {
	res, returnCode, err := makeElasticSearchRdsDailyRequest(ctx, params)
	if err != nil {
		return returnCode, nil, err
	}
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(params.AccountList, user, tx, es.IndexPrefixLineItems)
	if err != nil {
		return returnCode, nil, err
	}
	params.AccountList = accountsAndIndexes.Accounts
	params.indexList = accountsAndIndexes.Indexes
	costRes, _, _ := makeElasticSearchCostRequest(ctx, params)
	instances, err := prepareResponseRdsDaily(ctx, res, costRes)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, instances, nil
}

// GetRdsData gets RDS monthly reports based on query params, if there isn't a monthly report, it calls getRdsDailyInstances
func GetRdsData(ctx context.Context, parsedParams RdsQueryParams, user users.User, tx *sql.Tx) (int, []rds.InstanceReport, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, rds.IndexPrefixRDSReport)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	returnCode, monthlyInstances, err := GetRdsMonthlyInstances(ctx, parsedParams)
	if err != nil {
		return returnCode, nil, err
	} else if monthlyInstances != nil && len(monthlyInstances) > 0 {
		return returnCode, monthlyInstances, nil
	}
	returnCode, dailyInstances, err := GetRdsDailyInstances(ctx, parsedParams, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return returnCode, dailyInstances, nil
}

// GetRdsUnusedData gets RDS reports and parse them based on query params to have an array of unused instances
func GetRdsUnusedData(ctx context.Context, params RdsUnusedQueryParams, user users.User, tx *sql.Tx) (int, []rds.InstanceReport, error) {
	returnCode, instances, err := GetRdsData(ctx, RdsQueryParams{params.accountList, nil, params.date}, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return prepareResponseRdsUnused(params, instances)
}
