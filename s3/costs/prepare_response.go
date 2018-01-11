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

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"
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

type bucketsInfo = map[string]*S3BucketCost
type bucket = map[string]interface{}

// isValidBucket is used to filter errored requests that appears in the billing data
func isValidBucket(bucketName string) bool {
	return len(bucketName) > 0 && !strings.HasPrefix(bucketName, "[Error:")
}

// getBucketInfoByName returns the S3BucketCost associated to the bucketName
func getBucketInfoByName(buckets bucketsInfo, bucketName string) *S3BucketCost {
	if val, ok := buckets[bucketName]; ok {
		return val
	}
	buckets[bucketName] = &S3BucketCost{}
	return buckets[bucketName]
}

func parseBucketsStorage(buckets bucketsInfo, parsedDocument bucket) {
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(bucket)
		bucketInfo := getBucketInfoByName(buckets, bucketData["key"].(string))
		bucketInfo.GbMonth = bucketData["gb"].(bucket)["value"].(float64)
		bucketInfo.StorageCost = bucketData["cost"].(bucket)["value"].(float64)
	}
}

// parseStorage parses the result from GetS3SpaceElasticSearchParams
func parseStorage(ctx context.Context, buckets bucketsInfo, resStorage *elastic.SearchResult) error {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument bucket
	err := json.Unmarshal(*resStorage.Aggregations["buckets"], &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return err
	}
	parseBucketsStorage(buckets, parsedDocument)
	return nil
}

func parseBucketsRequests(buckets bucketsInfo, parsedDocument bucket) {
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(bucket)
		bucketName := bucketData["key"].(string)
		// The billing data can contain billings for errored requests that we do not want to see
		if isValidBucket(bucketName) {
			bucketInfo := getBucketInfoByName(buckets, bucketName)
			bucketInfo.Requests = bucketData["requests"].(bucket)["value"].(float64)
			bucketInfo.RequestsCost = bucketData["cost"].(bucket)["value"].(float64)
		}
	}
}

// parseRequests parses the result from GetS3RequestsElasticSearchParams
func parseRequests(ctx context.Context, buckets bucketsInfo, resRequests *elastic.SearchResult) error {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument bucket
	err := json.Unmarshal(*resRequests.Aggregations["buckets"], &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return err
	}
	parseBucketsRequests(buckets, parsedDocument)
	return nil
}

func parseBucketsBandwidth(buckets bucketsInfo, parsedDocument bucket, bwType string) {
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(bucket)
		bucketName := bucketData["key"].(string)
		// The billing data can contain billings for errored requests that we do not want to see
		if isValidBucket(bucketName) {
			bucketInfo := getBucketInfoByName(buckets, bucketName)
			if bwType == "In" {
				bucketInfo.DataIn = bucketData["gb"].(bucket)["value"].(float64)
			} else {
				bucketInfo.DataOut = bucketData["gb"].(bucket)["value"].(float64)
			}
			bucketInfo.BandwidthCost += bucketData["cost"].(bucket)["value"].(float64)
		}
	}
}

// parseBandwidth parses the result from GetS3BandwidthElasticSearchParams
func parseBandwidth(ctx context.Context, buckets bucketsInfo, resBandwidth *elastic.SearchResult, bwType string) error {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument bucket
	err := json.Unmarshal(*resBandwidth.Aggregations["buckets"], &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return err
	}
	parseBucketsBandwidth(buckets, parsedDocument, bwType)
	return nil
}

// prepareResponse parses the results from elasticsearch and returns a map of buckets with their usage informations
func prepareResponse(ctx context.Context, resStorage, resRequests, resBandwidthIn, resBandwidthOut *elastic.SearchResult) (interface{}, error) {
	buckets := make(bucketsInfo)
	err := parseStorage(ctx, buckets, resStorage)
	if err != nil {
		return nil, err
	}
	err = parseRequests(ctx, buckets, resRequests)
	if err != nil {
		return nil, err
	}
	err = parseBandwidth(ctx, buckets, resBandwidthIn, "In")
	if err != nil {
		return nil, err
	}
	err = parseBandwidth(ctx, buckets, resBandwidthOut, "Out")
	if err != nil {
		return nil, err
	}
	return buckets, nil
}
