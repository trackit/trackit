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

package ec2CoverageReports

import "github.com/trackit/trackit/es/indexes/common"

var Model = common.VersioningData{
	IndexSuffix:       "ec2-coverage-reports",
	Name:              "ec2-coverage-reports",
	Type:              "ec2-coverage-report",
	Version:           1,
	MappingProperties: properties,
}

type (
	// ReservationReport is saved in ES to have all the information of an EC2 reservation
	ReservationReport struct {
		common.ReportBase
		Reservation Reservation `json:"reservation"`
	}

	// Reservation contains basics information of an EC2 reservation
	Reservation struct {
		Type              string   `json:"type"`
		Platform          string   `json:"platform"`
		Tenancy           string   `json:"tenancy"`
		Region            string   `json:"region"`
		AverageCoverage   float64  `json:"averageCoverage"`
		CoveredHours      float64  `json:"coveredHours"`
		OnDemandHours     float64  `json:"onDemandHours"`
		TotalRunningHours float64  `json:"totalRunningHours"`
		InstancesNames    []string `json:"instancesNames"`
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
	"reservation": {
		"properties": {
			"type": {
				"type": "keyword"
			},
			"platform": {
				"type": "keyword"
			},
			"tenancy": {
				"type": "keyword"
			},
			"region": {
				"type": "keyword"
			},
			"averageCoverage": {
				"type": "double"
			},
			"coveredHours": {
				"type": "double"
			},
			"onDemandHours": {
				"type": "double"
			},
			"totalRunningHours": {
				"type": "double"
			},
			"instancesNames": {
				"type": "keyword"
			}
		}
	}
}
`
