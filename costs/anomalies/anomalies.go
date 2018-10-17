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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/config"
)

type (
	// esDatesBucket is used to store the raw ElasticSearch response.
	esDatesBucket struct {
		Key string `json:"key_as_string"`
		Cost struct {
			Value float64 `json:"value"`
		} `json:"cost"`
	}

	// esTypedResult is	used to store the raw ElasticSearch response.
	esTypedResult struct {
		Products struct {
			Buckets []struct {
				Key string `json:"key"`
				Dates struct {
					Buckets []esDatesBucket `json:"buckets"`
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
		Date        string  `json:"date"`
		Cost        float64 `json:"cost"`
		UpperBand   float64 `json:"upper_band"`
		Abnormal    bool    `json:"abnormal"`
		Level       int     `json:"level"`
		PrettyLevel string  `json:"pretty_level"`
	}

	// ProductsCostAnomalies is used as http response.
	// Keys are products and values are a slice of CostAnomaly.
	ProductsCostAnomalies map[string][]CostAnomaly
)

// getAnomalyLevel get the level of an anomaly.
func getAnomalyLevel(an CostAnomaly) int {
	percent := (an.Cost * 100) / an.UpperBand
	levels := strings.Split(config.AnomalyDetectionLevels, ",")
	for i, level := range levels[1:] {
		l, _ := strconv.ParseFloat(level, 64)
		if percent < l {
			return i
		}
	}
	return len(levels) - 1
}

// setAnomaliesLevel set the level for all anomalies.
func setAnomaliesLevel(costAnomalies []CostAnomaly) []CostAnomaly {
	prettyLevels := strings.Split(config.AnomalyDetectionPrettyLevels, ",")
	for index, an := range costAnomalies {
		if an.Abnormal {
			l := getAnomalyLevel(an)
			costAnomalies[index].Level = l
			costAnomalies[index].PrettyLevel = prettyLevels[l]
		}
	}
	return costAnomalies
}

// min returns the minimum between a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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

// getHighestSpenders gets a podium of the highest spenders.
func getHighestSpenders(c ProductsCostAnomalies, dateString string) (res []string) {
	date, err := time.Parse("2006-01-02T15:04:05.000Z", dateString)
	lowestDate := date.Add(-time.Duration(config.AnomalyDetectionDisturbanceCleaningHighestSpendingPeriod) * 24 * time.Hour)
	if err != nil {
		return
	}
	type costWithSpender struct {
		spender string
		cost    float64
	}
	var costWithSpenders []costWithSpender
	for spender, costAnomalies := range c {
		var totalCost float64
		for _, an := range costAnomalies {
			if cd, err := time.Parse("2006-01-02T15:04:05.000Z", an.Date); err == nil &&
				cd.After(lowestDate) && (cd.Before(date) || cd.Equal(date)) {
				totalCost += an.Cost
			}
		}
		costWithSpenders = append(costWithSpenders, costWithSpender{spender, totalCost})
	}
	sort.Slice(costWithSpenders, func(i, j int) bool {
		return costWithSpenders[i].cost > costWithSpenders[j].cost
	})
	for i := 0; i < config.AnomalyDetectionDisturbanceCleaningHighestSpendingMinRank; i++ {
		res = append(res, costWithSpenders[i].spender)
	}
	return
}

// clearDisturbances clears the fake alerts.
func clearDisturbances(service string, costAnomalies []CostAnomaly, totalCostAnomalies map[string]float64, c ProductsCostAnomalies) []CostAnomaly {
	for index, an := range costAnomalies {
		if an.Abnormal {
			date := an.Date
			increaseAmount := costAnomalies[index].Cost - costAnomalies[index].UpperBand
			if increaseAmount < totalCostAnomalies[date]*config.AnomalyDetectionDisturbanceCleaningMinPercentOfDailyBill/100 || increaseAmount < config.AnomalyDetectionDisturbanceCleaningMinAbsoluteCost {
				costAnomalies[index].Abnormal = false
			} else {
				spenderInPodium := false
				for _, spender := range getHighestSpenders(c, date) {
					if spender == service {
						spenderInPodium = true
					}
				}
				if spenderInPodium == false {
					costAnomalies[index].Abnormal = false
				}
			}
		}
	}
	return costAnomalies
}

// getTotalCostAnomalies gets the total cost by product.
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
		if len(key) > 0 {
			for index := range costAnomalies {
				if index > 0 {
					a := &costAnomalies[index]
					tempSliceSize := min(index, config.AnomalyDetectionBollingerBandPeriod)
					tempSlice := costAnomalies[index-tempSliceSize : index]
					avg := average(tempSlice)
					sigma := sigma(tempSlice, avg)
					deviation := deviation(sigma, tempSliceSize)
					a.UpperBand = avg*config.AnomalyDetectionBollingerBandUpperBandCoefficient + (deviation * config.AnomalyDetectionBollingerBandStandardDeviationCoefficient)
					if a.Cost > a.UpperBand {
						a.Abnormal = true
					}
				}
			}
		}
		c[key] = clearDisturbances(key, costAnomalies, totalCostAnomalies, c)
		c[key] = setAnomaliesLevel(c[key])
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
				0,
				"",
			})
		}
		c[product.Key] = costAnomalies
	}
	return analyseAnomalies(c)
}

// addPadding adds a padding if we ask from 10 to 15
// but ES has only from 12 to 15. So 10 11 will be padded.
func addPadding(typedDocument esTypedResult) esTypedResult {
	date := func() (minDate time.Time) {
		minDate = time.Unix(1<<31-1, 0)
		for _, product := range typedDocument.Products.Buckets {
			if len(product.Dates.Buckets) > 0 {
				if cd, err := time.Parse("2006-01-02T15:04:05.000Z", product.Dates.Buckets[0].Key); err == nil && minDate.After(cd) {
					minDate = cd
				}
			}
		}
		return
	}()
	for idx := range typedDocument.Products.Buckets {
		if len(typedDocument.Products.Buckets[idx].Dates.Buckets) > 0 {
			if cd, err := time.Parse("2006-01-02T15:04:05.000Z", typedDocument.Products.Buckets[idx].Dates.Buckets[0].Key); err == nil && date.Before(cd) {
				for i := int(cd.Sub(date).Hours()/24); i > 0; i-- {
					cd = cd.AddDate(0, 0, -1)
					var pad esDatesBucket
					pad.Key = cd.Format("2006-01-02T15:04:05.000Z")
					pad.Cost.Value = 0
					typedDocument.Products.Buckets[idx].Dates.Buckets = append([]esDatesBucket{pad}, typedDocument.Products.Buckets[idx].Dates.Buckets...)
				}
			}
		}
	}
	return typedDocument
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
	typedDocument = addPadding(typedDocument)
	return parseAnomalies(typedDocument), nil
}
