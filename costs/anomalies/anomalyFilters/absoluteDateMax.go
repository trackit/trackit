package anomalyFilters

import "github.com/trackit/trackit-server/costs/anomalies/anomalyType"

type (
	// absoluteDateMax will hide every entry after
	// the given date.
	//
	// Format (string):
	// 2006-01-02T15:04:05.000Z
	absoluteDateMax struct{}
)

func init() {
	registerFilter("absolute_date_max", absoluteDateMax{})
}

// valid verifies the validity of the data
func (f absoluteDateMax) valid(data interface{}) error {
	return genericValidDate(f, data)
}

// apply applies the filter to the anomaly and returns the result.
func (f absoluteDateMax) apply(data interface{}, res anomalyType.ProductAnomaly) bool {
	return false
}
