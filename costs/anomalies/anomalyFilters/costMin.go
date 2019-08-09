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
