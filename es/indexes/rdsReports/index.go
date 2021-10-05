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

package rdsReports

import "github.com/trackit/trackit/es/indexes/common"

var Model = common.VersioningData{
	IndexSuffix:       "rds-reports",
	Name:              "rds-reports",
	Type:              "rds-report",
	Version:           5,
	MappingProperties: properties,
}

type (
	// InstanceReport is saved in ES to have all the information of an RDS instance
	InstanceReport struct {
		common.ReportBase
		Instance Instance `json:"instance"`
	}

	// InstanceBase contains basics information of an RDS instance
	InstanceBase struct {
		DBInstanceIdentifier string `json:"id"`
		AvailabilityZone     string `json:"availabilityZone"`
		DBInstanceClass      string `json:"type"`
		Engine               string `json:"engine"`
		AllocatedStorage     int64  `json:"allocatedStorage"`
		MultiAZ              bool   `json:"multiAZ"`
	}

	// Instance contains the information of an RDS instance
	Instance struct {
		InstanceBase
		Tags  []common.Tag       `json:"tags"`
		Costs map[string]float64 `json:"costs"`
		Stats Stats              `json:"stats"`
	}

	// Stats contains statistics of an instance get on CloudWatch
	Stats struct {
		Cpu       Cpu       `json:"cpu"`
		FreeSpace FreeSpace `json:"freeSpace"`
	}

	// Cpu contains cpu statistics of an instance
	Cpu struct {
		Average float64 `json:"average"`
		Peak    float64 `json:"peak"`
	}

	// FreeSpace contains free space statistics of an instance
	FreeSpace struct {
		Minimum float64 `json:"minimum"`
		Maximum float64 `json:"maximum"`
		Average float64 `json:"average"`
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
	"instance": {
		"type": "object",
		"properties": {
			"id": {
				"type": "keyword"
			},
			"availabilityZone": {
				"type": "keyword"
			},
			"type": {
				"type": "keyword"
			},
			"engine": {
				"type": "keyword"
			},
			"allocatedStorage": {
				"type": "integer"
			},
			"multiAZ": {
				"type": "boolean"
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
			"costs": {
				"type": "object"
			},
			"stats": {
				"type": "object",
				"properties": {
					"cpu": {
						"type": "object",
						"properties": {
								"average": {
									"type": "double"
								},
								"peak": {
									"type": "double"
								}
						}
					},
					"freeSpace": {
						"type": "object",
						"properties": {
								"minimum": {
									"type": "double"
								},
								"maximum": {
									"type": "double"
								},
								"average": {
									"type": "double"
								}
						}
					}
				}
			}
		}
	}
}
`
