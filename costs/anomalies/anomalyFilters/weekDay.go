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
	"github.com/trackit/trackit/costs/anomalies/anomalyType"
)

type (
	// weekDay will only show entries dated
	// the given week days.
	//
	// Format (array of integer between 0 (Monday) and 6 (Sunday) ):
	// [0, 1, 2, 3, 4]
	weekDay struct{}
)

func init() {
	registerFilter("week_day", weekDay{})
}

// valid verifies the validity of the data
func (f weekDay) valid(data interface{}) error {
	return genericValidUnsignedIntegerArray(f, data, 0, 6)
}

// apply applies the filter to the anomaly and returns the result.
func (f weekDay) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.([]interface{}); !ok {
	} else {
		for _, day := range typed {
			if float64(an.Date.Weekday())-1 == day {
				return false
			}
		}
		return true
	}
	return false
}
