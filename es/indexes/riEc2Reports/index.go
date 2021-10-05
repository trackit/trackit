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

package riEc2Reports

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

var Model = common.VersioningData{
	IndexSuffix:       "ri-ec2-reports",
	Name:              "ri-ec2-reports",
	Type:              "ri-ec2-report",
	Version:           3,
	MappingProperties: properties,
}

type (
	// ReservationReport is saved in ES to have all the information of a reservation
	ReservationReport struct {
		common.ReportBase
		Reservation Reservation `json:"reservation"`
	}

	// ReservationBase contains basics information of a reserved instance
	ReservationBase struct {
		Id                 string             `json:"id"`
		Region             string             `json:"region"`
		AvailabilityZone   string             `json:"availabilityZone"`
		Type               string             `json:"type"`
		OfferingClass      string             `json:"offeringClass"`
		OfferingType       string             `json:"offeringType"`
		ProductDescription string             `json:"productDescription"`
		State              string             `json:"state"`
		Start              time.Time          `json:"start"`
		End                time.Time          `json:"end"`
		InstanceCount      int64              `json:"instanceCount"`
		Tenancy            string             `json:"tenancy"`
		UsagePrice         float64            `json:"usagePrice"`
		RecurringCharges   []RecurringCharges `json:"recurringCharges"`
	}

	// Reservation contains all the information of a reservation
	Reservation struct {
		ReservationBase
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
	"service": {
		"type": "keyword"
	},
	"reservation": {
		"properties": {
			"id": {
				"type": "keyword"
			},
			"region": {
				"type": "keyword"
			},
			"availabilityZone": {
				"type": "keyword"
			},
			"type": {
				"type": "keyword"
			},
			"offeringClass": {
				"type": "keyword"
			},
			"offeringType": {
				"type": "keyword"
			},
			"productDescription": {
				"type": "keyword"
			},
			"state":{
				"type": "keyword"
			},
			"start": {
				"type": "date"
			},
			"end": {
				"type": "date"
			},
			"instanceCount": {
				"type": "integer"
			},
			"tenancy": {
				"type": "keyword"
			},
			"usagePrice": {
				"type": "double"
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
			}
		}
	}
}
`
