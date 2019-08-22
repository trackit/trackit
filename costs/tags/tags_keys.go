//   Copyright 2017-2018 MSolution.IO
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

package tags

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

type (
	// struct that allows to parse ES result
	esTagsKeysResult struct {
		Keys struct {
			Buckets []struct {
				Key string `json:"key"`
			} `json:"buckets"`
		} `json:"keys"`
	}

	// result format of the endpoint
	TagsKeys []string
)

const maxAggregationSize = 0x7FFFFFFF

// getTagsKeysWithParsedParams will parse the data from ElasticSearch and return it
func GetTagsKeysWithParsedParams(ctx context.Context, params TagsKeysQueryParams) (int, TagsKeys, error) {
	var typedDocument esTagsKeysResult
	var response = TagsKeys{}
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	res, returnCode, err := makeElasticSearchRequestForTagsKeys(ctx, params, es.Client)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, response, err
		}
		return returnCode, nil, errors.GetErrorMessage(ctx, err)
	}
	err = json.Unmarshal(*res.Aggregations["data"], &typedDocument)
	if err != nil {
		l.Error("Error while unmarshaling", err)
		return http.StatusInternalServerError, nil, errors.GetErrorMessage(ctx, err)
	}
	for _, key := range typedDocument.Keys.Buckets {
		response = append(response, key.Key)
	}
	return http.StatusOK, response, nil
}

// makeElasticSearchRequestForTagsKeys will make the actual request to the ElasticSearch
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequestForTagsKeys(ctx context.Context, params TagsKeysQueryParams,
	client *elastic.Client) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := getTagsKeysQuery(params)
	index := strings.Join(params.IndexList, ",")
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("data", elastic.NewNestedAggregation().Path("tags").
		SubAggregation("keys", elastic.NewTermsAggregation().Field("tags.key").Size(maxAggregationSize)))
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, err
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, err
	}
	return res, http.StatusOK, nil
}

// getTagsKeysQuery will generate a query for the ElasticSearch based on params
func getTagsKeysQuery(params TagsKeysQueryParams) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilter(params.AccountList))
	}
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(params.DateBegin).To(params.DateEnd))
	return query
}
