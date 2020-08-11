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
	"strconv"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/accountPlugins"
	"github.com/trackit/trackit/es/indexes/anomaliesDetection"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/ebsReports"
	"github.com/trackit/trackit/es/indexes/ec2CoverageReports"
	"github.com/trackit/trackit/es/indexes/ec2Reports"
	"github.com/trackit/trackit/es/indexes/elasticacheReports"
	"github.com/trackit/trackit/es/indexes/esReports"
	"github.com/trackit/trackit/es/indexes/instanceCountReports"
	"github.com/trackit/trackit/es/indexes/lambdaReports"
	"github.com/trackit/trackit/es/indexes/lineItems"
	"github.com/trackit/trackit/es/indexes/odToRiEc2Reports"
	"github.com/trackit/trackit/es/indexes/rdsReports"
	"github.com/trackit/trackit/es/indexes/rdsRiReports"
	"github.com/trackit/trackit/es/indexes/riEc2Reports"
	"github.com/trackit/trackit/es/indexes/taggingCompliance"
	"github.com/trackit/trackit/es/indexes/taggingReports"
)

var versioningData = [...]common.VersioningData{
	lineItems.Model,
	ebsReports.Model,
	ec2Reports.Model,
	ec2CoverageReports.Model,
	elasticacheReports.Model,
	esReports.Model,
	instanceCountReports.Model,
	esReports.Model,
	lambdaReports.Model,
	rdsReports.Model,
	riEc2Reports.Model,
	rdsRiReports.Model,
	taggingReports.Model,
	taggingCompliance.Model,
	accountPlugins.Model,
	anomaliesDetection.Model,
	odToRiEc2Reports.Model,
}

// put the ElasticSearch index templates indices at startup.
func init() {
	for _, data := range versioningData {
		buildTemplatesAndMappings(&data)
		putTemplate(data.Name, data.Template)
	}
}

func buildTemplatesAndMappings(data *common.VersioningData) {
	data.Mapping = `
	{
		"` + data.Type + `": {
			"properties": ` + data.MappingProperties + `
			,
			"_all": {
				"enabled": false
			},
			"numeric_detection": false,
			"date_detection": false
		}
	}
	`

	data.Template = `
	{
		"template": "*-` + data.Name + `",
		"version": ` + strconv.Itoa(data.Version) + `,
		"mappings": ` + data.Mapping + `
	}
	`
}

func putTemplate(templateName string, template string) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(templateName).BodyString(template).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index template.", map[string]interface{}{
			"templateName": templateName,
			"error":        err,
		})
	} else {
		jsonlog.DefaultLogger.Info("Put ES index template.", map[string]interface{}{
			"templateName": templateName,
			"res":          res,
		})
		ctxCancel()
	}
}
