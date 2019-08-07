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

package costs

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/olivere/elastic"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/errors"
)

// S3BucketCost represents the detail of the costs for a given bucket
type S3BucketCost struct {
	GbMonth       float64
	StorageCost   float64
	BandwidthCost float64
	DataIn        float64
	DataOut       float64
	RequestsCost  float64
	Requests      float64
}

type BucketsInfo = map[string]*S3BucketCost
type bucket = map[string]interface{}

// bucketCostGetter represents a function used to retrieve infos from a bucket
type bucketCostGetter func(*S3BucketCost, bucket) *S3BucketCost

// resultTypeToBucketCostGetter maps result type names to their associated
// getter function pointer
var resultTypeToBucketCostGetter = map[string]bucketCostGetter{
	"storage":      getBucketStorage,
	"requests":     getBucketRequests,
	"bandwidthIn":  getBucketBandwidthIn,
	"bandwidthOut": getBucketBandwidthOut,
}

// isValidBucket is used to filter errored requests that appears in the billing data
func isValidBucket(bucketName string) bool {
	return len(bucketName) > 0 && !strings.HasPrefix(bucketName, "[Error:")
}

// getBucketInfoByName returns the S3BucketCost associated to the bucketName
func getBucketInfoByName(buckets BucketsInfo, bucketName string) *S3BucketCost {
	if val, ok := buckets[bucketName]; ok {
		return val
	}
	buckets[bucketName] = &S3BucketCost{}
	return buckets[bucketName]
}

// getBucketStorage gets the storage informations from the bucketData
// and returns an *S3BucketCost filled with these data
func getBucketStorage(bucketInfo *S3BucketCost, bucketData bucket) *S3BucketCost {
	bucketInfo.GbMonth = bucketData["usage"].(bucket)["value"].(float64)
	bucketInfo.StorageCost = bucketData["cost"].(bucket)["value"].(float64)
	return bucketInfo
}

// getBucketRequests gets the requests informations from the bucketData
// and returns an *S3BucketCost filled with these data
func getBucketRequests(bucketInfo *S3BucketCost, bucketData bucket) *S3BucketCost {
	bucketInfo.Requests = bucketData["usage"].(bucket)["value"].(float64)
	bucketInfo.RequestsCost = bucketData["cost"].(bucket)["value"].(float64)
	return bucketInfo
}

// getBucketBandwidthIn gets the bandwidthIn informations from the bucketData
// and returns an *S3BucketCost filled with these data
func getBucketBandwidthIn(bucketInfo *S3BucketCost, bucketData bucket) *S3BucketCost {
	bucketInfo.DataIn = bucketData["usage"].(bucket)["value"].(float64)
	bucketInfo.BandwidthCost += bucketData["cost"].(bucket)["value"].(float64)
	return bucketInfo
}

// getBucketBandwidthOut gets the bandwidthOut informations from the bucketData
// and returns an *S3BucketCost filled with these data
func getBucketBandwidthOut(bucketInfo *S3BucketCost, bucketData bucket) *S3BucketCost {
	bucketInfo.DataOut = bucketData["usage"].(bucket)["value"].(float64)
	bucketInfo.BandwidthCost += bucketData["cost"].(bucket)["value"].(float64)
	return bucketInfo
}

// parseBuckets iterates through all the buckets and calls the getter function corresponding
// to the resultType
func parseBuckets(buckets BucketsInfo, parsedDocument bucket, resultType string) BucketsInfo {
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(bucket)
		bucketName := bucketData["key"].(string)
		// The billing data can contain billings for errored requests that we do not want to see
		if isValidBucket(bucketName) {
			bucketInfo := getBucketInfoByName(buckets, bucketName)
			if resultTypePtr, ok := resultTypeToBucketCostGetter[resultType]; ok {
				bucketInfo = resultTypePtr(bucketInfo, bucketData)
			}
		}
	}
	return buckets
}

// parseESResult parses an *elastic.SearchResult according to it's resultType
func parseESResult(ctx context.Context, buckets BucketsInfo, res *elastic.SearchResult, resultType string) (BucketsInfo, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument bucket
	err := json.Unmarshal(*res.Aggregations["buckets"], &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return buckets, errors.GetErrorMessage(ctx, err)
	}
	buckets = parseBuckets(buckets, parsedDocument, resultType)
	return buckets, nil
}

// prepareResponse parses the results from elasticsearch and returns a map of buckets with their usage informations
func prepareResponse(ctx context.Context, resStorage, resRequests, resBandwidthIn, resBandwidthOut *elastic.SearchResult) (BucketsInfo, error) {
	buckets := make(BucketsInfo)
	var err error
	var components = [...]struct {
		k  string
		sr *elastic.SearchResult
	}{
		{"storage", resStorage},
		{"requests", resRequests},
		{"bandwidthIn", resBandwidthIn},
		{"bandwidthOut", resBandwidthOut},
	}
	for _, cpn := range components {
		buckets, err = parseESResult(ctx, buckets, cpn.sr, cpn.k)
		if err != nil {
			return nil, err
		}
	}
	return buckets, nil
}
