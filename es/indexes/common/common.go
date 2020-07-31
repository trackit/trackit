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

package common

import "time"

type (
	// CostPerResource associates a cost to an aws resourceId with a region
	CostPerResource struct {
		Resource string
		Cost     float64
		Region   string
	}

	// BaseReport contains basic information for any kin of usage report
	ReportBase struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		ReportType string    `json:"reportType"`
	}

	// Tag contains the key of a tag and his value
	Tag struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

type VersioningData struct {
	IndexSuffix string
	Name        string
	Version     int
	Template    string
}
