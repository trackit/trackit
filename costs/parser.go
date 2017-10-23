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

func flatSubAggregation(in *interface{}) {
	oldTab := (*in).([]interface{})
	newMap := make(map[string]interface{})
	for _, field := range oldTab {
		for name, subField := range field.(map[string]interface{}) {
			newMap[name] = subField
		}
	}
	*in = newMap
}

func cleanSubAggregation(aggregationName string, aggregationType interface{}, in map[string]interface{}) {
	key, ok := in["key_as_string"]
	if !ok {
		key, ok = in["key"]
	}
	skey := fmt.Sprintf("%v", key)
	if ok {
		if aggregationMapType, ok := aggregationType.(map[string]interface{}); ok && aggregationName[0] == '*' {
			aggregationType := aggregationMapType["buckets"]
			flatSubAggregation(&aggregationType)
		}
		in[skey] = aggregationType
		delete(in, "key")
		delete(in, "key_as_string")
	}
	delete(in, aggregationName)
}

func recurMap(in map[string]interface{}) {
	for aggregationName, field := range in {
		if aggregationName[0] != '&' && aggregationName[0] != '*' {
			continue
		}
		field := field.(map[string]interface{})
		metadata := strings.Split(aggregationName[1:], ".")
		aggregationTypeName := strings.Replace(metadata[0], ",", ".", -1)
		aggregationType := field[aggregationTypeName]
		if aggregationName[0] == '&' {
			recurTab(&aggregationType)
		} else if aggregationName[0] == '*' {
			recurMap(field)
		}
		cleanSubAggregation(aggregationName, aggregationType, in)
	}
}

func recurTab(in *interface{}) {
	tab := (*in).([]interface{})
	for _, field := range tab {
		field := field.(map[string]interface{})
		delete(field, "doc_count")
		recurMap(field)
	}
	flatSubAggregation(in)
}

func prepareRecursiveParsing(result map[string]interface{}, aggregationName string, aggregation *json.RawMessage) {
	var root map[string]interface{}
	err := json.Unmarshal(*aggregation, &root)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	metadata := strings.Split(aggregationName[1:], ".")
	name := metadata[len(metadata)-1]
	aggregationTypeName := strings.Replace(metadata[0], ",", ".", -1)
	aggregationType := root[aggregationTypeName]
	if aggregationName[0] == '&' {
		recurTab(&aggregationType)
		result[name] = aggregationType
	} else if aggregationName[0] == '*' {
		recurMap(root)
		bucketAggregation := aggregationType.(map[string]interface{})["buckets"]
		flatSubAggregation(&bucketAggregation)
		result[name] = bucketAggregation
	} else {
		result[aggregationName] = root
	}
}

// GetParsedElasticSearchResult is used to parse an *elastic.SearchResult
// to a human readable map[string]interface{} that will be able to be
// transform to JSON and send to the UI
func GetParsedElasticSearchResult(esResult *elastic.SearchResult) map[string]interface{} {
	res := make(map[string]interface{})
	for name, aggregation := range esResult.Aggregations {
		prepareRecursiveParsing(res, name, aggregation)
	}
	return res
}
