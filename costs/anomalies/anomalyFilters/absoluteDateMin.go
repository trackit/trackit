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

// apply applies the filter to the anomaly results
func (f absoluteDateMin) apply(data interface{}, res anomalyType.AnomaliesDetectionResponse) anomalyType.AnomaliesDetectionResponse {
	if typed, ok := data.(string); !ok {
	} else if date, err := time.Parse("2006-01-02T15:04:05.000Z", typed); err != nil {
	} else {
		for account := range res {
			for product := range res[account] {
				for anomaly, an := range res[account][product] {
					if an.Abnormal && an.Date.Before(date) {
						res[account][product][anomaly].Filtered = true
					}
				}
			}
		}
	}
	return res
}
