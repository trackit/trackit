//   Copyright 2020 MSolution.IO
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

package tagging

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const typeTaggingReport = "tagging-reports"
const indexPrefixTaggingReport = "tagging-reports"
const templateNameTaggingReport = "tagging-reports"

// put the ElasticSearch index for *-tagging-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(templateNameTaggingReport).BodyString(templateTaggingReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES indext tagging-reports.", err)
		ctxCancel()
	} else {
		jsonlog.DefaultLogger.Info("Put ES index tagging-reports.", res)
		ctxCancel()
	}
}

const templateTaggingReport = `
{
    "template":"*-tagging-reports",
    "version":1,
    "mappings":{
        "tagging-reports":{
            "properties":{
                "account":{
                    "type":"keyword"
                },
                "region":{
                    "type":"keyword"
                },
                "reportDate":{
                    "type":"date"
                },
                "resourceId":{
                    "type":"keyword"
                },
                "resourceType":{
                    "type":"keyword"
                },
                "tags":{
                    "type":"nested",
                    "properties":{
                        "key":{
                            "type":"keyword"
                        },
                        "value":{
                            "type":"keyword"
                        }
                    }
                },
                "url":{
                    "type":"keyword"
                }
			},
			"_all": {
				"enabled": false
			},
			"date_detection": false,
			"numeric_detection": false
        }
    }
}
`
