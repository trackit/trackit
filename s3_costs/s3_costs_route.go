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

package s3_costs

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
	"gopkg.in/olivere/elastic.v5"
)

const (
	iso8601TimeFormat = "2006-01-02"
)

// esQueryParams will store the parsed query params
type esQueryParams struct {
	dateBegin   time.Time
	dateEnd     time.Time
	accountList []string
}

type queryParamsParser func([]string, esQueryParams) (esQueryParams, error)

// formFieldsToHandleFunc maps a query param to a queryParamsParser function
// to parse it.
var formFieldsToHandlerFunc = map[string]queryParamsParser{
	"begin":    parseBeginQueryParam,
	"end":      parseEndQueryParam,
	"accounts": parseAccountsQueryParam,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getS3CostData).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{},
			routes.Documentation{
				Summary:     "get the s3 costs data",
				Description: "Responds with cost data based on the queryparams passed to it",
			},
		),
	}.H().Register("/s3_costs")
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

// parseQueryParams parses the queryParams arguments.
// It parses those arguments into a esQueryParams struct to be used by the
// elasticsearch request constructor functions that will create the ElasticSearch requests.
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
		return parsedParams, fmt.Errorf("'begin' is not set")
	} else if parsedParams.dateEnd.IsZero() {
		return parsedParams, fmt.Errorf("'end' is not set")
	} else if len(parsedParams.accountList) == 0 {
		return parsedParams, fmt.Errorf("'accounts' is not set")
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

// makeElasticSearchStorageRequest prepares and run the request to retrieve storage usage/cost
func makeElasticSearchStorageRequest(ctx context.Context, parsedParams esQueryParams,
	user users.User) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, "lineitems")
	searchService := GetS3SpaceElasticSearchParams(
		parsedParams.accountList,
		parsedParams.dateBegin,
		parsedParams.dateEnd,
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		l.Error("Query execution failed : "+err.Error(), nil)
		return nil, fmt.Errorf("could not execute the ElasticSearch query")
	}
	return res, nil
}

// makeElasticSearchRequestsRequest prepares and run the request to retrieve requests usage/cost
func makeElasticSearchRequestsRequest(ctx context.Context, parsedParams esQueryParams,
	user users.User) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, "lineitems")
	searchService := GetS3RequestsElasticSearchParams(
		parsedParams.accountList,
		parsedParams.dateBegin,
		parsedParams.dateEnd,
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		l.Error("Query execution failed : "+err.Error(), nil)
		return nil, fmt.Errorf("could not execute the ElasticSearch query")
	}
	return res, nil
}

// makeElasticSearchBandwidthRequest prepares and run the request to retrieve bandwidth usage/cost
func makeElasticSearchBandwidthRequest(ctx context.Context, parsedParams esQueryParams,
	user users.User, bwType string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, "lineitems")
	searchService := GetS3BandwidthElasticSearchParams(
		parsedParams.accountList,
		parsedParams.dateBegin,
		parsedParams.dateEnd,
		es.Client,
		index,
		bwType,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		l.Error("Query execution failed : "+err.Error(), nil)
		return nil, fmt.Errorf("could not execute the ElasticSearch query")
	}
	return res, nil
}

// getS3CostData returns the s3 cost data based on the query params, in JSON format.
func getS3CostData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams, err := parseFormAndQueryParams(request)
	if err != nil {
		return http.StatusBadRequest, err
	}
	resStorage, err := makeElasticSearchStorageRequest(request.Context(), parsedParams, user)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	resRequests, err := makeElasticSearchRequestsRequest(request.Context(), parsedParams, user)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	resBandwidthIn, err := makeElasticSearchBandwidthRequest(request.Context(), parsedParams, user, "In")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	resBandwidthOut, err := makeElasticSearchBandwidthRequest(request.Context(), parsedParams, user, "Out")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	res, err := prepareResponse(request.Context(), resStorage, resRequests, resBandwidthIn, resBandwidthOut)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}
