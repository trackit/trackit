package anomalyType

import (
	"time"
)

type (
	// AnomalyEsQueryParams will store the parsed query params
	AnomalyEsQueryParams struct {
		DateBegin   time.Time
		DateEnd     time.Time
		AccountList []string
		IndexList   []string
		AnomalyType string
	}

	// ProductAnomaly represents one anomaly returned.
	ProductAnomaly struct {
		Id          string    `json:"id"`
		Date        time.Time `json:"date"`
		Cost        float64   `json:"cost"`
		UpperBand   float64   `json:"upper_band"`
		Abnormal    bool      `json:"abnormal"`
		Filtered    bool      `json:"filtered"`
		Snoozed     bool      `json:"snoozed"`
		Level       int       `json:"level"`
		PrettyLevel string    `json:"pretty_level"`
	}

	// ProductAnomalies is used to respond to the request.
	// Key is a product name.
	ProductAnomalies map[string][]ProductAnomaly

	// AnomaliesDetectionResponse is used to respond to the request.
	// Key is an AWS Account Identity.
	AnomaliesDetectionResponse map[string]ProductAnomalies

	// Filter represents a filter.
	// A filter contains the rule and the associated data.
	Filter struct {
		Name     string      `json:"name" req:"nonzero"`
		Desc     string      `json:"desc"`
		Disabled bool        `json:"disabled" req:"nonzero"`
		Rule     string      `json:"rule" req:"nonzero"`
		Data     interface{} `json:"data" req:"nonzero"`
	}

	// Filters represents an array of filter.
	Filters []Filter
)
