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

package costs

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// simpleCriterionMap will map simple criterion to the boolean true.
// This will be used in parseCriterionQueryParams to validate the queryParam.
// It does not take into account the 'tag:*' criterion as it is not fixed.
var simpleCriterionMap = map[string]bool{
	"year":             true,
	"month":            true,
	"week":             true,
	"day":              true,
	"account":          true,
	"product":          true,
	"region":           true,
	"availabilityzone": true,
}

// EsQueryParams will store the parsed query params
type EsQueryParams struct {
	DateBegin         time.Time
	DateEnd           time.Time
	AccountList       []string
	IndexList         []string
	AggregationParams []string
}

// costQueryArgs allows to get required queryArgs params
var costsQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
	routes.QueryArg{
		Name:        "by",
		Description: "Criteria for the ES aggregation, comma separated. Possible values are year, month, week, day, account, product, region, tag(soon)",
		Type:        routes.QueryArgStringSlice{},
		Optional:    false,
	},
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getCostData).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(costsQueryArgs),
			cache.UsersCache{},
			routes.Documentation{
				Summary:     "get the costs data",
				Description: "Responds with cost data based on the query args passed to it",
			},
		),
	}.H().Register("/costs")
}

// validateCriteraParam will validate the different criterions.
// It validate the criterion by checking its presence in the simpleCriterionMap
// or, in the case of the special criterion tag, will check if it is in the
// correct format : 'tag:*' (with no more than one ':')
// Right now the tags are not enabled and will generate an error if they are
// used because they are not yet implemented in the new ElasticSearch mapping
func validateCriteriaParam(parsedParams EsQueryParams) error {
	for _, criterion := range parsedParams.AggregationParams {
		if !simpleCriterionMap[criterion] {
			if len(criterion) >= 5 && criterion[:4] == "tag:" && strings.Count(criterion, ":") == 1 {
				return fmt.Errorf("tags not yet implemented")
			}
			return fmt.Errorf("Error parsing criterion : %s", criterion)
		}
	}
	return nil
}

// MakeElasticSearchRequestAndParseIt will make the actual request to the ElasticSearch parse the results and return them
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func MakeElasticSearchRequestAndParseIt(ctx context.Context, parsedParams EsQueryParams) (es.SimplifiedCostsDocument, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.IndexList, ",")
	searchService := GetElasticSearchParams(
		parsedParams.AccountList,
		parsedParams.DateBegin,
		parsedParams.DateEnd,
		parsedParams.AggregationParams,
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
			return es.SimplifiedCostsDocument{}, http.StatusOK, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return es.SimplifiedCostsDocument{}, http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	simplifiedCostDocument, err := es.SimplifyCostsDocument(ctx, res)
	if err != nil {
		l.Error("Error parsing cost response : "+err.Error(), nil)
		return simplifiedCostDocument, http.StatusInternalServerError, fmt.Errorf("could not parse ElasticSearch response")
	}
	return simplifiedCostDocument, http.StatusOK, nil
}

// getCostsData returns the cost data based on the query params, in JSON format.
func getCostData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := EsQueryParams{
		AccountList:       []string{},
		DateBegin:         a[costsQueryArgs[1]].(time.Time),
		DateEnd:           a[costsQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
		AggregationParams: a[costsQueryArgs[3]].([]string),
	}
	if a[costsQueryArgs[0]] != nil {
		parsedParams.AccountList = a[costsQueryArgs[0]].([]string)
	}
	if err := validateCriteriaParam(parsedParams); err != nil {
		return http.StatusBadRequest, err
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, s3.IndexPrefixLineItem)
	if err != nil {
		return returnCode, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	simplifiedCostDocument, returnCode, err := MakeElasticSearchRequestAndParseIt(request.Context(), parsedParams)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, es.SimplifiedCostsDocument{}.ToJsonable()
		} else {
			return returnCode, err
		}
	}
	return http.StatusOK, simplifiedCostDocument.ToJsonable()
}
