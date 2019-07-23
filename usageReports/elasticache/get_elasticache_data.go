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

package elasticache

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"errors"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit/aws/usageReports/elasticache"
	terrors "github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/users"
)

// makeElasticSearchRequest prepares and run an ES request
// based on the elastiCacheQueryParams and search params
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams ElastiCacheQueryParams,
	esSearchParams func(ElastiCacheQueryParams, *elastic.Client, string) *elastic.SearchService) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.IndexList, ",")
	searchService := esSearchParams(
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

// GetElastiCacheMonthlyInstances does an elastic request and returns an array of instances monthly report based on query params
func GetElastiCacheMonthlyInstances(ctx context.Context, params ElastiCacheQueryParams) (int, []InstanceReport, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSearchElastiCacheMonthlyParams)
	if err != nil {
		return returnCode, nil, err
	} else if res == nil {
		return http.StatusInternalServerError, nil, errors.New("Error while getting data. Please check again in few hours.")
	}
	instances, err := prepareResponseElastiCacheMonthly(ctx, res)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, instances, nil
}

// GetElastiCacheDailyInstances does an elastic request and returns an array of instances daily report based on query params
func GetElastiCacheDailyInstances(ctx context.Context, params ElastiCacheQueryParams, user users.User, tx *sql.Tx) (int, []InstanceReport, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSearchElastiCacheDailyParams)
	if err != nil {
		return returnCode, nil, err
	} else if res == nil {
		return http.StatusInternalServerError, nil, errors.New("Error while getting data. Please check again in few hours.")
	}
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(params.AccountList, user, tx, es.IndexPrefixLineItems)
	if err != nil {
		return returnCode, nil, err
	}
	params.AccountList = accountsAndIndexes.Accounts
	params.IndexList = accountsAndIndexes.Indexes
	costRes, _, _ := makeElasticSearchRequest(ctx, params, getElasticSearchCostParams)
	instances, err := prepareResponseElastiCacheDaily(ctx, res, costRes)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, instances, nil
}

// GetElastiCacheData gets ElastiCache monthly reports based on query params, if there isn't a monthly report, it gets daily reports
func GetElastiCacheData(ctx context.Context, parsedParams ElastiCacheQueryParams, user users.User, tx *sql.Tx) (int, []InstanceReport, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, elasticache.IndexPrefixElastiCacheReport)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	returnCode, monthlyInstances, err := GetElastiCacheMonthlyInstances(ctx, parsedParams)
	if err != nil {
		return returnCode, nil, err
	} else if monthlyInstances != nil && len(monthlyInstances) > 0 {
		return returnCode, monthlyInstances, nil
	}
	returnCode, dailyInstances, err := GetElastiCacheDailyInstances(ctx, parsedParams, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return returnCode, dailyInstances, nil
}

// GetElastiCacheUnusedData gets ElastiCache reports and parse them based on query params to have an array of unused instances
func GetElastiCacheUnusedData(ctx context.Context, params ElastiCacheUnusedQueryParams, user users.User, tx *sql.Tx) (int, []InstanceReport, error) {
	returnCode, instances, err := GetElastiCacheData(ctx, ElastiCacheQueryParams{params.AccountList, nil, params.Date}, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return prepareResponseElastiCacheUnused(params, instances)
}
