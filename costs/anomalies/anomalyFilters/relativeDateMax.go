package anomalyFilters

import "github.com/trackit/trackit-server/costs/anomalies/anomalyType"

type (
	// relativeDateMax will hide every entry after
	// today minus the given duration.
	//
	// Format (positive integer, in seconds):
	// 3600
	relativeDateMax struct{}
)

func init() {
	registerFilter("relative_date_max", relativeDateMax{})
}

// valid verifies the validity of the data
func (f relativeDateMax) valid(data interface{}) error {
	return genericValidUnsignedInteger(f, data)
}

// apply applies the filter to the anomaly and returns the result.
func (f relativeDateMax) apply(data interface{}, res anomalyType.ProductAnomaly) bool {
	return false
}
