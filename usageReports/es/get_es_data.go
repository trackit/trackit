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

package es

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"errors"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	tes "github.com/trackit/trackit/aws/usageReports/es"
	terrors "github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/users"
)

// makeElasticSearchRequest prepares and run an ES request
// based on the esQueryParams and search params
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams EsQueryParams,
	esSearchParams func(EsQueryParams, *elastic.Client, string) *elastic.SearchService) (*elastic.SearchResult, int, error) {
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

// GetEsMonthlyDomains does an elastic request and returns an array of domains monthly report based on query params
func GetEsMonthlyDomains(ctx context.Context, params EsQueryParams) (int, []DomainReport, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSearchEsMonthlyParams)
	if err != nil {
		return returnCode, nil, err
	} else if res == nil {
		return http.StatusInternalServerError, nil, errors.New("Error while getting data. Please check again in few hours.")
	}
	domains, err := prepareResponseEsMonthly(ctx, res)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, domains, nil
}

// GetEsDailyDomains does an elastic request and returns an array of domains daily report based on query params
func GetEsDailyDomains(ctx context.Context, params EsQueryParams, user users.User, tx *sql.Tx) (int, []DomainReport, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSearchEsDailyParams)
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
	domains, err := prepareResponseEsDaily(ctx, res, costRes)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, domains, nil
}

// GetEsData gets ES monthly reports based on query params, if there isn't a monthly report, it gets daily reports
func GetEsData(ctx context.Context, parsedParams EsQueryParams, user users.User, tx *sql.Tx) (int, []DomainReport, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, tes.IndexPrefixESReport)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	returnCode, monthlyDomains, err := GetEsMonthlyDomains(ctx, parsedParams)
	if err != nil {
		return returnCode, nil, err
	} else if monthlyDomains != nil && len(monthlyDomains) > 0 {
		return returnCode, monthlyDomains, nil
	}
	returnCode, dailyDomains, err := GetEsDailyDomains(ctx, parsedParams, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return returnCode, dailyDomains, nil
}

// GetEsUnusedData gets ES reports and parse them based on query params to have an array of unused domains
func GetEsUnusedData(ctx context.Context, params EsUnusedQueryParams, user users.User, tx *sql.Tx) (int, []DomainReport, error) {
	returnCode, reports, err := GetEsData(ctx, EsQueryParams{params.AccountList, nil, params.Date}, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return prepareResponseEsUnused(params, reports)
}
