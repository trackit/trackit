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
	"errors"
	"net/http"
	"strings"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/es"
)

type (
	// struct that allows to parse ES result
	esTagsValuesResult struct {
		Tags struct {
			Buckets []struct {
				Key  string `json:"key"`
				Cost struct {
					Value float64 `json:"value"`
				} `json:"cost"`
			} `json:"buckets"`
		}
	}

	// contain a tag and his cost
	TagValue struct {
		Tag  string  `json:"tag"`
		Cost float64 `json:"cost"`
	}

	// contain a key and the list of tags associated
	TagsValues struct {
		Key        string     `json:"key"`
		TagsValues []TagValue `json:"values"`
	}

	// result format of the endpoint
	TagsValuesResponse []TagsValues
)

// getTagsValuesWithParsedParams will parse the data from ElasticSearch and return it
func getTagsValuesWithParsedParams(ctx context.Context, params tagsValuesQueryParams) (int, interface{}) {
	var response = TagsValuesResponse{}
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	for i := range params.TagsKey {
		var typedDocument esTagsValuesResult
		res, returnCode, err := makeElasticSearchRequestForTagsValues(ctx, params, es.Client, i)
		if err != nil {
			if returnCode == http.StatusOK {
				return returnCode, response
			}
			return returnCode, errors.New("Internal server error")
		}
		err = json.Unmarshal(*res.Aggregations["tags"], &typedDocument.Tags)
		if err != nil {
			l.Error("Error while unmarshaling", err)
			return http.StatusInternalServerError, errors.New("Internal server error")
		}
		tagsValues := TagsValues{params.TagsKey[i], nil}
		for _, tag := range typedDocument.Tags.Buckets {
			tagsValues.TagsValues = append(tagsValues.TagsValues, TagValue{tag.Key, tag.Cost.Value})
		}
		response = append(response, tagsValues)
	}
	return http.StatusOK, response
}

// makeElasticSearchRequestForTagsValues will make the actual request to the ElasticSearch
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequestForTagsValues(ctx context.Context, params tagsValuesQueryParams, client *elastic.Client, i int) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := getTagsValuesQuery(params)
	index := strings.Join(params.IndexList, ",")
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("tags", elastic.NewTermsAggregation().Field("tags."+params.TagsKey[i]).
		SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost")))
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists : "+index, err)
			return nil, http.StatusOK, err
		}
		l.Error("Query execution failed : "+err.Error(), nil)
		return nil, http.StatusInternalServerError, err
	}
	return res, http.StatusOK, nil
}

// getTagsValuesQuery will generate a query for the ElasticSearch based on params
func getTagsValuesQuery(params tagsValuesQueryParams) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilter(params.AccountList))
	}
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(params.DateBegin).To(params.DateEnd))
	return query
}

// createQueryAccountFilter creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilter(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("usageAccountId", accountListFormatted...)
}
