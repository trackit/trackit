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

package taggingReports

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

var Model = common.VersioningData{
	IndexSuffix:       "tagging-reports",
	Name:              "tagging-reports",
	Type:              "tagging-reports",
	Version:           1,
	MappingProperties: properties,
}

// TaggingReportDocument is an entry in ES' tagging index
type TaggingReportDocument struct {
	Account      string       `json:"account"`
	ReportDate   time.Time    `json:"reportDate"`
	ResourceID   string       `json:"resourceId"`
	ResourceType string       `json:"resourceType"`
	Region       string       `json:"region"`
	URL          string       `json:"url"`
	Tags         []common.Tag `json:"tags"`
}

const properties = `
{
	"account":{
		"type":"keyword"
	},
	"region":{
		"type":"keyword"
	},
	"reportDate":{
		"type":"date"
	},
	"resourceId":{
		"type":"keyword"
	},
	"resourceType":{
		"type":"keyword"
	},
	"tags":{
		"type":"nested",
		"properties":{
			"key":{
				"type":"keyword"
			},
			"value":{
				"type":"keyword"
			}
		}
	},
	"url":{
		"type":"keyword"
	}
}
`
