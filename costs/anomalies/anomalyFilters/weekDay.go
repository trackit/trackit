package anomalyFilters

import "github.com/trackit/trackit-server/costs/anomalies/anomalyType"

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
	return genericValidUnsignedIntegerArray(f, data, 6)
}

// apply applies the filter to the anomaly results
func (f weekDay) apply(data interface{}, res anomalyType.AnomaliesDetectionResponse) anomalyType.AnomaliesDetectionResponse {
	return res
}
