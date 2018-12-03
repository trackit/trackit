//   Copyright 2018 MSolution.IO
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

package reservedInstances

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
)

const TypeReservedInstancesReport = "reserved-instances-report"
const IndexPrefixReservedInstancesReport = "reserved-instances-reports"
const TemplateNameReservedInstancesReport = "reserved-instances-reports"

// put the ElasticSearch index for *-reserved-instances-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameReservedInstancesReport).BodyString(TemplateLineItem).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index ReservedInstancesReport.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index ReservedInstancesReport.", res)
		ctxCancel()
	}
}

const TemplateLineItem = `
{
	"template": "*-reserved-instances-reports",
	"version": 1,
	"mappings": {
		"reservedInstances-report": {
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
				"service": {
					"type": "keyword"
				},
				"reservation": {
					"properties": {
						"id": {
							"type": "keyword"
						},
						"region": {
							"type": "keyword"
						},
						"type": {
							"type": "keyword"
						},
						"fixedPrice": {
							"type": "keyword"
						},
						"usagePrice": {
							"type": "keyword"
						},
						"usageDuration": {
							"type": "keyword"
						},
						"start": {
							"type": "keyword"
						},
						"end": {
							"type": "keyword"
						},
						"instanceCount": {
							"type": "keyword"
						},
						"instanceTenancy": {
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
