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

package ebsReports

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

var Model = common.VersioningData{
	IndexSuffix:       "ebs-reports",
	Name:              "ebs-reports",
	Type:              "ebs-report",
	Version:           2,
	MappingProperties: properties,
}

type (
	// SnapshotReport is saved in ES to have all the information of an EBS snapshot
	SnapshotReport struct {
		common.ReportBase
		Snapshot Snapshot `json:"snapshot"`
	}

	// SnapshotBase contains basics information of an EBS snapshot
	SnapshotBase struct {
		Id          string    `json:"id"`
		Description string    `json:"description"`
		State       string    `json:"state"`
		Encrypted   bool      `json:"encrypted"`
		StartTime   time.Time `json:"startTime"`
		Region      string    `json:"region"`
	}

	// Snapshot contains all the information of an EBS snapshot
	Snapshot struct {
		SnapshotBase
		Tags   []common.Tag `json:"tags"`
		Volume Volume       `json:"volume"`
		Cost   float64      `json:"cost"`
	}

	// Volume contains information about an EBS volume
	Volume struct {
		Id   string `json:"id"`
		Size int64  `json:"size"`
	}
)

const properties = `
{
	"account": {
		"type": "keyword"
	},
	"reportDate": {
		"type": "date"
	},
	"reportType": {
		"type": "keyword"
	},
	"snapshot": {
		"properties": {
			"id": {
				"type": "keyword"
			},
			"description": {
				"type": "keyword"
			},
			"state": {
				"type": "keyword"
			},
			"encrypted": {
				"type": "boolean"
			},
			"startTime": {
				"type": "date"
			},
			"region": {
				"type": "keyword"
			},
			"tags": {
				"type": "nested",
				"properties": {
					"key": {
						"type": "keyword"
					},
					"value": {
						"type": "keyword"
					}
				}
			},
			"volume": {
				"type": "object",
				"properties": {
					"id": {
						"type": "keyword"
					},
					"size": {
						"type": "integer"
					}
				}
			},
			"cost": {
				"type": "double"
			}
		}
	}
}
`
