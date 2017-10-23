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
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v5"
)

func flattenSubAggregation(in *interface{}) {
	oldTab := (*in).([]interface{})
	newMap := make(map[string]interface{})
	for _, field := range oldTab {
		for name, subField := range field.(map[string]interface{}) {
			newMap[name] = subField
		}
	}
	*in = newMap
}

func stripAggregationArchitecture(aggregationName string, aggregationType interface{}, in map[string]interface{}) {
	key, ok := in["key_as_string"]
	if !ok {
		key, ok = in["key"]
	}
	skey := fmt.Sprintf("%v", key)
	if ok {
		if aggregationMapType, ok := aggregationType.(map[string]interface{}); ok && aggregationName[0] == '*' {
			aggregationType := aggregationMapType["buckets"]
			flattenSubAggregation(&aggregationType)
		}
		in[skey] = aggregationType
		delete(in, "key")
		delete(in, "key_as_string")
	}
	delete(in, aggregationName)
}

func browseRecursivlyMap(in map[string]interface{}) {
	for aggregationName, field := range in {
		if aggregationName[0] != '&' && aggregationName[0] != '*' {
			continue
		}
		field := field.(map[string]interface{})
		metadata := strings.Split(aggregationName[1:], ".")
		aggregationTypeName := strings.Replace(metadata[0], ",", ".", -1)
		aggregationType := field[aggregationTypeName]
		if aggregationName[0] == '&' {
			browseRecursivlyTab(&aggregationType)
		} else if aggregationName[0] == '*' {
			browseRecursivlyMap(field)
		}
		stripAggregationArchitecture(aggregationName, aggregationType, in)
	}
}

func browseRecursivlyTab(in *interface{}) {
	tab := (*in).([]interface{})
	for _, field := range tab {
		field := field.(map[string]interface{})
		delete(field, "doc_count")
		browseRecursivlyMap(field)
	}
	flattenSubAggregation(in)
}

func prepareRecursiveParsing(result map[string]interface{}, aggregationName string, aggregation *json.RawMessage) error {
	var root map[string]interface{}
	err := json.Unmarshal(*aggregation, &root)
	if err != nil {
		return err
	}
	metadata := strings.Split(aggregationName[1:], ".")
	name := metadata[len(metadata)-1]
	aggregationTypeName := strings.Replace(metadata[0], ",", ".", -1)
	aggregationType := root[aggregationTypeName]
	if aggregationName[0] == '&' {
		browseRecursivlyTab(&aggregationType)
		result[name] = aggregationType
	} else if aggregationName[0] == '*' {
		browseRecursivlyMap(root)
		bucketAggregation := aggregationType.(map[string]interface{})["buckets"]
		flattenSubAggregation(&bucketAggregation)
		result[name] = bucketAggregation
	} else {
		result[aggregationName] = root
	}
	return nil
}

// GetParsedElasticSearchResult is used to parse an *elastic.SearchResult
// to a human readable map[string]interface{} that will be able to be
// transformed to JSON and send to the UI
func GetParsedElasticSearchResult(esResult *elastic.SearchResult) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	for name, aggregation := range esResult.Aggregations {
		if err := prepareRecursiveParsing(res, name, aggregation); err != nil {
			return nil, err
		}
	}
	return res, nil
}
