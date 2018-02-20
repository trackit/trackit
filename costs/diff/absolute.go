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
	"context"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/trackit/jsonlog"

	"gopkg.in/olivere/elastic.v5"
)

type pricePoint struct {
	Date string
	Cost float64
}

type absoluteCost map[string][]pricePoint

// ToCSVable generates the CSV content from a absoluteCost
func (ac absoluteCost) ToCSVable() [][]string {
	csv := [][]string{}
	keys := []string{}
	for k, _ := range ac {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, usageTypeName := range keys {
		if len(csv) == 0 {
			header := []string{"usageType"}
			for _, costEntry := range ac[usageTypeName] {
				header = append(header, costEntry.Date)
			}
			csv = append(csv, header)
		}
		row := []string{usageTypeName}
		for _, costEntry := range ac[usageTypeName] {
			row = append(row, strconv.FormatFloat(costEntry.Cost, 'f', -1, 64))
		}
		csv = append(csv, row)
	}
	return csv
}

func parseAbsolutePricePoints(bucketData usageType) []pricePoint {
	pricePoints := []pricePoint{}
	dateAgg := bucketData["dateAgg"].(usageType)
	for _, bucketAgg := range dateAgg["buckets"].([]interface{}) {
		pricePoints = append(pricePoints, pricePoint{
			Date: bucketAgg.(usageType)["key_as_string"].(string),
			Cost: bucketAgg.(usageType)["cost"].(map[string]interface{})["value"].(float64),
		})
	}
	return pricePoints
}

func parseAbsoluteUsageTypes(parsedDocument usageType) absoluteCost {
	absolute := absoluteCost{}
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(usageType)
		usageTypeName := bucketData["key"].(string)
		absolute[usageTypeName] = parseAbsolutePricePoints(bucketData)
	}
	return absolute
}

func prepareAbsoluteData(ctx context.Context, sr *elastic.SearchResult) (absoluteCost, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument usageType
	err := json.Unmarshal(*sr.Aggregations["usageType"], &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return absoluteCost{}, err
	}
	return parseAbsoluteUsageTypes(parsedDocument), nil
}
