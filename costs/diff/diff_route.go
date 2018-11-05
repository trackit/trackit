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

package diff

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/s3"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

type usageType = map[string]interface{}

// validAggregationPeriodMap is a map that defines the aggregation period
// accepted by the diff route
var validAggregationPeriodMap = map[string]struct{}{
	"month": struct{}{},
	"week":  struct{}{},
}

// esQueryParams will store the parsed query params
type esQueryParams struct {
	dateBegin         time.Time
	dateEnd           time.Time
	accountList       []string
	indexList         []string
	aggregationPeriod string
}

// diffQueryArgs allows to get required queryArgs params
var diffQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
	routes.QueryArg{
		Name:        "by",
		Description: "Criteria for the ES aggregation. Possible values are month, week",
		Type:        routes.QueryArgString{},
		Optional:    false,
	},
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(prepareGetDiffData).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(diffQueryArgs),
			routes.Documentation{
				Summary:     "get the cost diff",
				Description: "Responds with the cost diff based on the query args passed to it",
			},
		),
	}.H().Register("/costs/diff")
}

// makeElasticSearchRequest prepares and run the request to retrieve the billing costs
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
		parsedParams.dateBegin,
		parsedParams.dateEnd,
		parsedParams.aggregationPeriod,
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

// getDiffData returns the cost diff based on the query params, in JSON or CSV format.
func getDiffData(ctx context.Context, parsedParams esQueryParams) (int, interface{}) {
	sr, returnCode, err := makeElasticSearchRequest(ctx, parsedParams)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, nil
		} else {
			return returnCode, err
		}
	}
	res, err := prepareDiffData(ctx, sr)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}

func convertDiffData(ctx context.Context, diffData interface{}) (data costDiff, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Convert running", nil)
	report, ok := diffData.(costDiff)
	if ok == false {
		logger.Error("An error occured while converting to diffData", err)
	}
	return report, err
}

// TaskDiffData prepares an elasticsearch query and retrieves cost differentiator data
func TaskDiffData(ctx context.Context, aa aws.AwsAccount) (data costDiff, err error) {
	now := time.Now()
	dateBegin := time.Date(now.Year(), now.Month() - 1, 1, 0, 0, 0, 0, now.Location()).UTC()
	dateEnd := time.Date(now.Year(), now.Month(), 1, 23, 59, 59, 999999999, now.Location()).UTC().AddDate(0, 0, -1)
	parsedParams := esQueryParams{
		accountList:       []string{es.GetAccountIdFromRoleArn(aa.RoleArn)},
		dateBegin: dateBegin,
		dateEnd: dateEnd,
		aggregationPeriod: "day",
	}
	_, diffData := getDiffData(ctx, parsedParams)
	return convertDiffData(ctx, diffData)
}

func prepareGetDiffData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := esQueryParams{
		accountList:       []string{},
		dateBegin:         a[diffQueryArgs[1]].(time.Time),
		dateEnd:           a[diffQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
		aggregationPeriod: a[diffQueryArgs[3]].(string),
	}
	if a[diffQueryArgs[0]] != nil {
		parsedParams.accountList = a[diffQueryArgs[0]].([]string)
	}
	if _, ok := validAggregationPeriodMap[parsedParams.aggregationPeriod]; ok == false {
		return http.StatusBadRequest, fmt.Errorf("invalid aggregation period : %s", parsedParams.aggregationPeriod)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, s3.IndexPrefixLineItem)
	if err != nil {
		return returnCode, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	sr, returnCode, err := makeElasticSearchRequest(request.Context(), parsedParams)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, nil
		} else {
			return returnCode, err
		}
	}
	_, err = prepareDiffData(request.Context(), sr)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return getDiffData(request.Context(), parsedParams)
}
