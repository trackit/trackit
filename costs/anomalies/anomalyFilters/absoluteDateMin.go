package anomalyFilters

import (
	"time"

	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
)

type (
	// absoluteDateMin will hide every entry before
	// the given date.
	//
	// Format (string):
	// 2006-01-02T15:04:05.000Z
	absoluteDateMin struct{}
)

func init() {
	registerFilter("absolute_date_min", absoluteDateMin{})
}

// valid verifies the validity of the data
func (f absoluteDateMin) valid(data interface{}) error {
	return genericValidDate(f, data)
}

// apply applies the filter to the anomaly and returns the result.
func (f absoluteDateMin) apply(data interface{}, an anomalyType.ProductAnomaly, product string) bool {
	if typed, ok := data.(string); !ok {
	} else if date, err := time.Parse("2006-01-02T15:04:05.000Z", typed); err != nil {
	} else if an.Date.Before(date) {
		return true
	}
	return false
}
