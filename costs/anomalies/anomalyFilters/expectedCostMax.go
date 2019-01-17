package anomalyFilters

import (
	"fmt"

	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
)

type (
	// costMax will hide every entry whose
	// the expected cost is greater than
	// the given value.
	//
	// Format (number):
	// 514.5
	expectedCostMax struct{}
)

func init() {
	registerFilter("expected_cost_max", expectedCostMax{})
}

// valid verifies the validity of the data
func (f expectedCostMax) valid(data interface{}) error {
	if _, ok := data.(float64); !ok {
		return fmt.Errorf("%s: not a number", filtersName[f])
	}
	return nil
}

// apply applies the filter to the anomaly and returns the result.
func (f expectedCostMax) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.(float64); !ok {
	} else if an.UpperBand > typed {
		return true
	}
	return false
}
