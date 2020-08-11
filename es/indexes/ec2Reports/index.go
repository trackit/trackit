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

package ec2Reports

import "github.com/trackit/trackit/es/indexes/common"

var Model = common.VersioningData{
	IndexSuffix:       "ec2-reports",
	Name:              "ec2-reports",
	Type:              "ec2-report",
	Version:           11,
	MappingProperties: properties,
}

type (
	// InstanceReport is saved in ES to have all the information of an EC2 instance
	InstanceReport struct {
		common.ReportBase
		Instance Instance `json:"instance"`
	}

	// InstanceBase contains basics information of an EC2 instance
	InstanceBase struct {
		Id         string `json:"id"`
		Region     string `json:"region"`
		State      string `json:"state"`
		Purchasing string `json:"purchasing"`
		KeyPair    string `json:"keyPair"`
		Type       string `json:"type"`
		Platform   string `json:"platform"`
	}

	// Instance contains all the information of an EC2 instance
	Instance struct {
		InstanceBase
		Tags           []common.Tag       `json:"tags"`
		Costs          map[string]float64 `json:"costs"`
		Stats          Stats              `json:"stats"`
		Recommendation Recommendation     `json:"recommendation"`
	}

	// Recommendation contains all recommendation of an EC2 instance
	Recommendation struct {
		InstanceType  string `json:"instancetype"`
		Reason        string `json:"reason"`
		NewGeneration string `json:"newgeneration"`
	}

	// Stats contains statistics of an instance get on CloudWatch
	Stats struct {
		Cpu     Cpu      `json:"cpu"`
		Network Network  `json:"network"`
		Volumes []Volume `json:"volumes"`
	}

	// Cpu contains cpu statistics of an instance
	Cpu struct {
		Average float64 `json:"average"`
		Peak    float64 `json:"peak"`
	}

	// Network contains network statistics of an instance
	Network struct {
		In  float64 `json:"in"`
		Out float64 `json:"out"`
	}

	// Volume contains information about an EBS volume
	Volume struct {
		Id    string  `json:"id"`
		Read  float64 `json:"read"`
		Write float64 `json:"write"`
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
		"properties": {
			"id": {
				"type": "keyword"
			},
			"region": {
				"type": "keyword"
			},
			"state": {
				"type": "keyword"
			},
			"purchasing": {
				"type": "keyword"
			},
			"keyPair": {
				"type": "keyword"
			},
			"type": {
				"type": "keyword"
			},
			"platform": {
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
					"network": {
						"type": "object",
						"properties": {
								"in": {
									"type": "double"
								},
								"out": {
									"type": "double"
								}
						}
					},
					"volumes": {
						"type": "nested",
						"properties": {
							"id": {
								"type": "keyword"
							},
							"read": {
								"type": "double"
							},
							"write": {
								"type": "double"
							}
						}
					}
				}
			},
			"recommendation": {
				"type": "object",
				"properties": {
					"instancetype": {
						"type": "keyword"
					},
					"reason": {
						"type": "keyword"
					},
					"newgeneration": {
						"type": "keyword"
					}
				}
			}
		}
	}
}
`
