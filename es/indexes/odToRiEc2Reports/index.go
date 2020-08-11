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

package odToRiEc2Reports

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

var Model = common.VersioningData{
	IndexSuffix:       "od-to-ri-ec2-reports",
	Name:              "od-to-ri-ec2-reports",
	Type:              "od-to-ri-ec2-report",
	Version:           1,
	MappingProperties: properties,
}

type (
	Cost struct {
		PerUnit float64 `json:"perUnit"`
		Total   float64 `json:"total"`
	}

	OnDemandCost struct {
		Monthly    Cost `json:"monthly"`
		OneYear    Cost `json:"oneYear"`
		ThreeYears Cost `json:"threeYears"`
	}

	ReservationCost struct {
		Monthly Cost `json:"monthly"`
		Global  Cost `json:"global"`
		Saving  Cost `json:"saving"`
	}

	OnDemandTotalCost struct {
		MonthlyTotal    float64 `json:"monthly"`
		OneYearTotal    float64 `json:"oneYear"`
		ThreeYearsTotal float64 `json:"threeYears"`
	}

	ReservationTotalCost struct {
		MonthlyTotal float64 `json:"monthly"`
		GlobalTotal  float64 `json:"global"`
		SavingTotal  float64 `json:"saving"`
	}

	// InstancesSpecs stores the costs calculated for a given region/instance/platform
	// combination
	InstancesSpecs struct {
		Region        string       `json:"region"`
		Type          string       `json:"instanceType"`
		Platform      string       `json:"platform"`
		InstanceCount int          `json:"instanceCount"`
		OnDemand      OnDemandCost `json:"onDemand"`
		Reservation   struct {
			Type      string          `json:"type"`
			OneYear   ReservationCost `json:"oneYear"`
			ThreeYear ReservationCost `json:"threeYears"`
		} `json:"reservation"`
	}

	// OdToRiEc2Report stores all the on demand to RI EC2 report infos
	OdToRiEc2Report struct {
		Account     string            `json:"account"`
		ReportDate  time.Time         `json:"reportDate"`
		OnDemand    OnDemandTotalCost `json:"onDemand"`
		Reservation struct {
			OneYear   ReservationTotalCost `json:"oneYear"`
			ThreeYear ReservationTotalCost `json:"threeYears"`
		} `json:"reservation"`
		Instances []InstancesSpecs `json:"instances"`
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
	"onDemand": {
		"properties": {
			"monthly": {
				"type": "double"
			},
			"oneYear": {
				"type": "double"
			},
			"threeYears": {
				"type": "double"
			}
		}
	},
	"reservation": {
		"properties": {
			"oneYear": {
				"properties": {
					"monthly": {
						"type": "double"
					},
					"global": {
						"type": "double"
					},
					"saving": {
						"type": "double"
					}
				}
			},
			"threeYears": {
				"properties": {
					"monthly": {
						"type": "double"
					},
					"global": {
						"type": "double"
					},
					"saving": {
						"type": "double"
					}
				}
			}
		}
	},
	"instances": {
		"type": "nested",
		"properties": {
			"region": {
				"type": "keyword"
			},
			"instanceType": {
				"type": "keyword"
			},
			"platform": {
				"type": "keyword"
			},
			"instanceCount": {
				"type": "integer"
			},
			"onDemand": {
				"properties": {
					"monthly": {
						"properties": {
							"perUnit": {
								"type": "double"
							},
							"total": {
								"type": "double"
							}
						}
					},
					"oneYear": {
						"properties": {
							"perUnit": {
								"type": "double"
							},
							"total": {
								"type": "double"
							}
						}
					},
					"threeYears": {
						"properties": {
							"perUnit": {
								"type": "double"
							},
							"total": {
								"type": "double"
							}
						}
					}
				}
			},
			"reservation": {
				"properties": {
					"type": {
						"type": "keyword"
					},
					"oneYear": {
						"properties": {
							"monthly": {
								"properties": {
									"perUnit": {
										"type": "double"
									},
									"total": {
										"type": "double"
									}
								}
							},
							"global": {
								"properties": {
									"perUnit": {
										"type": "double"
									},
									"total": {
										"type": "double"
									}
								}
							},
							"saving": {
								"properties": {
									"perUnit": {
										"type": "double"
									},
									"total": {
										"type": "double"
									}
								}
							}
						}
					},
					"threeYears": {
						"properties": {
							"monthly": {
								"properties": {
									"perUnit": {
										"type": "double"
									},
									"total": {
										"type": "double"
									}
								}
							},
							"global": {
								"properties": {
									"perUnit": {
										"type": "double"
									},
									"total": {
										"type": "double"
									}
								}
							},
							"saving": {
								"properties": {
									"perUnit": {
										"type": "double"
									},
									"total": {
										"type": "double"
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
`
