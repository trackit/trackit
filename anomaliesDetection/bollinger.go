//   Copyright 2018 MSolution.IO
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
	"math"
	"time"

	"github.com/trackit/trackit-server/config"
)

// min returns the minimum between a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// sum adds every element of a CostAnomaly slice.
func sum(aCosts AnalyzedCosts) float64 {
	var sum float64
	for _, a := range aCosts {
		sum += a.Cost
	}
	return sum
}

// average calculates the average of a CostAnomaly slice.
func average(aCosts AnalyzedCosts) float64 {
	return sum(aCosts) / float64(len(aCosts))
}

// sigma calculates the sigma in the standard deviation formula.
func sigma(aCosts AnalyzedCosts, avg float64) float64 {
	var sigma float64
	for _, a := range aCosts {
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

// cleanAnomalies removes non-abnormal values.
func cleanAnomalies(aCosts AnalyzedCosts) AnalyzedCosts {
	res := make(AnalyzedCosts, 0)
	for index := range aCosts {
		if aCosts[index].Anomaly {
			res = append(res, aCosts[index])
		}
	}
	return res
}

// analyseAnomalies calculates anomalies with Bollinger Bands algorithm and
// const values above. It consists in generating an upper band, which, if
// exceeded, make an alert.
func analyseAnomalies(aCosts AnalyzedCosts) AnalyzedCosts {
	for index := range aCosts {
		if index > 0 {
			a := &aCosts[index]
			tempSliceSize := min(index, config.AnomalyDetectionBollingerBandPeriod)
			tempSlice := aCosts[index-tempSliceSize : index]
			avg := average(tempSlice)
			sigma := sigma(tempSlice, avg)
			deviation := deviation(sigma, tempSliceSize)
			a.UpperBand = avg*config.AnomalyDetectionBollingerBandUpperBandCoefficient + (deviation * config.AnomalyDetectionBollingerBandStandardDeviationCoefficient)
			if a.Cost > a.UpperBand {
				a.Anomaly = true
			}
		}
	}
	return aCosts
}

// addPadding adds a padding if we ask from 10 to 15
// but ES has only from 12 to 15. So 10 11 will be padded.
func addPadding(aCosts AnalyzedCosts, dateBegin time.Time) AnalyzedCosts {
	if cd, err := time.Parse("2006-01-02T15:04:05.000Z", aCosts[0].Meta.Date); err == nil && dateBegin.Before(cd) {
		for i := int(cd.Sub(dateBegin).Hours() / 24); i > 0; i-- {
			cd = cd.AddDate(0, 0, -1)
			pad := AnalyzedCost{
				Meta: AnalyzedCostEssentialMeta{
					Date: cd.Format("2006-01-02T15:04:05.000Z"),
				},
			}
			aCosts = append(AnalyzedCosts{pad}, aCosts...)
		}
	}
	return aCosts
}

// computeAnomalies calls every functions to well format
// AnalyzedCosts and do BollingerBand.
func computeAnomalies(ctx context.Context, aCosts AnalyzedCosts, dateBegin time.Time) AnalyzedCosts {
	aCosts = addPadding(aCosts, dateBegin)
	aCosts = analyseAnomalies(aCosts)
	aCosts = cleanAnomalies(aCosts)
	return aCosts
}
