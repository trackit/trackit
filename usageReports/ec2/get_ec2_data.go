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
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws/usageReports/ec2"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/users"
)

// makeElasticSearchRequest prepares and run an ES request
// based on the ec2QueryParams and search params
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams Ec2QueryParams,
	esSearchParams func(Ec2QueryParams, *elastic.Client, string) *elastic.SearchService) (*elastic.SearchResult, int, error) {
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

// GetEc2MonthlyInstances does an elastic request and returns an array of instances monthly report based on query params
func GetEc2MonthlyInstances(ctx context.Context, params Ec2QueryParams, user users.User, tx *sql.Tx) (int, []InstanceReport, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSearchEc2MonthlyParams)
	if err != nil {
		return returnCode, nil, err
	}
	instances, err := prepareResponseEc2Monthly(ctx, res)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, instances, nil
}

// GetEc2DailyInstances does an elastic request and returns an array of instances daily report based on query params
func GetEc2DailyInstances(ctx context.Context, params Ec2QueryParams, user users.User, tx *sql.Tx) (int, []InstanceReport, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSearchEc2DailyParams)
	if err != nil {
		return returnCode, nil, err
	}
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(params.AccountList, user, tx, es.IndexPrefixLineItems)
	if err != nil {
		return returnCode, nil, err
	}
	params.AccountList = accountsAndIndexes.Accounts
	params.IndexList = accountsAndIndexes.Indexes
	costRes, _, _ := makeElasticSearchRequest(ctx, params, getElasticSearchCostParams)
	instances, err := prepareResponseEc2Daily(ctx, res, costRes)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, instances, nil
}

// GetEc2Data gets EC2 monthly reports based on query params, if there isn't a monthly report, it gets daily reports
func GetEc2Data(ctx context.Context, parsedParams Ec2QueryParams, user users.User, tx *sql.Tx) (int, []InstanceReport, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, ec2.IndexPrefixEC2Report)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	returnCode, monthlyInstances, err := GetEc2MonthlyInstances(ctx, parsedParams, user, tx)
	if err != nil {
		return returnCode, nil, err
	} else if monthlyInstances != nil && len(monthlyInstances) > 0 {
		return returnCode, monthlyInstances, nil
	}
	returnCode, dailyInstances, err := GetEc2DailyInstances(ctx, parsedParams, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return returnCode, dailyInstances, nil
}

// GetEc2UnusedData gets EC2 reports and parse them based on query params to have an array of unused instances
func GetEc2UnusedData(ctx context.Context, params Ec2UnusedQueryParams, user users.User, tx *sql.Tx) (int, []InstanceReport, error) {
	returnCode, instances, err := GetEc2Data(ctx, Ec2QueryParams{params.AccountList, nil, params.Date}, user, tx)
	if err != nil {
		return returnCode, nil, err
	}
	return prepareResponseEc2Unused(params, instances)
}
