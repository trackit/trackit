//   Copyright 2019 MSolution.IO
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

package ec2Coverage

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
)

const TypeEC2CoverageReport = "ec2-coverage-report"
const IndexPrefixEC2CoverageReport = "ec2-coverage-reports"
const TemplateNameEC2CoverageReport = "ec2-coverage-reports"

// put the ElasticSearch index for *-ec2-coverage-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameEC2CoverageReport).BodyString(TemplateEc2CoverageReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index EC2 Coverage Report.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index EC2 Coverage Report.", res)
		ctxCancel()
	}
}

const TemplateEc2CoverageReport = `
{
	"template": "*-ec2-coverage-reports",
	"version": 1,
	"mappings": {
		"ec2-coverage-report": {
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
			},
			"_all": {
				"enabled": false
			},
			"numeric_detection": false,
			"date_detection": false
		}
	}
}
`
