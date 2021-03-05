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

package stepfunction

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeStepFunctionReport = "stepfunction-report"
const IndexPrefixStepFunctionReport = "stepfunction-reports"
const TemplateNameStepFunctionReport = "stepfunction-reports"

// put the ElasticSearch index for *-stepfunction-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameStepFunctionReport).BodyString(TemplateEbsReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index stepfunction-reports.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index stepfunction-reports.", res)
		ctxCancel()
	}
}

const TemplateEbsReport = `
{
	"template": "*-stepfunction-reports",
	"version": 2,
	"mappings": {
		"stepfunction-report": {
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
				"step": {
					"properties": {
						"name": {
							"type": "keyword"
						},
						"creationDate": {
							"type": "date"
						},
						"arn": {
							"type": "keyword"
						},
						"region": {
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
