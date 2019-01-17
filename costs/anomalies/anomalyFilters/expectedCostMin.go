package anomalyFilters

import (
	"fmt"

	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
)

type (
	// expectedCostMin will hide every entry whose
	// the expected cost is lower than
	// the given value.
	//
	// Format (number):
	// 514.5
	expectedCostMin struct{}
)

func init() {
	registerFilter("expected_cost_min", expectedCostMin{})
}

// valid verifies the validity of the data
func (f expectedCostMin) valid(data interface{}) error {
	if _, ok := data.(float64); !ok {
		return fmt.Errorf("%s: not a number", filtersName[f])
	}
	return nil
}

// apply applies the filter to the anomaly and returns the result.
func (f expectedCostMin) apply(data interface{}, res anomalyType.ProductAnomaly) bool {
	return false
}
