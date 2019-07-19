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

package ec2

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeEC2Report = "ec2-report"
const IndexPrefixEC2Report = "ec2-reports"
const TemplateNameEC2Report = "ec2-reports"

// put the ElasticSearch index for *-ec2-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameEC2Report).BodyString(TemplateEc2Report).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index EC2Report.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index EC2Report.", res)
		ctxCancel()
	}
}

const TemplateEc2Report = `
{
	"template": "*-ec2-reports",
	"version": 8,
	"mappings": {
		"ec2-report": {
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
				"instance": {
					"properties": {
						"id": {
							"type": "keyword"
						},
						"region": {
							"type": "keyword"
						},
						"state": {
							"type": "keyword"
						},
						"purchasing": {
							"type": "keyword"
						},
						"keyPair": {
							"type": "keyword"
						},
						"type": {
							"type": "keyword"
						},
						"platform": {
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
								"network": {
									"type": "object",
									"properties": {
											"in": {
												"type": "double"
											},
											"out": {
												"type": "double"
											}
									}
								},
								"volumes": {
									"type": "nested",
									"properties": {
										"id": {
											"type": "keyword"
										},
										"read": {
											"type": "double"
										},
										"write": {
											"type": "double"
										}
									}
								}
							}
						},
						"recommendation": {
							"instancetype": {
								"type": "keyword"
							},
							"reason": {
								"type": "keyword"
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
