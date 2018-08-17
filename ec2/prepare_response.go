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

package ec2

import (
	"context"
	"encoding/json"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
)

type bucket = map[string]interface{}

// parseBuckets iterates through all the buckets to retrieve the top hits
func parseBuckets(reports []interface{}, parsedTopReports bucket) []interface{} {
	buckets := parsedTopReports["buckets"].([]interface{})
	for _, bucketTmp := range buckets {
		bucketData := bucketTmp.(bucket)
		bucketTopHits := bucketData["top_reports_hits"].(bucket)
		topHitsHits := bucketTopHits["hits"].(bucket)
		hitsList := topHitsHits["hits"].([]interface{})
		if len(hitsList) > 0 {
			topHit := hitsList[0].(bucket)
			reports = append(reports, topHit["_source"])
		}
	}
	return reports
}

// parseESResult parses an *elastic.SearchResult according to it's resultType
func parseESResult(ctx context.Context, res *elastic.SearchResult) ([]interface{}, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	reports := make([]interface{}, 0)
	var parsedTopReports bucket
	err := json.Unmarshal(*res.Aggregations["top_reports"], &parsedTopReports)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return reports, err
	}
	reports = parseBuckets(reports, parsedTopReports)
	return reports, nil
}

// prepareResponse parses the results from elasticsearch and returns the RDS report
func prepareResponse(ctx context.Context, res *elastic.SearchResult) (interface{}, error) {
	reports, err := parseESResult(ctx, res)
	if err != nil {
		return nil, err
	}
	return reports, nil
}
