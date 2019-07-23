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

package anomalies

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
)

const TypeProductAnomaliesDetection = "product-anomalies-detection"
const IndexPrefixAnomaliesDetection = "anomalies-detection"
const TemplateNameAnomaliesDetection = "anomalies-detection"

// put the ElasticSearch index for *-anomalies-detection indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameAnomaliesDetection).BodyString(TemplateAnomaliesDetection).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index AnomaliesDetection.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index AnomaliesDetection.", res)
		ctxCancel()
	}
}

const TemplateAnomaliesDetection = `
{
	"template": "*-` + IndexPrefixAnomaliesDetection + `",
	"version": 2,
	"mappings": {
		"` + TypeProductAnomaliesDetection + `": {
			"properties": {
				"account": {
					"type": "keyword"
				},
				"date": {
					"type": "date"
				},
				"product" : {
					"type": "keyword"
				},
				"abnormal" : {
					"type": "boolean"
				},
				"recurrent" : {
					"type": "boolean"
				},
				"cost": {
					"type": "object",
					"properties": {
						"value": {
							"type": "double"
						},
						"maxExpected": {
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
