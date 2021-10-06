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
	// monthDay will only show entries dated
	// the given month days.
	//
	// Format (array of integer between 1 (the 1st) and 31 (the 31st) ):
	// [1, 2, 25, 31]
	monthDay struct{}
)

func init() {
	registerFilter("month_day", monthDay{})
}

// valid verifies the validity of the data
func (f monthDay) valid(data interface{}) error {
	return genericValidUnsignedIntegerArray(f, data, 1, 31)
}

// apply applies the filter to the anomaly and returns the result.
func (f monthDay) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.([]interface{}); !ok {
	} else {
		for _, day := range typed {
			if float64(an.Date.Day()) == day {
				return false
			}
		}
		return true
	}
	return false
}
