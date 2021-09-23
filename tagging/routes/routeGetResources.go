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
	"errors"
	"fmt"
	"net/http"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	terrors "github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/tagging"
	"github.com/trackit/trackit/tagging/utils"
	"github.com/trackit/trackit/users"
)

type Tag struct {
	Key   string `json:"key"   req:"nonzero"`
	Value string `json:"value" req:"nonzero"`
}

type ResourcesRequestBody struct {
	Accounts      []string `json:"accounts"      req:"nonzero"`
	Regions       []string `json:"regions"       req:"nonzero"`
	ResourceTypes []string `json:"resourceTypes" req:"nonzero"`
	Tags          []Tag    `json:"tags"          req:"nonzero"`
	MissingTags   []Tag    `json:"missingTags"   req:"nonzero"`
}

// routeGetResources returns the list of resources based on the request body, in JSON format.
func routeGetResources(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body ResourcesRequestBody
	routes.MustRequestBody(a, &body)
	user := a[users.AuthenticatedUser].(users.User)
	returnCode, report, err := GetResources(request.Context(), body, user)
	if err != nil {
		return returnCode, err
	} else {
		return returnCode, report
	}
}

// GetResources does an elastic request and returns an array of resources report based on the request body
func GetResources(ctx context.Context, params ResourcesRequestBody, user users.User) (int, []utils.TaggingReportDocument, error) {
	res, returnCode, err := makeElasticSearchRequest(ctx, params, user, getElasticSeachResourcesParams)
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
// based on the ResourcesRequestBody and search params
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 Internal Server Error status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, params ResourcesRequestBody, user users.User,
	esSearchParams func(ResourcesRequestBody, *elastic.Client, string) *elastic.SearchService) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, tagging.IndexPrefixTaggingReport)
	searchService := esSearchParams(
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
