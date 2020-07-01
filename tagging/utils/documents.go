package utils

import "time"

// TagDocument is a key/value pair which represents a tag and it's value
type TagDocument struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TaggingReportDocument is an entry in ES' tagging index
type TaggingReportDocument struct {
	Account      string        `json:"account"`
	ReportDate   time.Time     `json:"reportDate"`
	ResourceID   string        `json:"resourceId"`
	ResourceType string        `json:"resourceType"`
	Region       string        `json:"region"`
	URL          string        `json:"url"`
	Tags         []TagDocument `json:"tags"`
}
