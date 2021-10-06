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
	// level will only show entries with
	// the given level.
	//
	// Format (array of integer between 0 (low) and 3 (critical) ):
	// [0, 1, 2, 3]
	level struct{}
)

func init() {
	registerFilter("level", level{})
}

// valid verifies the validity of the data
func (f level) valid(data interface{}) error {
	return genericValidUnsignedIntegerArray(f, data, 0, 3)
}

// apply applies the filter to the anomaly and returns the result.
func (f level) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.([]interface{}); !ok {
	} else {
		for _, level := range typed {
			if float64(an.Level) == level {
				return false
			}
		}
		return true
	}
	return false
}
