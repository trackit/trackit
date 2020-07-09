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

package ebs

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeEBSReport = "ebs-report"
const IndexPrefixEBSReport = "ebs-reports"
const TemplateNameEBSReport = "ebs-reports"

// put the ElasticSearch index for *-ebs-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameEBSReport).BodyString(TemplateEbsReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index EBSReport.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index EBSReport.", res)
		ctxCancel()
	}
}

const TemplateEbsReport = `
{
	"index_patterns": ["*-ebs-reports"],
	"version": 3,
	"mappings": {
		"ebs-report": {
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
				"snapshot": {
					"properties": {
						"id": {
							"type": "keyword"
						},
						"description": {
							"type": "keyword"
						},
						"state": {
							"type": "keyword"
						},
						"encrypted": {
							"type": "boolean"
						},
						"startTime": {
							"type": "date"
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
						"volume": {
							"type": "object",
							"properties": {
								"id": {
									"type": "keyword"
								},
								"size": {
									"type": "integer"
								}
							}
						},
						"cost": {
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
}
`
