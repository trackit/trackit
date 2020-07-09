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

package instanceCount

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeInstanceCountReport = "instancecount-report"
const IndexPrefixInstanceCountReport = "instancecount-reports"
const TemplateNameInstanceCountReport = "instancecount-reports"

// put the ElasticSearch index for *-instanceCount-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameInstanceCountReport).BodyString(TemplateInstanceCountReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index InstanceCountReport.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index InstanceCountReport.", res)
		ctxCancel()
	}
}

const TemplateInstanceCountReport = `
{
	"index_patterns": ["*-instancecount-reports"],
	"version": 2,
	"mappings": {
		"instancecount-report": {
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
				"instanceCount": {
					"properties": {
						"instanceType": {
							"type": "keyword"
						},
						"hours": {
							"type": "nested",
							"properties": {
								"hour": {
									"type": "date"
								},
								"count": {
									"type": "integer"
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
