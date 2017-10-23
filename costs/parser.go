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

// flattenSubAggregation allows to transform an array of map to one flattened map
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

// stripAggregationArchitecture strips an array by deleting some values and calling flattenSubAggregation
func stripAggregationArchitecture(aggrName string, aggrType interface{}, in map[string]interface{}) {
	key, ok := in["key_as_string"]
	if !ok {
		key, ok = in["key"]
	}
	if skey := fmt.Sprintf("%v", key); ok {
		if aggrMapType, ok := aggrType.(map[string]interface{}); ok && aggrName[0] == '*' {
			aggrType := aggrMapType["buckets"]
			flattenSubAggregation(&aggrType)
		}
		in[skey] = aggrType
		delete(in, "key")
		delete(in, "key_as_string")
	}
	delete(in, aggrName)
}

func browseMapRecursivly(in map[string]interface{}) {
	for aggrName, field := range in {
		if aggrName[0] != '&' && aggrName[0] != '*' {
			continue
		}
		field := field.(map[string]interface{})
		metadata := strings.Split(aggrName[1:], ".")
		aggrTypeName := strings.Replace(metadata[0], ",", ".", -1)
		aggrType := field[aggrTypeName]
		if aggrName[0] == '&' {
			browseTabRecursivly(&aggrType)
		} else if aggrName[0] == '*' {
			browseMapRecursivly(field)
		}
		stripAggregationArchitecture(aggrName, aggrType, in)
	}
}

func browseTabRecursivly(in *interface{}) {
	tab := (*in).([]interface{})
	for _, field := range tab {
		field := field.(map[string]interface{})
		delete(field, "doc_count")
		browseMapRecursivly(field)
	}
	flattenSubAggregation(in)
}

func prepareRecursiveParsing(result map[string]interface{}, aggrName string, aggr *json.RawMessage) error {
	var root map[string]interface{}
	err := json.Unmarshal(*aggr, &root)
	if err != nil {
		return err
	}
	metadata := strings.Split(aggrName[1:], ".")
	name := metadata[len(metadata)-1]
	aggrTypeName := strings.Replace(metadata[0], ",", ".", -1)
	aggrType := root[aggrTypeName]
	if aggrName[0] == '&' {
		browseTabRecursivly(&aggrType)
		result[name] = aggrType
	} else if aggrName[0] == '*' {
		browseMapRecursivly(root)
		bucketAggr := aggrType.(map[string]interface{})["buckets"]
		flattenSubAggregation(&bucketAggr)
		result[name] = bucketAggr
	} else {
		result[aggrName] = root
	}
	return nil
}

// GetJSONSimplifiedElasticSearchResult is used to parse an *elastic.SearchResult
// to a human readable map[string]interface{} that will be able to be
// transformed to JSON and send to the UI
func GetJSONSimplifiedElasticSearchResult(esResult *elastic.SearchResult) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	for name, aggregation := range esResult.Aggregations {
		if err := prepareRecursiveParsing(res, name, aggregation); err != nil {
			return nil, err
		}
	}
	return res, nil
}
