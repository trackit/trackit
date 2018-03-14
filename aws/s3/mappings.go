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

package s3

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/es"
)

const TypeLineItem = "lineitem"
const IndexPrefixLineItem = "lineitems"
const TemplateNameLineItem = "lineitems"

// put the ElasticSearch index for *-lineitems indices at startup.
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
	"template": "*-lineitems",
	"version": 4,
	"mappings": {
		"lineitem": {
			"properties": {
				"billRepositoryId": {
					"type": "integer"
				},
				"lineItemId": {
					"type": "keyword",
					"norms": false
				},
				"timeInterval": {
					"type": "keyword",
					"norms": false
				},
				"invoiceId": {
					"type": "keyword",
					"norms": false
				},
				"usageAccountId": {
					"type": "keyword",
					"norms": false
				},
				"productCode": {
					"type": "keyword",
					"norms": false
				},
				"usageType": {
					"type": "keyword",
					"norms": false
				},
				"operation": {
					"type": "keyword",
					"norms": false
				},
				"availabilityZone": {
					"type": "keyword",
					"norms": false
				},
				"resourceId": {
					"type": "keyword",
					"norms": false
				},
				"usageAmount": {
					"type": "float",
					"index": false
				},
				"serviceCode": {
					"type": "keyword",
					"norms": false
				},
				"currencyCode": {
					"type": "keyword",
					"norms": false
				},
				"unblendedCost": {
					"type": "float",
					"index": false
				},
				"usageStartDate": {
					"type": "date"
				},
				"usageEndDate": {
					"type": "date"
				}
			},
			"dynamic_templates": [
				{
					"tags": {
						"match_mapping_type": "string",
						"path_match": "tags.*",
						"mapping": {
							"type": "keyword"
						}
					}
				}
			],
			"_all": {
				"enabled": false
			},
			"numeric_detection": false,
			"date_detection": false
		}
	}
}
`
