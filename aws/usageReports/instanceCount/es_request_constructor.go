//   Copyright 2019 MSolution.IO
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

// Package costs gets billing information from an ElasticSearch.
package instanceCount

import (
	"time"

	"gopkg.in/olivere/elastic.v5"
)

// paramAggrAndName is a structure containing the name of the parameter and
// corresponding aggregation. A parameter is a string that is passed to the
// GetElasticSearchParams and represents an aggregation. A list of those
// parameters can be found in the paramNameToFuncPtr map.
type paramAggrAndName struct {
	name string
	aggr elastic.Aggregation
}

// aggregationMaxSize is the maximum size of an Elastic Search Aggregation
const aggregationMaxSize = 0x7FFFFFFF

// createQueryAccountFilter creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilter(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("usageAccountId", accountListFormatted...)
}

// createQueryTimeRange creates and return a new *elastic.RangeQuery based on the duration
// defined by durationBegin and durationEnd
func createQueryTimeRange(durationBegin time.Time, durationEnd time.Time) *elastic.RangeQuery {
	return elastic.NewRangeQuery("usageStartDate").
		From(durationBegin).To(durationEnd)
}

// createQueryTimeRange creates and return a new *elastic.RangeQuery based on the duration
// defined by durationBegin and durationEnd
func createQueryWildcard(field, pattern string) *elastic.WildcardQuery {
	return elastic.NewWildcardQuery(field, pattern)
}

// getElasticSearchParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
func getElasticSearchParams(accountList []string, durationBegin time.Time,
	durationEnd time.Time, client *elastic.Client, index string) *elastic.SearchService {
 	query := elastic.NewBoolQuery()
	if len(accountList) > 0 {
		query = query.Filter(createQueryAccountFilter(accountList))
	}
	query = query.Filter(createQueryTimeRange(durationBegin, durationEnd))
	query = query.Filter(createQueryWildcard("usageType", "*BoxUsage*"))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("region", elastic.NewTermsAggregation().Field("region").Size(aggregationMaxSize).
		SubAggregation("type", elastic.NewTermsAggregation().Field("usageType").Size(aggregationMaxSize).
			SubAggregation("date", elastic.NewDateHistogramAggregation().Field("usageStartDate").MinDocCount(0).Interval("hour").
				SubAggregation("amount", elastic.NewTermsAggregation().Field("usageAmount").Size(aggregationMaxSize)))))
	return search
}
