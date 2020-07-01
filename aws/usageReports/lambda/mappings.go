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

package lambda

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeLambdaReport = "lambda-report"
const IndexPrefixLambdaReport = "lambda-reports"
const TemplateNameLambdaReport = "lambda-reports"

// put the ElasticSearch index for *-lambda-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameLambdaReport).BodyString(TemplateLineItem).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index LambdaReport.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index LambdaReport.", res)
		ctxCancel()
	}
}

const TemplateLineItem = `
{
	"template": "*-lambda-reports",
	"version": 2,
	"mappings": {
		"lambda-report": {
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
				"function": {
					"properties": {
						"name": {
							"type": "keyword"
						},
						"description": {
							"type": "keyword"
						},
						"size": {
							"type": "integer"
						},
						"memory": {
							"type": "integer"
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
						},
						"stats": {
							"type": "object",
							"properties": {
								"invocations": {
									"type": "object",
									"properties": {
											"total": {
												"type": "double"
											},
											"error": {
												"type": "double"
											}
									}
								},
								"duration": {
									"type": "object",
									"properties": {
											"average": {
												"type": "double"
											},
											"maximum": {
												"type": "double"
											}
									}
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
