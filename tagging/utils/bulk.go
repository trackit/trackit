//   Copyright 2020 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/trackit/trackit/es/indexes/taggingReports"
)

// GenerateBulkID generates a unique ID for a Tagging Report
func GenerateBulkID(doc taggingReports.TaggingReportDocument) (string, error) {
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
