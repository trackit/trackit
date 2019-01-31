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

package onDemandToRiEc2

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
)

const TypeOdToRiEC2Report = "od-to-ri-ec2-report"
const IndexPrefixOdToRiEC2Report = "od-to-ri-ec2-reports"
const TemplateNameOdToRiEC2Report = "od-to-ri-ec2-reports"

// put the ElasticSearch index for *-od-to-ri-ec2-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameOdToRiEC2Report).BodyString(TemplateOdToRiEc2Report).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index OdToRiEC2Report.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index OdToRiEC2Report.", res)
		ctxCancel()
	}
}

const TemplateOdToRiEc2Report = `
{
	"template": "*-od-to-ri-ec2-reports",
	"version": 1,
	"mappings": {
		"od-to-ri-ec2-report": {
			"properties": {
				"account": {
					"type": "keyword"
				},
				"reportDate": {
					"type": "date"
				},
        "onDemandMonthlyCostTotal": {
          "type": "double"
        },
        "onDemand1yrCostTotal": {
          "type": "double"
        },
        "onDemand3yrCostTotal": {
          "type": "double"
        },
        "ri1yrMonthlyCostTotal": {
          "type": "double"
        },
        "ri1yrCostTotal": {
          "type": "double"
        },
        "ri1yrSavingTotal": {
          "type": "double"
        },
        "ri3yrMonthlyCostTotal": {
          "type": "double"
        },
        "ri3yrCostTotal": {
          "type": "double"
        },
        "ri3yrSavingTotal": {
          "type": "double"
        },
        "instances": {
          "type": "nested",
          "properties": {
            "region": {
              "type": "keyword"
            },
            "instanceType": {
              "type": "keyword"
            },
            "platform": {
              "type": "keyword"
            },
            "instanceCount": {
              "type": "integer"
            },
            "onDemandMonthlyCostPerUnit": {
              "type": "double"
            },
            "onDemandMonthlyCostTotal": {
              "type": "double"
            },
            "onDemand1yrCostPerUnit": {
              "type": "double"
            },
            "onDemand1yrCostTotal": {
              "type": "double"
            },
            "onDemand3yrCostPerUnit": {
              "type": "double"
            },
            "onDemand3yrCostTotal": {
              "type": "double"
            },
            "reservationType": {
              "type": "keyword"
            },
            "ri1yrMonthlyCostPerUnit": {
              "type": "double"
            },
            "ri1yrMonthlyCostTotal": {
              "type": "double"
            },
            "ri1yrCostPerUnit": {
              "type": "double"
            },
            "ri1yrCostTotal": {
              "type": "double"
            },
            "ri1yrSavingPerUnit": {
              "type": "double"
            },
            "ri1yrSavingTotal": {
              "type": "double"
            },
            "ri3yrMonthlyCostPerUnit": {
              "type": "double"
            },
            "ri3yrMonthlyCostTotal": {
              "type": "double"
            },
            "ri3yrCostPerUnit": {
              "type": "double"
            },
            "ri3yrCostTotal": {
              "type": "double"
            },
            "ri3yrSavingPerUnit": {
              "type": "double"
            },
            "ri3yrSavingTotal": {
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
