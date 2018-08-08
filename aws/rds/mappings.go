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

package rds

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
)

const TypeRDSReport = "rds-report"
const IndexPrefixRDSReport = "rds-reports"
const TemplateNameRDSReport = "rds-reports"

// put the ElasticSearch index for *-lineitems indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameRDSReport).BodyString(TemplateRDSReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index rds-reports.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index rds-reports.", res)
		ctxCancel()
	}
}

const TemplateRDSReport = `
{
	"template": "*-rds-reports",
	"version": 1,
	"mappings": {
		"rds-report": {
			"properties": {
				"account": {
 					"type": "keyword"
 				},
        "reportDate": {
					"type": "date"
				},
        "instances": {
          "type": "nested",
          "properties": {
            "dbInstanceIdentifier": {
              "type": "keyword",
              "norms": false
            },
            "dbInstanceClass": {
              "type": "keyword",
              "norms": false
            },
            "allocatedStorage": {
              "type": "integer"
            },
            "engine": {
              "type": "keyword",
              "norms": false
            },
            "availabilityZone": {
              "type": "keyword",
              "norms": false
            },
            "multiAZ": {
              "type": "boolean"
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
