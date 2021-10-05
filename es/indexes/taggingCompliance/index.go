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

package taggingCompliance

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

var Model = common.VersioningData{
	IndexSuffix:       "tagging-compliance",
	Name:              "tagging-compliance",
	Type:              "tagging-compliance",
	Version:           2,
	MappingProperties: properties,
}

type ComplianceReport struct {
	ReportDate      time.Time `json:"reportDate"`
	Total           int64     `json:"total"`
	TotallyTagged   int64     `json:"totallyTagged"`
	PartiallyTagged int64     `json:"partiallyTagged"`
	NotTagged       int64     `json:"notTagged"`
	MostUsedTagsId  string    `json:"mostUsedTagsId"`
	MostUsedTags    []string  `json:"mostUsedTags"`
}

const properties = `
{
	"reportDate":{
		"type":"date"
	},
	"total":{
		"type":"long"
	},
	"totallyTagged":{
		"type":"long"
	},
	"partiallyTagged":{
		"type":"long"
	},
	"notTagged":{
		"type":"long"
	},
	"mostUsedTagsId":{
		"type":"keyword"
	},
	"mostUsedTags":{
		"type":"text"
	}
}
`
