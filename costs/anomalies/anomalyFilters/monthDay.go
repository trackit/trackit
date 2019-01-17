package anomalyFilters

import "github.com/trackit/trackit-server/costs/anomalies/anomalyType"

type (
	// monthDay will only show entries dated
	// the given month days.
	//
	// Format (array of integer between 0 (the 1st) and 30 (the 31st) ):
	// [0, 1, 2, 25, 30]
	monthDay struct{}
)

func init() {
	registerFilter("month_day", monthDay{})
}

// valid verifies the validity of the data
func (f monthDay) valid(data interface{}) error {
	return genericValidUnsignedIntegerArray(f, data, 30)
}

// apply applies the filter to the anomaly results
func (f monthDay) apply(data interface{}, res anomalyType.AnomaliesDetectionResponse) anomalyType.AnomaliesDetectionResponse {
	return res
}
