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

package anomalies

import (
	"time"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/config"
)

const (
	// aggregationMaxSize is the maximum size of an Elastic Search Aggregation
	aggregationMaxSize = 0x7FFFFFFF

	// queryMaxSize is the maximum size of an Elastic Search Query
	queryMaxSize = 10000
)

// createQueryAccountFilter creates and return a new *elastic.TermQuery on the account
func createQueryAccountFilter(account string) *elastic.TermQuery {
	return elastic.NewTermQuery("usageAccountId", account)
}

// createQueryTimeRange creates and return a new *elastic.RangeQuery based on the duration
// defined by durationBegin and durationEnd.
// durationBegin is reduced by period. This offset is deleted later.
func createQueryTimeRange(durationBegin time.Time, durationEnd time.Time) *elastic.RangeQuery {
	periodDuration := time.Duration(config.AnomalyDetectionBollingerBandPeriod) * 24 * time.Hour
	durationBegin = durationBegin.Add(-periodDuration - 1)
	return elastic.NewRangeQuery("usageStartDate").
		From(durationBegin).To(durationEnd)
}

// getProductElasticSearchParams is used to construct an ElasticSearch *elastic.SearchService
// used to retrieve the average cost by usageType for each day.
// It takes as parameters :
// 	- account string : A string representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//	- durationBeing time.Time : A time.Time struct representing the beginning of the time range in the query
//	- durationEnd time.Time : A time.Time struct representing the end of the time range in the query
//  - aggregationPeriod string : An aggregation period, can be "day"
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "awsdetailedlineitems"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getProductElasticSearchParams(account string, durationBegin time.Time,
	durationEnd time.Time, aggregationPeriod string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	query = query.Filter(createQueryAccountFilter(account))
	query = query.Filter(createQueryTimeRange(durationBegin, durationEnd))
	search := client.Search().Index(index).Size(0).Query(query)

	search.Aggregation("products", elastic.NewTermsAggregation().Field("productCode").Size(aggregationMaxSize).
		SubAggregation("dates", elastic.NewDateHistogramAggregation().Field("usageStartDate").ExtendedBounds(durationBegin, durationEnd).Interval(aggregationPeriod).
			SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost"))))
	return search
}

// getDateRangeElasticSearchParams is used to construct an ElasticSearch *elastic.SearchService
// used to retrieve the start and end date depending on the first and the last product in ElasticSearch.
// It takes as parameters :
// 	- account string : A string representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//  - begin bool : Returns the begin date if true, else returns the end date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "awsdetailedlineitems"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getDateRangeElasticSearchParams(account string, begin bool, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	query = query.Filter(createQueryAccountFilter(account))
	search := client.Search().Index(index).Size(1).Sort("usageStartDate", begin).Query(query)
	return search
}

// getAnomalyElasticSearchParams is used to construct an ElasticSearch *elastic.SearchService
// used to retrieve the anomalies.
// It takes as parameters :
// 	- account string : A string representing aws account number
//	- durationBeing time.Time : A time.Time struct representing the beginning of the time range in the query
//	- durationEnd time.Time : A time.Time struct representing the end of the time range in the query
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	- index string : The Elastic Search index on which to execute the query.
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getAnomalyElasticSearchParams(account string, durationBegin time.Time,
	durationEnd time.Time, client *elastic.Client, index string, anomalyType string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewRangeQuery("date").From(durationBegin).To(durationEnd))
	query = query.Filter(elastic.NewTermQuery("abnormal", true))
	search := client.Search().Index(index).Type(anomalyType).Size(queryMaxSize).Sort("date", true).Query(query)
	return search
}

// addDocToBulkProcessor adds a document in a bulk processor to ingest them in ES
func addDocToBulkProcessor(bp *elastic.BulkProcessor, doc interface{}, docType, index, id string) *elastic.BulkProcessor {
	rq := elastic.NewBulkIndexRequest()
	rq = rq.Index(index)
	rq = rq.Type(docType)
	rq = rq.Id(id)
	rq = rq.Doc(doc)
	bp.Add(rq)
	return bp
}
