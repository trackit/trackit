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

package anomaliesDetection

const Template = `
{
	"template": "*-` + IndexSuffix + `",
	"version": 2,
	"mappings": ` + Mappings + `
}
`

const Mappings = `
{
	"product-anomalies-detection": {
		"properties": {
			"account": {
				"type": "keyword"
			},
			"date": {
				"type": "date"
			},
			"product" : {
				"type": "keyword"
			},
			"abnormal" : {
				"type": "boolean"
			},
			"recurrent" : {
				"type": "boolean"
			},
			"cost": {
				"type": "object",
				"properties": {
					"value": {
						"type": "double"
					},
					"maxExpected": {
						"type": "double"
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
