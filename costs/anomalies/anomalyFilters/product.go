package anomalyFilters

import (
	"fmt"

	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
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
func (f product) apply(data interface{}, res anomalyType.ProductAnomaly) bool {
	return false
}
