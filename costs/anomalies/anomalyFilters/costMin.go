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
	// costMin will hide every entry whose
	// the cost is lower than the given value.
	//
	// Format (number):
	// 514.5
	costMin struct{}
)

func init() {
	registerFilter("cost_min", costMin{})
}

// valid verifies the validity of the data
func (f costMin) valid(data interface{}) error {
	if _, ok := data.(float64); !ok {
		return fmt.Errorf("%s: not a number", filtersName[f])
	}
	return nil
}

// apply applies the filter to the anomaly and returns the result.
func (f costMin) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.(float64); !ok {
	} else if an.Cost < typed {
		return true
	}
	return false
}
