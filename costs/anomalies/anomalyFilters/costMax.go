package anomalyFilters

import (
	"fmt"

	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
)

type (
	// costMax will hide every entry whose
	// the cost is greater than the given value.
	//
	// Format (number):
	// 514.5
	costMax struct{}
)

func init() {
	registerFilter("cost_max", costMax{})
}

// valid verifies the validity of the data
func (f costMax) valid(data interface{}) error {
	if _, ok := data.(float64); !ok {
		return fmt.Errorf("%s: not a number", filtersName[f])
	}
	return nil
}

// apply applies the filter to the anomaly and returns the result.
func (f costMax) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.(float64); !ok {
	} else if an.Cost > typed {
		return true
	}
	return false
}
