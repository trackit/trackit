//   Copyright 2021 MSolution.IO
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
package anomalyFilters

import (
	"fmt"

	"github.com/trackit/trackit/costs/anomalies/anomalyType"
)

type (
	// genericFilter implements valid called to validate
	// the data received by postAnomaliesFilters and apply
	// to apply the filter to anomaly results.
	// All filters have to implement genericFilter.
	genericFilter interface {
		valid(data interface{}) error
		apply(data interface{}, res anomalyType.ProductAnomaly, product string) bool
	}
)

var (
	filters     = make(map[string]genericFilter)
	filtersName = make(map[genericFilter]string)
)

// registerFilter has to be called by every filters to register them.
func registerFilter(filterName string, filter genericFilter) {
	filters[filterName] = filter
	filtersName[filter] = filterName
}

// genericValidUnsignedIntegerArray is a generic validation function to
// validate an array of positive integer.
func genericValidUnsignedIntegerArray(filter genericFilter, data interface{}, minBound, maxBound int) error {
	if typed, ok := data.([]interface{}); !ok {
		return fmt.Errorf("%s: not an array", filtersName[filter])
	} else if len(typed) == 0 {
		return fmt.Errorf("%s: empty array", filtersName[filter])
	} else {
		for i := range typed {
			if elemTyped, ok := typed[i].(float64); !ok {
				return fmt.Errorf("%s: not an array of number", filtersName[filter])
			} else if elemTyped < float64(minBound) || elemTyped > float64(maxBound) {
				return fmt.Errorf("%s: not an array of number between %d and %d", filtersName[filter], minBound, maxBound)
			} else if elemTyped != float64(int64(elemTyped)) {
				return fmt.Errorf("%s: not an array of integer", filtersName[filter])
			}
		}
	}
	return nil
}

// Valid verifies the given couple filter / data.
func Valid(filterName string, data interface{}) error {
	if filter, ok := filters[filterName]; !ok {
		return fmt.Errorf("%s: rule not found", filterName)
	} else {
		return filter.valid(data)
	}
}

// Apply applies filters on the response
func Apply(flts []anomalyType.Filter, res anomalyType.AnomaliesDetectionResponse) anomalyType.AnomaliesDetectionResponse {
	for account := range res {
		for product := range res[account] {
			for anomaly, an := range res[account][product] {
				if an.Abnormal && !an.Filtered {
					for _, flt := range flts {
						if filter, ok := filters[flt.Rule]; ok && !flt.Disabled {
							if filter.apply(flt.Data, an, product) {
								res[account][product][anomaly].Filtered = true
								break
							}
						}
					}
				}
			}
		}
	}
	return res
}
