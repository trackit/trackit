package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"time"
)

// GenerateBulkID generates a unique ID for a Tagging Report
func GenerateBulkID(doc TaggingReportDocument) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		ID         string    `json:"id"`
	}{
		doc.Account,
		doc.ReportDate,
		doc.ResourceID,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}
