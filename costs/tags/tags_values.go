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
	FilterType struct {
		Filter string
		Type   string
	}

	// structs that allows to parse ES result
	esTagsValuesDetailedResult struct {
		Keys struct {
			Buckets []struct {
				Key  string               `json:"key"`
				Tags esTagsDetailedResult `json:"tags"`
			} `json:"buckets"`
		} `json:"keys"`
	}

	esTagsDetailedResult struct {
		Buckets []struct {
			Tag string `json:"key"`
			Rev struct {
				Filter esFilterDetailedResult `json:"filter"`
			} `json:"rev"`
		} `json:"buckets"`
	}

	esFilterDetailedResult struct {
		Buckets []struct {
			Time string      `json:"key_as_string"`
			Item interface{} `json:"key"`
			Type struct {
				Buckets []struct {
					Key  string `json:"key"`
					Cost struct {
						Value float64 `json:"value"`
					} `json:"cost"`
				} `json:"buckets"`
			} `json:"type,omitempty"`
			Cost struct {
				Value float64 `json:"value"`
			} `json:"cost,omitempty"`
		} `json:"buckets"`
	}

	// TagValue contains a product and its cost
	TagValue struct {
		Item string  `json:"item"`
		Cost float64 `json:"cost"`
	}

	ValueDetailed struct {
		UsageType string  `json:"usagetype"`
		Cost      float64 `json:"cost"`
	}

	TagValueDetailed struct {
		Item       string          `json:"item"`
		UsageTypes []ValueDetailed `json:"usagetypes"`
	}

	// TagsValues contains a tag and the list of products associated
	TagsValues struct {
		Tag   string             `json:"tag"`
		Costs []TagValue         `json:"costs,omitempty"`
		Items []TagValueDetailed `json:"items,omitempty"`
	}

	// TagsValuesResponse is the response format of the endpoint
	TagsValuesResponse map[string][]TagsValues
)

// GetTagsValuesWithParsedParams will parse the data from ElasticSearch and return it
func GetTagsValuesWithParsedParams(ctx context.Context, params TagsValuesQueryParams) (int, TagsValuesResponse, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	var typedDocument esTagsValuesDetailedResult
	res, returnCode, err := makeElasticSearchRequestForTagsValues(ctx, params, es.Client)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, nil, err
		}
		return returnCode, nil, errors.GetErrorMessage(ctx, err)
	}
	err = json.Unmarshal(*res.Aggregations["data"], &typedDocument)
	if err != nil {
		l.Error("Error while unmarshaling", err)
		return http.StatusInternalServerError, nil, errors.GetErrorMessage(ctx, err)
	}

	var response TagsValuesResponse
	if params.Detailed {
		response = getTagsResponseDetailed(typedDocument, params)
	} else {
		response = getTagsResponse(typedDocument, params)
	}
	return http.StatusOK, response, err
}

//getTagsResponseDetailed get response for tagging when detailed is true
func getTagsResponseDetailed(typedDocument esTagsValuesDetailedResult, params TagsValuesQueryParams) TagsValuesResponse {
	response := TagsValuesResponse{}
	for _, key := range typedDocument.Keys.Buckets {
		var values []TagsValues
		for _, tag := range key.Tags.Buckets {
			var products []TagValueDetailed
			for _, filter := range tag.Rev.Filter.Buckets {
				var valueDetailed []ValueDetailed
				for _, usageType := range filter.Type.Buckets {
					valueDetailed = append(valueDetailed, ValueDetailed{
						UsageType: usageType.Key,
						Cost:      usageType.Cost.Value,
					})
				}
				if filter.Time != "" {
					products = append(products, TagValueDetailed{Item: filter.Time, UsageTypes: valueDetailed})
				} else {
					products = append(products, TagValueDetailed{Item: filter.Item.(string), UsageTypes: valueDetailed})
				}
			}
			values = append(values, TagsValues{Tag: tag.Tag, Items: products})
		}
		if len(params.TagsKeys) == 0 || arrayContainsString(params.TagsKeys, key.Key) {
			response[key.Key] = values
		}
	}
	return response
}

//getTagsResponseDetailed get response for tagging when detailed is false
func getTagsResponse(typedDocument esTagsValuesDetailedResult, params TagsValuesQueryParams) TagsValuesResponse {
	response := TagsValuesResponse{}
	for _, key := range typedDocument.Keys.Buckets {
		var values []TagsValues
		for _, tag := range key.Tags.Buckets {
			var costs []TagValue
			for _, cost := range tag.Rev.Filter.Buckets {
				if cost.Time != "" {
					costs = append(costs, TagValue{cost.Time, cost.Cost.Value})
				} else {
					costs = append(costs, TagValue{cost.Item.(string), cost.Cost.Value})
				}
			}
			values = append(values, TagsValues{Tag: tag.Tag, Costs: costs})
		}
		if len(params.TagsKeys) == 0 || arrayContainsString(params.TagsKeys, key.Key) {
			response[key.Key] = values
		}
	}
	return response
}

//getAggregationForTagsValues get NewReversedNestedAggregation if detailed is true or false
func getAggregationForTagsValues(params TagsValuesQueryParams, filter FilterType) (aggregation *elastic.ReverseNestedAggregation) {
	if params.Detailed {
		aggregation = elastic.NewReverseNestedAggregation().
			SubAggregation("filter", elastic.NewTermsAggregation().Field(filter.Filter).Size(maxAggregationSize).
				SubAggregation("type", elastic.NewTermsAggregation().Field("usageType").Size(maxAggregationSize).
					SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost"))))
		if filter.Type == "time" {
			aggregation = elastic.NewReverseNestedAggregation().
				SubAggregation("filter", elastic.NewDateHistogramAggregation().
					Field("usageStartDate").MinDocCount(0).Interval(filter.Filter).
					SubAggregation("type", elastic.NewTermsAggregation().Field("usageType").Size(maxAggregationSize).
						SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost"))))
		}
		return
	} else {
		aggregation = elastic.NewReverseNestedAggregation().
			SubAggregation("filter", elastic.NewTermsAggregation().Field(filter.Filter).Size(maxAggregationSize).
				SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost")))
		if filter.Type == "time" {
			aggregation = elastic.NewReverseNestedAggregation().
				SubAggregation("filter", elastic.NewDateHistogramAggregation().
					Field("usageStartDate").MinDocCount(0).Interval(filter.Filter).
					SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost")))
		}
		return
	}
}

// makeElasticSearchRequestForTagsValues will make the actual request to the ElasticSearch
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 Internal Server Error status code, it will return the provided status code
// with empty data
func makeElasticSearchRequestForTagsValues(ctx context.Context, params TagsValuesQueryParams, client *elastic.Client) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	filter := getTagsValuesFilter(params.By)
	query := getTagsValuesQuery(params)
	index := strings.Join(params.IndexList, ",")
	aggregation := getAggregationForTagsValues(params, filter)
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("data", elastic.NewNestedAggregation().Path("tags").
		SubAggregation("keys", elastic.NewTermsAggregation().Field("tags.key").Size(maxAggregationSize).
			SubAggregation("tags", elastic.NewTermsAggregation().Field("tags.tag").Size(maxAggregationSize).
				SubAggregation("rev", aggregation))))
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

// getTagsValuesQuery will generate a query for the ElasticSearch based on params
func getTagsValuesQuery(params TagsValuesQueryParams) *elastic.BoolQuery {
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

// getTagsValuesFilter returns a string of the field to filter
func getTagsValuesFilter(filter string) FilterType {
	var filters = map[string]FilterType{
		"product":          {"productCode", "term"},
		"region":           {"region", "term"},
		"account":          {"usageAccountId", "term"},
		"availabilityzone": {"availabilityZone", "term"},
		"day":              {"day", "time"},
		"week":             {"week", "time"},
		"month":            {"month", "time"},
		"year":             {"year", "time"},
	}
	for i := range filters {
		if i == filter {
			return filters[i]
		}
	}
	return FilterType{"error", "error"}
}

// arrayContainsString returns true if a string is present in an array of string
func arrayContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
