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

const typeTaggingCompliance = "tagging-compliance"
const indexPrefixTaggingCompliance = "tagging-compliance"
const templateNameTaggingCompliance = "tagging-compliance"

// put the ElasticSearch index for *-tagging-compliance indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(templateNameTaggingCompliance).BodyString(templateTaggingCompliance).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index tagging-compliance.", err)
		ctxCancel()
	} else {
		jsonlog.DefaultLogger.Info("Put ES index tagging-compliance.", res)
		ctxCancel()
	}
}

const templateTaggingCompliance = `
{
    "template":"*-tagging-compliance",
    "version":1,
    "mappings":{
        "tagging-compliance":{
            "properties":{
                "reportDate":{
                    "type":"date"
                },
                "total":{
                    "type":"long"
                },
                "totallyTagged":{
                    "type":"long"
                },
                "partiallyTagged":{
                    "type":"long"
                },
                "notTagged":{
                    "type":"long"
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
