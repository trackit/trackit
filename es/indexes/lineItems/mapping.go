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

package lineItems

const Template = `
{
	"template": "*-` + IndexSuffix + `",
	"version": 8,
	"mappings": ` + Mappings + `
}
`

const Mappings = `
{
	"` + Type + `": {
		"properties": {
			"billRepositoryId": {
				"type": "integer"
			},
			"lineItemId": {
				"type": "keyword",
				"norms": false
			},
			"timeInterval": {
				"type": "keyword",
				"norms": false
			},
			"invoiceId": {
				"type": "keyword",
				"norms": false
			},
			"usageAccountId": {
				"type": "keyword",
				"norms": false
			},
			"lineItemType": {
				"type": "keyword",
				"norms": false
			},
			"productCode": {
				"type": "keyword",
				"norms": false
			},
			"usageType": {
				"type": "keyword",
				"norms": false
			},
			"operation": {
				"type": "keyword",
				"norms": false
			},
			"availabilityZone": {
				"type": "keyword",
				"norms": false
			},
			"region": {
				"type": "keyword",
				"norms": false
			},
			"resourceId": {
				"type": "keyword",
				"norms": false
			},
			"usageAmount": {
				"type": "float",
				"index": false
			},
			"serviceCode": {
				"type": "keyword",
				"norms": false
			},
			"currencyCode": {
				"type": "keyword",
				"norms": false
			},
			"unblendedCost": {
				"type": "float",
				"index": false
			},
			"taxType": {
				"type": "keyword",
				"norms": false
			},
			"usageStartDate": {
				"type": "date"
			},
			"usageEndDate": {
				"type": "date"
			},
			"tags": {
				"type": "nested",
				"properties": {
					"key": {
						"type": "keyword",
						"norms": false
					},
					"tag": {
						"type": "keyword",
						"norms": false
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
