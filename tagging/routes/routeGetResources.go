//   Copyright 2020 MSolution.IO
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

package routes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/db"
	terrors "github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/taggingReports"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type (
	// ResourcesQueryParams will store the parsed query params
	ResourcesQueryParams struct {
		AccountsList  []string
		IndexesList   []string
		Regions       []string
		ResourceTypes []string
	}
)

// routeGetResources returns the list of resources based on the query params, in JSON format.
func routeGetResources(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	parsedParams := ResourcesQueryParams{
		AccountsList:  []string{},
		Regions:       []string{},
		ResourceTypes: []string{},
	}
	if a[routes.AwsAccountsOptionalQueryArg] != nil {
		parsedParams.AccountsList = a[routes.AwsAccountsOptionalQueryArg].([]string)
	}
	if a[resourcesQueryArgs[1]] != nil {
		parsedParams.Regions = a[resourcesQueryArgs[1]].([]string)
	}
	if a[resourcesQueryArgs[2]] != nil {
		parsedParams.ResourceTypes = a[resourcesQueryArgs[2]].([]string)
	}
	returnCode, report, err := GetResourcesData(request.Context(), parsedParams, user, tx)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}

// GetResourcesData gets resources report based on query params
func GetResourcesData(ctx context.Context, parsedParams ResourcesQueryParams, user users.User, tx *sql.Tx) (int, []taggingReports.TaggingReportDocument, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountsList, user, tx, taggingReports.Model.IndexSuffix)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.AccountsList = accountsAndIndexes.Accounts
	parsedParams.IndexesList = accountsAndIndexes.Indexes
	returnCode, resources, err := GetResources(ctx, parsedParams)
	if err != nil {
		return returnCode, nil, err
	}
	return returnCode, resources, nil
}

// GetResources does an elastic request and returns an array of resources report based on query params
func GetResources(ctx context.Context, params ResourcesQueryParams) (int, []taggingReports.TaggingReportDocument, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, getElasticSeachResourcesParams)
	if err != nil {
		return returnCode, nil, err
	} else if res == nil {
		return http.StatusInternalServerError, nil, errors.New("Error while getting data. Please check again in few hours.")
	}
	resources, err := prepareResponseResources(ctx, res)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, resources, nil
}

// makeElasticSearchRequest prepares and run an ES request
// based on the ResourcesQueryParams and search params
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams ResourcesQueryParams,
	esSearchParams func(ResourcesQueryParams, *elastic.Client, string) *elastic.SearchService) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.IndexesList, ",")
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
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
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
