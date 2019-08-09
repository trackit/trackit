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

package mediastore

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeMediaStoreReport = "mediastore-report"
const IndexPrefixMediaStoreReport = "mediastore-reports"
const TemplateNameMediaStoreReport = "mediastore-reports"

// put the ElasticSearch index for *-mediastore-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameMediaStoreReport).BodyString(TemplateMediaStoreReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index MediaStoreReport.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index MediaStoreReport.", res)
		ctxCancel()
	}
}

const TemplateMediaStoreReport = `
{
	"template": "*-mediastore-reports",
	"version": 7,
	"mappings": {
		"mediastore-report": {
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
				"container": {
					"properties": {
						"arn": {
							"type": "keyword"
						},
						"region": {
							"type": "keyword"
						},
						"name": {
							"type": "keyword"
						},
						"costs": {
							"type": "nested",
							"properties": {
								"key": {
									"type": "date"
								},
								"value": {
									"type": "double"
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
