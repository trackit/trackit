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

package route53

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeRoute53Report = "route53-report"
const IndexPrefixRoute53Report = "route53-reports"
const TemplateNameRoute53Report = "route53-reports"

// put the ElasticSearch index for *-route53-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	res, err := es.Client.IndexPutTemplate(TemplateNameRoute53Report).BodyString(TemplateRoute53Report).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index route53-reports.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index route53-reports.", res)
	}
}

const TemplateRoute53Report = `
{
	"template": "*-route53-reports",
	"version": 2,
	"mappings": {
		"route53-report": {
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
				"hostedZone": {
					"properties": {
						"name": {
							"type": "keyword"
						},
						"id": {
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
