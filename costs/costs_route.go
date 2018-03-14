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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/aws"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/es"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"

	"gopkg.in/olivere/elastic.v5"
)

// simpleCriterionMap will map simple criterion to the boolean true.
// This will be used in parseCriterionQueryParams to validate the queryParam.
// It does not take into account the 'tag:*' criterion as it is not fixed.
var simpleCriterionMap = map[string]bool{
	"year":    true,
	"month":   true,
	"week":    true,
	"day":     true,
	"account": true,
	"product": true,
	"region":  true,
}

// esQueryParams will store the parsed query params
type esQueryParams struct {
	dateBegin         time.Time
	dateEnd           time.Time
	accountList       []int
	aggregationParams []string
}

// costQueryArgs allows to get required queryArgs params
var costsQueryArgs = []routes.QueryArg{
	// TODO (BREAKING CHANGE): replace by routes.AwsAccountsOptionalQueryArg
	routes.QueryArg{
		Name:        "accounts",
		Description: "List of comma separated AWS accounts ids",
		Type:        routes.QueryArgIntSlice{},
		Optional:    true,
	},
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
			users.RequireAuthenticatedUser{},
			routes.QueryArgs(costsQueryArgs),
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
func validateCriteriaParam(parsedParams esQueryParams) error {
	for _, criterion := range parsedParams.aggregationParams {
		if !simpleCriterionMap[criterion] {
			if len(criterion) >= 5 && criterion[:4] == "tag:" && strings.Count(criterion, ":") == 1 {
				return fmt.Errorf("tags not yet implemented")
			}
			return fmt.Errorf("Error parsing criterion : %s", criterion)
		}
	}
	return nil
}

// makeElasticSearchRequestAndParseIt will make the actual request to the ElasticSearch parse the results and return them
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empy data
func makeElasticSearchRequestAndParseIt(ctx context.Context, parsedParams esQueryParams, user users.User) (es.SimplifiedCostsDocument, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, "lineitems")
	searchService := GetElasticSearchParams(
		parsedParams.accountList,
		parsedParams.dateBegin,
		parsedParams.dateEnd,
		parsedParams.aggregationParams,
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists : "+index, err)
			return es.SimplifiedCostsDocument{}, http.StatusOK, err
		}
		l.Error("Query execution failed : "+err.Error(), nil)
		return es.SimplifiedCostsDocument{}, http.StatusInternalServerError, fmt.Errorf("could not execute the ElasticSearch query")
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
	parsedParams := esQueryParams{
		accountList:       []int{},
		dateBegin:         a[costsQueryArgs[1]].(time.Time),
		dateEnd:           a[costsQueryArgs[2]].(time.Time),
		aggregationParams: a[costsQueryArgs[3]].([]string),
	}
	if a[costsQueryArgs[0]] != nil {
		parsedParams.accountList = a[costsQueryArgs[0]].([]int)
	}
	if err := validateCriteriaParam(parsedParams); err != nil {
		return http.StatusBadRequest, err
	}
	if err := aws.ValidateAwsAccounts(parsedParams.accountList); err != nil {
		return http.StatusBadRequest, err
	}
	simplifiedCostDocument, returnCode, err := makeElasticSearchRequestAndParseIt(request.Context(), parsedParams, user)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, es.SimplifiedCostsDocument{}.ToJsonable()
		} else {
			return returnCode, err
		}
	}
	return http.StatusOK, simplifiedCostDocument.ToJsonable()
}
