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

import (
	"time"
	"context"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
)

const TypeLineItem = "ec2-report"
const IndexPrefixLineItem = "ec2-reports"
const TemplateNameLineItem = "ec2-reports"

// put the ElasticSearch index for *-ec2-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameLineItem).BodyString(TemplateLineItem).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index lineitems.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index lineitems.", res)
		ctxCancel()
	}
}

const TemplateLineItem = `
{
	"template": "*-ec2-reports",
	"version": 5,
	"mappings": {
		"ec2-report": {
			"properties": {
				"account": {
					"type": "keyword"
				},
				"reportDate": {
					"type": "date"
				},
				"endDate": {
					"type": "date"
				},
				"instances" : {
					"type": "nested"
					"properties": {
						"id": {
							"type": "keyword"
						},
						"region": {
							"type": "keyword"
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

