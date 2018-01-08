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
	"net/url"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/es"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

const (
	iso8601TimeFormat = "2006-01-02"
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
	accountList       []string
	aggregationParams []string
}

type queryParamsParser func([]string, esQueryParams) (esQueryParams, error)

// formFieldsToHandleFunc maps a query param to a queryParamsParser function
// to parse it.
var formFieldsToHandlerFunc = map[string]queryParamsParser{
	"begin":    parseBeginQueryParam,
	"end":      parseEndQueryParam,
	"accounts": parseAccountsQueryParam,
	"by":       parseCriteriaQueryParam,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getCostData).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{},
			routes.Documentation{
				Summary:     "get the costs data",
				Description: "Responds with cost data based on the queryparams passed to it",
			},
		),
	}.H().Register("/costs")
}

func parseDateQueryParam(queryParam []string) (time.Time, error) {
	if len(queryParam) != 1 {
		return time.Now(), fmt.Errorf("expected queryParam len of 1 but got %d", len(queryParam))
	}
	timeToRet, err := time.Parse(iso8601TimeFormat, queryParam[0])
	if err != nil {
		return timeToRet, fmt.Errorf("could not parse time : %s", err.Error())
	}
	return timeToRet, nil
}

// parseBeginQueryParam parses and validates the begining of the time range.
func parseBeginQueryParam(queryParam []string, parsedParams esQueryParams) (esQueryParams, error) {
	var err error
	parsedParams.dateBegin, err = parseDateQueryParam(queryParam)
	return parsedParams, err
}

// parseEndQueryParam parses and validates the end of the time range.
func parseEndQueryParam(queryParam []string, parsedParams esQueryParams) (esQueryParams, error) {
	var err error
	parsedParams.dateEnd, err = parseDateQueryParam(queryParam)
	return parsedParams, err
}

// parseAccountsQueryPAram parses the accounts passed in the queryParams.
func parseAccountsQueryParam(queryParam []string, parsedParams esQueryParams) (esQueryParams, error) {
	if len(queryParam) != 1 {
		return parsedParams, fmt.Errorf("expected queryParam len of 1 but got %d", len(queryParam))
	}
	parsedParams.accountList = strings.Split(queryParam[0], ",")
	return parsedParams, nil
}

// parseCriterionQueryParam parses and validate the different cirterions.
// It validate the criterion by checking its presence in the simpleCriterionMap
// or, in the case of the special criterion tag, will check if it is in the
// correct format : 'tag:*' (with no more than one ':')
// Right now the tags are not enabled and will generate an error if they are
// used because they are not yet implemented in the new ElasticSearch mapping
func parseCriteriaQueryParam(queryParam []string, parsedParams esQueryParams) (esQueryParams, error) {
	if len(queryParam) != 1 {
		return parsedParams, fmt.Errorf("expected queryParam len of 1 but got %d", len(queryParam))
	}
	criteriaList := strings.Split(queryParam[0], ",")
	var resCriterions []string
	for _, criterion := range criteriaList {
		if simpleCriterionMap[criterion] {
			resCriterions = append(resCriterions, criterion)
		} else if len(criterion) >= 5 && criterion[:4] == "tag:" && strings.Count(criterion, ":") == 1 {
			return parsedParams, fmt.Errorf("tags not yet implemented")
		} else {
			return parsedParams, fmt.Errorf("error parsing criterion : %s", criterion)
		}
	}
	parsedParams.aggregationParams = resCriterions
	return parsedParams, nil
}

// parseQueryParams parses the queryParams arguments.
// It parses those arguments into a esQueryParams struct to be used by the
// GetElasticSearchParams function that will create the ElasticSearch request.
func parseQueryParams(queryParams url.Values) (esQueryParams, error) {
	var parsedParams esQueryParams
	var err error
	for key, val := range queryParams {
		if parserFunc := formFieldsToHandlerFunc[key]; parserFunc != nil {
			parsedParams, err = parserFunc(val, parsedParams)
			if err != nil {
				return parsedParams, err
			}
		} else {
			return parsedParams, fmt.Errorf("could not parse param : %s=%s", key, val)
		}
	}
	if parsedParams.dateBegin.IsZero() {
		return parsedParams, fmt.Errorf("'end' is not set")
	} else if parsedParams.dateEnd.IsZero() {
		return parsedParams, fmt.Errorf("'until' is not set")
	} else if len(parsedParams.aggregationParams) < 1 {
		return parsedParams, fmt.Errorf("'by' should at least have one criterion")
	}
	return parsedParams, nil
}

func parseFormAndQueryParams(request *http.Request) (esQueryParams, error) {
	l := jsonlog.LoggerFromContextOrDefault(request.Context())
	err := request.ParseForm()
	if err != nil {
		l.Info("Error parsing form data : "+err.Error(), nil)
		return esQueryParams{}, err
	}
	parsedParams, err := parseQueryParams(request.Form)
	if err != nil {
		l.Info("Error parsing user submitted forms : "+err.Error(), nil)
		return esQueryParams{}, err
	}
	return parsedParams, nil
}

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
		if err.Error() == "elastic: Error 404 (Not Found): no such index [type=index_not_found_exception]" {
			l.Warning("Query execution failed, ES index does not exists : "+index, nil)
			return es.SimplifiedCostsDocument{}, http.StatusOK, err
		}
		l.Error("Query execution failed : "+err.Error(), nil)
		return es.SimplifiedCostsDocument{}, http.StatusInternalServerError, fmt.Errorf("could not execute the ElasticSearch query")
	}
	simplifiedCostDocument, err := es.SimplifyCostsDocument(ctx, res)
	fmt.Printf("%v", simplifiedCostDocument)
	if err != nil {
		l.Error("Error parsing cost response : "+err.Error(), nil)
		return simplifiedCostDocument, http.StatusInternalServerError, fmt.Errorf("could not parse ElasticSearch response")
	}
	return simplifiedCostDocument, http.StatusOK, nil
}

// getCostsData returns the cost data based on the query params, in JSON format.
func getCostData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams, err := parseFormAndQueryParams(request)
	if err != nil {
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
