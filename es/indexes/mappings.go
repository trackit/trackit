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

package indexes

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/lineItems"
)

var mappings = map[string]string{
	lineItems.TemplateName: lineItems.Template,
}

// put the ElasticSearch index templates indices at startup.
func init() {
	for templateName, template := range mappings {
		putTemplate(templateName, template)
	}
}

func putTemplate(templateName string, template string) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(templateName).BodyString(template).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index lineitems.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index lineitems.", res)
		ctxCancel()
	}
}
