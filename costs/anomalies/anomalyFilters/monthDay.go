package anomalyFilters

import (
	"github.com/trackit/trackit/costs/anomalies/anomalyType"
)

type (
	// monthDay will only show entries dated
	// the given month days.
	//
	// Format (array of integer between 1 (the 1st) and 31 (the 31st) ):
	// [1, 2, 25, 31]
	monthDay struct{}
)

func init() {
	registerFilter("month_day", monthDay{})
}

// valid verifies the validity of the data
func (f monthDay) valid(data interface{}) error {
	return genericValidUnsignedIntegerArray(f, data, 1, 31)
}

// apply applies the filter to the anomaly and returns the result.
func (f monthDay) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.([]interface{}); !ok {
	} else {
		for _, day := range typed {
			if float64(an.Date.Day()) == day {
				return false
			}
		}
		return true
	}
	return false
}
