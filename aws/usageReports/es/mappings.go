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

package es

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeESReport = "es-report"
const IndexPrefixESReport = "es-reports"
const TemplateNameESReport = "es-reports"

// put the ElasticSearch index for *-es-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameESReport).BodyString(TemplateEsReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index ESReport.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index ESReport.", res)
		ctxCancel()
	}
}

const TemplateEsReport = `
{
	"index_patterns": ["*-es-reports"],
	"version": 2,
	"mappings": {
		"es-report": {
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
				"domain": {
					"properties": {
						"arn": {
							"type": "keyword"
						},
						"region": {
							"type": "keyword"
						},
						"domainId": {
							"type": "keyword"
						},
						"instanceType": {
							"type": "keyword"
						},
						"instanceCount": {
							"type": "integer"
						},
						"totalStorageSpace": {
							"type": "integer"
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
						"costs": {
							"type": "object"
						},
						"stats": {
							"type": "object",
							"properties": {
								"cpu": {
									"type": "object",
									"properties": {
											"average": {
												"type": "double"
											},
											"peak": {
												"type": "double"
											}
									}
								},
								"freeSpace": {
									"type": "double"
								},
								"JVMMemoryPressure": {
									"type": "object",
									"properties": {
											"in": {
												"type": "double"
											},
											"out": {
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
