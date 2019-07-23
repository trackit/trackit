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

package plugins_account_core

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
)

const TypeAccountPlugins = "account-plugin"
const IndexPrefixAccountPlugin = "account-plugins"
const TemplateNameAccountPlugin = "account-plugins"

// put the ElasticSearch index for *-account-plugins indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(TemplateNameAccountPlugin).BodyString(TemplateAccountPlugin).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index account-plugins.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index account-plugins.", res)
		ctxCancel()
	}
}

const TemplateAccountPlugin = `
{
  "template": "*-account-plugins",
  "version": 3,
  "mappings": {
    "account-plugin": {
      "properties": {
        "accountPluginIdx": {
          "type": "keyword"
        },
        "account": {
          "type": "keyword"
        },
        "reportDate": {
          "type": "date"
        },
        "pluginName": {
          "type": "keyword"
        },
        "category": {
          "type": "keyword"
        },
        "label": {
          "type": "keyword"
        },
        "result": {
          "type": "keyword"
        },
        "status": {
          "type": "keyword"
        },
        "details": {
          "type": "keyword"
        },
        "error": {
          "type": "keyword"
        },
        "checked": {
          "type": "integer"
        },
        "passed": {
          "type": "integer"
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
