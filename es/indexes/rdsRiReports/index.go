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

package rdsRiReports

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

var Model = common.VersioningData{
	IndexSuffix:       "rds-ri-reports",
	Name:              "rds-ri-reports",
	Type:              "rds-ri-report",
	Version:           2,
	MappingProperties: properties,
}

type (
	// InstanceReport is saved in ES to have all the information of an RDS reserved instance
	InstanceReport struct {
		common.ReportBase
		Instance Instance `json:"instance"`
	}

	// InstanceBase contains basics information of an RDS reserved instance
	InstanceBase struct {
		DBInstanceIdentifier string             `json:"id"`
		DBInstanceOfferingId string             `json:"offeringId"`
		AvailabilityZone     string             `json:"availabilityZone"`
		DBInstanceClass      string             `json:"type"`
		DBInstanceCount      int64              `json:"dbInstanceCount"`
		Duration             int64              `json:"duration"`
		MultiAZ              bool               `json:"multiAZ"`
		ProductDescription   string             `json:"productDescription"`
		OfferingType         string             `json:"offeringType"`
		State                string             `json:"state"`
		StartTime            time.Time          `json:"startTime"`
		RecurringCharges     []RecurringCharges `json:"recurringCharges"`
	}

	// Instance contains the information of an RDS reserved instance
	Instance struct {
		InstanceBase
		Tags []common.Tag `json:"tags"`
	}

	//RecurringCharges contains recurring charges informations of a reservation
	RecurringCharges struct {
		Amount    float64
		Frequency string
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
			"offeringId": {
				"type": "keyword"
			},
			"availabilityZone": {
				"type": "keyword"
			},
			"instanceClass": {
				"type": "keyword"
			},
			"instanceCount": {
				"type": "integer"
			},
			"duration": {
				"type": "integer"
			},
			"multiAz": {
				"type": "boolean"
			},
			"productDescription": {
				"type": "keyword"
			},
			"offeringType": {
				"type": "keyword"
			},
			"state": {
				"type": "keyword"
			},
			"startTime": {
				"type": "date"
			},
			"recurringCharges": {
				"type": "nested",
				"properties": {
					"amount": {
						"type": "double"
					},
					"frequency": {
						"type": "keyword"
					}
				}
			}
		}
	}
}
`
