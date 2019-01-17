package anomalyFilters

import (
	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
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
