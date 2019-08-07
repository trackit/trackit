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

package plugins_account_s3_traffic

import (
	"encoding/json"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	core "github.com/trackit/trackit/plugins/account/core"
)

type bucketsInfos = map[string]float64
type bucket = map[string]interface{}

// parseBuckets iterates through all the buckets and parses their values
func parseBuckets(bandwidthInfos bucketsInfos, parsedDocument bucket) bucketsInfos {
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(bucket)
		if bucketData["key"].(string) != "" {
			bandwidthInfos[bucketData["key"].(string)] = bucketData["usage"].(bucket)["value"].(float64)
		}
	}
	return bandwidthInfos
}

// parseESResult parses an *elastic.SearchResult and returns a map of the bucket names associated with their usage
func parseESResult(pluginParams core.PluginParams, res *elastic.SearchResult) (bucketsInfos, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(pluginParams.Context)
	var parsedDocument bucket
	bandwidthInfos := bucketsInfos{}
	err := json.Unmarshal(*res.Aggregations["buckets"], &parsedDocument)
	if err != nil {
		logger.Error("S3 traffic failed to parse elasticsearch document.", err.Error())
		return bandwidthInfos, err
	}
	bandwidthInfos = parseBuckets(bandwidthInfos, parsedDocument)
	return bandwidthInfos, nil
}
