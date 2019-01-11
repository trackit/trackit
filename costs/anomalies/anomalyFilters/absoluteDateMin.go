package anomalyFilters

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

// apply applies the filter to the anomaly results
func (f absoluteDateMin) apply() {
}