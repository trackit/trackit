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
	"fmt"
	"sort"
	"strconv"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/errors"
)

type pricePoint struct {
	Date             string
	Cost             float64
	PercentVariation float64
}

type costDiff map[string][]pricePoint

// ToCSVable generates the CSV content from a costDiff
func (cd costDiff) ToCSVable() [][]string {
	csv := [][]string{}
	keys := []string{}
	for k, _ := range cd {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, usageTypeName := range keys {
		if len(csv) == 0 {
			header := []string{"usageType"}
			for _, costEntry := range cd[usageTypeName] {
				header = append(header, fmt.Sprintf("cost-%s", costEntry.Date), fmt.Sprintf("variation-%s", costEntry.Date))
			}
			csv = append(csv, header)
		}
		row := []string{usageTypeName}
		for _, costEntry := range cd[usageTypeName] {
			variationStr := ""
			if costEntry.PercentVariation > 0.0 {
				variationStr = fmt.Sprintf("+%s%%", strconv.FormatFloat(costEntry.PercentVariation, 'f', 3, 64))
			} else {
				variationStr = fmt.Sprintf("%s%%", strconv.FormatFloat(costEntry.PercentVariation, 'f', 3, 64))
			}
			row = append(row, strconv.FormatFloat(costEntry.Cost, 'f', -1, 64), variationStr)
		}
		csv = append(csv, row)
	}
	return csv
}

// getVariations compute the percentage of variation between each pair of consecutive
// week/month in the interval selected by the user.
func getVariations(pricePoints []pricePoint) []pricePoint {
	for i := 1; i < len(pricePoints); i += 1 {
		if pricePoints[i-1].Cost == 0.0 {
			pricePoints[i].PercentVariation = 0.0
		} else {
			pricePoints[i].PercentVariation = ((pricePoints[i].Cost - pricePoints[i-1].Cost) / pricePoints[i-1].Cost) * 100
		}
	}
	return pricePoints
}

func parseDiffPricePoints(bucketData usageType) []pricePoint {
	pricePoints := []pricePoint{}
	dateAgg := bucketData["dateAgg"].(usageType)
	for _, bucketAgg := range dateAgg["buckets"].([]interface{}) {
		pricePoints = append(pricePoints, pricePoint{
			Date: bucketAgg.(usageType)["key_as_string"].(string),
			Cost: bucketAgg.(usageType)["cost"].(map[string]interface{})["value"].(float64),
		})
	}
	return getVariations(pricePoints)
}

func parseDiffUsageTypes(parsedDocument usageType) costDiff {
	absolute := costDiff{}
	bucketsField := parsedDocument["buckets"].([]interface{})
	for _, bucketData := range bucketsField {
		bucketData := bucketData.(usageType)
		usageTypeName := bucketData["key"].(string)
		absolute[usageTypeName] = parseDiffPricePoints(bucketData)
	}
	return absolute
}

func prepareDiffData(ctx context.Context, sr *elastic.SearchResult) (costDiff, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument usageType
	err := json.Unmarshal(*sr.Aggregations["usageType"], &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return costDiff{}, errors.GetErrorMessage(ctx, err)
	}
	return parseDiffUsageTypes(parsedDocument), nil
}
