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

type window struct {
	date    string
	oldCost float64
	newCost float64
}

type variation struct {
	Date      string
	Variation string
}

type costVariations map[string][]variation

// ToCSVable generates the CSV content from a costVariations
func (cv costVariations) ToCSVable() [][]string {
	csv := [][]string{}
	keys := []string{}
	for k, _ := range cv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, usageTypeName := range keys {
		if len(csv) == 0 {
			header := []string{"usageType"}
			for _, variation := range cv[usageTypeName] {
				header = append(header, variation.Date)
			}
			csv = append(csv, header)
		}
		row := []string{usageTypeName}
		for _, variation := range cv[usageTypeName] {
			row = append(row, variation.Variation)
		}
		csv = append(csv, row)
	}
	return csv
}

// getWindows takes a list of pricePoint that represent the cost for each week/month
// in the interval selected by the user and returns pairs of costs that will be
// used to calculate the variation from week to week or from month to month
func getWindows(pricePoints []pricePoint) []window {
	windows := make([]window, len(pricePoints)-1)
	for i := 0; i < len(pricePoints)-1; i += 1 {
		windows[i] = window{
			date:    pricePoints[i+1].Date,
			oldCost: pricePoints[i].Cost,
			newCost: pricePoints[i+1].Cost,
		}
	}
	return windows
}

// getVariations generates a list of variation between each pair of consecutive
// week/month in the interval selected by the user.
// The variation is the coefficient applied to the first item to obtain the second one
func getVariations(pricePoints []pricePoint) []variation {
	variations := make([]variation, len(pricePoints)-1)
	for idx, window := range getWindows(pricePoints) {
		variations[idx].Date = window.date
		if window.oldCost == 0.0 {
			variations[idx].Variation = ""
		} else {
			variations[idx].Variation = strconv.FormatFloat((window.newCost/window.oldCost)-1, 'f', -1, 64)
		}
	}
	return variations
}

func parseVariationPricePoints(bucketData usageType) []variation {
	dateAgg := bucketData["dateAgg"].(usageType)
	bucketAggs := dateAgg["buckets"].([]interface{})
	tmpPricePoints := make([]pricePoint, len(bucketAggs)+1)
	for idx, bucketAgg := range bucketAggs {
		tmpPricePoints[idx+1].Date = bucketAgg.(usageType)["key_as_string"].(string)
		tmpPricePoints[idx+1].Cost = bucketAgg.(usageType)["cost"].(map[string]interface{})["value"].(float64)
	}
	return getVariations(tmpPricePoints)
}

func parseVariationUsageTypes(parsedDocument usageType) costVariations {
	absolute := costVariations{}
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(usageType)
		usageTypeName := bucketData["key"].(string)
		absolute[usageTypeName] = parseVariationPricePoints(bucketData)
	}
	return absolute
}

func prepareVariationsData(ctx context.Context, sr *elastic.SearchResult) (costVariations, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument usageType
	err := json.Unmarshal(*sr.Aggregations["usageType"], &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return costVariations{}, err
	}
	return parseVariationUsageTypes(parsedDocument), nil
}
