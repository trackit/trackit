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
	"strings"

	"github.com/trackit/trackit/costs/anomalies/anomalyType"
)

type (
	// product will only show entries in the
	// given string array.
	//
	// Format (array of string):
	// ["AmazonEC2", "AmazonES"]
	product struct{}
)

func init() {
	registerFilter("product", product{})
}

// valid verifies the validity of the data
func (f product) valid(data interface{}) error {
	if typed, ok := data.([]interface{}); !ok {
		return fmt.Errorf("%s: not an array", filtersName[f])
	} else if len(typed) == 0 {
		return fmt.Errorf("%s: empty array", filtersName[f])
	} else {
		for i := range typed {
			if _, ok := typed[i].(string); !ok {
				return fmt.Errorf("%s: not an array of string", filtersName[f])
			}
		}
	}
	return nil
}

// apply applies the filter to the anomaly and returns the result.
func (f product) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.([]interface{}); !ok {
	} else {
		for _, p := range typed {
			if ps, ok := p.(string); ok && strings.Contains(strings.ToLower(product), strings.ToLower(ps)) {
				return false
			}
		}
		return true
	}
	return false
}
