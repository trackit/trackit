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

package s3

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeS3Report = "s3-report"
const IndexPrefixS3Report = "s3-reports"
const TemplateNameS3Report = "s3-reports"

// put the ElasticSearch index for *-s3-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	res, err := es.Client.IndexPutTemplate(TemplateNameS3Report).BodyString(TemplateS3Report).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index s3-reports.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index s3-reports.", res)
	}
}

const TemplateS3Report = `
{
	"template": "*-s3-reports",
	"version": 2,
	"mappings": {
		"s3-report": {
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
				"bucket": {
					"properties": {
						"name": {
							"type": "keyword"
						},
						"creationDate": {
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
