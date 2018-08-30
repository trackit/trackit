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

package anomalies

import (
	"context"
	"encoding/json"
	"math"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"
)

type (
	// esTypedResult is	used to store the raw ElasticSearch response.
	esTypedResult struct {
		Products struct {
			Buckets []struct {
				Key string `json:"key"`
				Dates struct {
					Buckets []struct {
						Key string `json:"key_as_string"`
						Cost struct {
							Value float64 `json:"value"`
						} `json:"cost"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		}
	}

	// CostAnomaly represents a day and contains
	// the date of beginning, the cost for the 24h,
	// the value of the upper band and an abnormal value
	// representing if the day is tagged as abnormal,
	// which is an alert.
	CostAnomaly struct {
		Date      string  `json:"date"`
		Cost      float64 `json:"cost"`
		UpperBand float64 `json:"upper_band"`
		Abnormal  bool    `json:"abnormal"`
	}

	// ProductsCostAnomalies is used as http response.
	// Keys are products and values are a slice of CostAnomaly.
	ProductsCostAnomalies map[string][]CostAnomaly
)

// const values used by the Bollinger Bands algorithm.
const (
	// period is the number of day took.
	// A bigger period means more stability in the cost so
	// it will be more sensitive to the picks.
	period = 3

	// standardDeviationCoefficient allows to add a
	// coefficient to the standard deviation.
	// A standardDeviationCoefficient bigger makes
	// the algorithm more flexible.
	standardDeviationCoefficient = 3.0

	// margin of error set to 5% of the price.
	// Upper band will be higher by 5%.
	margin = 1.05

	// minCostPercent is set to 4%.
	// If an anomaly is detected, the cost has to be
	// higher than 4% of the total bill.
	minCostPercent = 0.04
)

// sum adds every element of a CostAnomaly slice.
func sum(costAnomalies []CostAnomaly) float64 {
	var sum float64
	for _, a := range costAnomalies {
		sum += a.Cost
	}
	return sum
}

// average calculates the average of a CostAnomaly slice.
func average(costAnomalies []CostAnomaly) float64 {
	return sum(costAnomalies) / float64(len(costAnomalies))
}

// sigma calculates the sigma in the standard deviation formula.
func sigma(costAnomalies []CostAnomaly, avg float64) float64 {
	var sigma float64
	for _, a := range costAnomalies {
		sigma += math.Pow(a.Cost-avg, 2)
	}
	return sigma
}

// deviation calculates the standard deviation.
func deviation(sigma float64, period int) float64 {
	var deviation float64
	deviation = 1 / float64(period) * math.Pow(sigma, 0.5)
	return deviation
}

// clearDisturbances clears the fake alerts.
// Alerts below the minCostPercent are removed.
func clearDisturbances(costAnomalies []CostAnomaly, totalCostAnomalies map[string]float64) []CostAnomaly {
	for index := range costAnomalies {
		date := costAnomalies[index].Date
		if costAnomalies[index].Cost < totalCostAnomalies[date]*minCostPercent {
			costAnomalies[index].Abnormal = false
		}
	}
	return costAnomalies
}

func getTotalCostAnomalies(c ProductsCostAnomalies) map[string]float64 {
	totalCostAnomalies := map[string]float64{}
	for _, costAnomalies := range c {
		for _, an := range costAnomalies {
			totalCostAnomalies[an.Date] += an.Cost
		}
	}
	return totalCostAnomalies
}

// analyseAnomalies calculates anomalies with Bollinger Bands algorithm and
// the const above. It consists in generating an upper band, which, if
// exceeded, make an alert.
func analyseAnomalies(c ProductsCostAnomalies) ProductsCostAnomalies {
	totalCostAnomalies := getTotalCostAnomalies(c)
	for key, costAnomalies := range c {
		for index := range costAnomalies {
			if index > 0 {
				a := &costAnomalies[index]
				tempSliceSize := int(math.Min(float64(index), period))
				tempSlice := costAnomalies[index-tempSliceSize : index]
				avg := average(tempSlice)
				sigma := sigma(tempSlice, avg)
				deviation := deviation(sigma, tempSliceSize)
				a.UpperBand = avg*margin + (deviation * standardDeviationCoefficient)
				if a.Cost > a.UpperBand {
					a.Abnormal = true
				}
			}
		}
		c[key] = clearDisturbances(costAnomalies, totalCostAnomalies)
	}
	return c
}

// parseAnomalies transforms the esTypedResult in a ProductsCostAnomalies
// empty of alerts. It calls analyseAnomalies then to fill the alerts.
func parseAnomalies(typedDocument esTypedResult) ProductsCostAnomalies {
	c := ProductsCostAnomalies{}
	for _, product := range typedDocument.Products.Buckets {
		costAnomalies := make([]CostAnomaly, 0, len(product.Dates.Buckets))
		for _, date := range product.Dates.Buckets {
			costAnomalies = append(costAnomalies, CostAnomaly{
				date.Key,
				date.Cost.Value,
				0,
				false,
			})
		}
		c[product.Key] = costAnomalies
	}
	return analyseAnomalies(c)
}

// prepareAnomalyData calls ElasticSearch and stores
// the result in a esTypedResult type. It calls parseAnomalies
// then.
func prepareAnomalyData(ctx context.Context, sr *elastic.SearchResult) (ProductsCostAnomalies, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var typedDocument esTypedResult
	err := json.Unmarshal(*sr.Aggregations["products"], &typedDocument.Products)
	if err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return ProductsCostAnomalies{}, err
	}
	return parseAnomalies(typedDocument), nil
}
