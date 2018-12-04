//   Copyright 2018 MSolution.IO
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

package plugins_account_s3_traffic

import (
	"time"

	"gopkg.in/olivere/elastic.v5"
)

// aggregationMaxSize is the maximum size of an Elastic Search Aggregation
const aggregationMaxSize = 0x7FFFFFFF

// createQueryTimeRange creates and return a new *elastic.RangeQuery based on the duration
// defined by durationBegin and durationEnd
func createQueryTimeRange(durationBegin time.Time, durationEnd time.Time) *elastic.RangeQuery {
	return elastic.NewRangeQuery("usageStartDate").
		From(durationBegin).To(durationEnd)
}

// GetS3StorageUsage prepared a search query to retrieve the storage usage for each s3 bucket
// Parameters:
// - durationBegin and durationEnd: parameters to define the time interval
// - client: preconfigured elasticsearch client
// - account: the aws account id to use for filtering
// - index: the index on which the query should be executed
// It returns an elastic.SearchService that can be used to execute the search
func GetS3StorageUsage(durationBegin time.Time, durationEnd time.Time, client *elastic.Client, account, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("usageAccountId", account))
	query = query.Filter(createQueryTimeRange(durationBegin, durationEnd))
	query = query.Filter(elastic.NewTermQuery("productCode", "AmazonS3"))
	// We only want to get standard storage
	query = query.Filter(elastic.NewWildcardQuery("usageType", "*-TimedStorage-ByteHrs"))
	search := client.Search().Index(index).Size(0).Query(query)

	search.Aggregation("buckets", elastic.NewTermsAggregation().Field("resourceId").Size(aggregationMaxSize).
		SubAggregation("usage", elastic.NewSumAggregation().Field("usageAmount")))
	return search
}

// GetS3BandwidthUsage prepared a search query to retrieve the bandwidth usage for each s3 bucket
// Parameters:
// - durationBegin and durationEnd: parameters to define the time interval
// - client: preconfigured elasticsearch client
// - account: the aws account id to use for filtering
// - index: the index on which the query should be executed
// It returns an elastic.SearchService that can be used to execute the search
func GetS3BandwidthUsage(durationBegin time.Time, durationEnd time.Time, client *elastic.Client, account, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("usageAccountId", account))
	query = query.Filter(createQueryTimeRange(durationBegin, durationEnd))
	query = query.Filter(elastic.NewTermQuery("productCode", "AmazonS3"))
	query = query.Filter(elastic.NewWildcardQuery("usageType", "*-Bytes"))
	// We only want to get the bandwidth related to objects uploading/downloading
	query = query.Filter(elastic.NewWildcardQuery("operation", "*Object"))
	search := client.Search().Index(index).Size(0).Query(query)

	search.Aggregation("buckets", elastic.NewTermsAggregation().Field("resourceId").Size(aggregationMaxSize).
		SubAggregation("usage", elastic.NewSumAggregation().Field("usageAmount")))
	return search
}
