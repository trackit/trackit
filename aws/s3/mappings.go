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
	"version": 2,
	"mappings": {
		"lineitem": {
			"properties": {
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
						"path_match": "tag.*",
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
