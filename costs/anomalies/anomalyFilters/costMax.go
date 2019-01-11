package anomalyFilters

import "fmt"

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

// apply applies the filter to the anomaly results
func (f costMax) apply() {
}
