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

const Template = `
{
	"template": "*-` + IndexSuffix + `",
	"version": 3,
	"mappings": ` + Mappings + `
}
`

const Mappings = `
{
	"ri-ec2-report": {
		"properties": {
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
		},
		"_all": {
			"enabled": false
		},
		"numeric_detection": false,
		"date_detection": false
	}
}
`
