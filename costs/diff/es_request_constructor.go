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

package diff

import (
	"time"

	"gopkg.in/olivere/elastic.v5"
)

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

// GetElasticSearchParams is used to construct an ElasticSearch *elastic.SearchService
// used to retrieve the cost by usageType for each week/month in the time range
// It takes as paramters :
// 	- accountList []string : A slice of string representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//	- durationBeing time.Time : A time.Time struct representing the begining of the time range in the query
//	- durationEnd time.Time : A time.Time struct representing the end of the time range in the query
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "awsdetailedlineitems"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchParams(accountList []string, durationBegin time.Time,
	durationEnd time.Time, aggregationPeriod string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(accountList) > 0 {
		query = query.Filter(createQueryAccountFilter(accountList))
	}
	query = query.Filter(createQueryTimeRange(durationBegin, durationEnd))
	search := client.Search().Index(index).Size(0).Query(query)

	search.Aggregation("usageType", elastic.NewTermsAggregation().Field("usageType").Size(aggregationMaxSize).
		SubAggregation("dateAgg", elastic.NewDateHistogramAggregation().Field("usageStartDate").MinDocCount(0).ExtendedBounds(durationBegin, durationEnd).Interval(aggregationPeriod).
			SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost"))))
	return search
}
