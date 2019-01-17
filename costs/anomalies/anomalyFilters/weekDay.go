package anomalyFilters

import (
	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
)

type (
	// weekDay will only show entries dated
	// the given week days.
	//
	// Format (array of integer between 0 (Monday) and 6 (Sunday) ):
	// [0, 1, 2, 3, 4]
	weekDay struct{}
)

func init() {
	registerFilter("week_day", weekDay{})
}

// valid verifies the validity of the data
func (f weekDay) valid(data interface{}) error {
	return genericValidUnsignedIntegerArray(f, data, 0, 6)
}

// apply applies the filter to the anomaly and returns the result.
func (f weekDay) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.([]interface{}); !ok {
	} else {
		for _, day := range typed {
			if float64(an.Date.Weekday())-1 == day {
				return false
			}
		}
		return true
	}
	return false
}
