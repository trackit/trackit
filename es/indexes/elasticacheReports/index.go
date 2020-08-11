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

package elasticacheReports

import "github.com/trackit/trackit/es/indexes/common"

var Model = common.VersioningData{
	IndexSuffix:       "elasticache-reports",
	Name:              "elasticache-reports",
	Type:              "elasticache-report",
	Version:           1,
	MappingProperties: properties,
}

type (
	// InstanceReport is saved in ES to have all the information of an ElastiCache instance
	InstanceReport struct {
		common.ReportBase
		Instance Instance `json:"instance"`
	}

	// InstanceBase contains basics information of an ElastiCache instance
	InstanceBase struct {
		Id            string `json:"id"`
		Status        string `json:"status"`
		Region        string `json:"region"`
		NodeType      string `json:"nodeType"`
		Nodes         []Node `json:"nodes"`
		Engine        string `json:"engine"`
		EngineVersion string `json:"engineVersion"`
	}

	// Instance contains all the information of an ElastiCache instance
	Instance struct {
		InstanceBase
		Tags  []common.Tag       `json:"tags"`
		Costs map[string]float64 `json:"costs"`
		Stats Stats              `json:"stats"`
	}

	Node struct {
		Id     string `json:"id"`
		Status string `json:"status"`
		Region string `json:"region"`
	}

	// Stats contains statistics of an instance get on CloudWatch
	Stats struct {
		Cpu     Cpu     `json:"cpu"`
		Network Network `json:"network"`
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
			"status": {
				"type": "keyword"
			},
			"region": {
				"type": "keyword"
			},
			"nodeType": {
				"type": "keyword"
			},
			"nodes": {
				"type": "nested",
				"properties": {
					"id": {
						"type": "keyword"
					},
					"status": {
						"type": "keyword"
					},
					"region": {
						"type": "keyword"
					}
				}
			},
			"engine": {
				"type": "keyword"
			},
			"engineVersion": {
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
					}
				}
			}
		}
	}
}
`
