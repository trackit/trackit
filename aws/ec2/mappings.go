//   Copyright 2017 MSolution.IO
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

package ec2

const IndexHourlyInstanceUsage = "awshourlyinstanceusage"
const DocumentEc2 = "ec2instance"

const mappingEc2Instance = `{
	"properties": {
		"account": {
			"type": "keyword"
		},
		"service": {
			"type": "keyword"
		},
		"id": {
			"type": "keyword"
		},
		"region": {
			"type": "keyword"
		},
		"startDate": {
			"type": "date"
		},
		"endDate": {
			"type": "date"
		},
		"cpuAverage": {
			"type": "double"
		},
		"cpuPeak": {
			"type": "double"
		},
		"keyPair": {
			"type": "keyword"
		},
		"type": {
			"type": "keyword"
		},
		"tags": {
			"type": "nested"
		}
	},
	"_all": {
		"enabled": false
	},
	"numeric_detection": false,
	"date_detection": false
}`
