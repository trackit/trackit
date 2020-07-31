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
	common.VersioningData{
		Name:        lineItems.TemplateName,
		IndexSuffix: lineItems.IndexSuffix,
		Template:    lineItems.Template,
	},
	common.VersioningData{
		Name:        ebsReports.TemplateName,
		IndexSuffix: ebsReports.IndexSuffix,
		Template:    ebsReports.Template,
	},
	common.VersioningData{
		Name:        ec2Reports.TemplateName,
		IndexSuffix: ec2Reports.IndexSuffix,
		Template:    ec2Reports.Template,
	},
	common.VersioningData{
		Name:        ec2CoverageReports.TemplateName,
		IndexSuffix: ec2CoverageReports.IndexSuffix,
		Template:    ec2CoverageReports.Template,
	},
	common.VersioningData{
		Name:        elasticacheReports.TemplateName,
		IndexSuffix: elasticacheReports.IndexSuffix,
		Template:    elasticacheReports.Template,
	},
	common.VersioningData{
		Name:        esReports.TemplateName,
		IndexSuffix: esReports.IndexSuffix,
		Template:    esReports.Template,
	},
	common.VersioningData{
		Name:        instanceCountReports.TemplateName,
		IndexSuffix: instanceCountReports.IndexSuffix,
		Template:    instanceCountReports.Template,
	},
	common.VersioningData{
		Name:        esReports.TemplateName,
		IndexSuffix: esReports.IndexSuffix,
		Template:    esReports.Template,
	},
	common.VersioningData{
		Name:        lambdaReports.TemplateName,
		IndexSuffix: lambdaReports.IndexSuffix,
		Template:    lambdaReports.Template,
	},
	common.VersioningData{
		Name:        rdsReports.TemplateName,
		IndexSuffix: rdsReports.IndexSuffix,
		Template:    rdsReports.Template,
	},
	common.VersioningData{
		Name:        riEc2Reports.TemplateName,
		IndexSuffix: riEc2Reports.IndexSuffix,
		Template:    riEc2Reports.Template,
	},
	common.VersioningData{
		Name:        rdsRiReports.TemplateName,
		IndexSuffix: rdsRiReports.IndexSuffix,
		Template:    rdsRiReports.Template,
	},
	common.VersioningData{
		Name:        taggingReports.TemplateName,
		IndexSuffix: taggingReports.IndexSuffix,
		Template:    taggingReports.Template,
	},
	common.VersioningData{
		Name:        taggingCompliance.TemplateName,
		IndexSuffix: taggingCompliance.IndexSuffix,
		Template:    taggingCompliance.Template,
	},
	common.VersioningData{
		Name:        accountPlugins.TemplateName,
		IndexSuffix: accountPlugins.IndexSuffix,
		Template:    accountPlugins.Template,
	},
	common.VersioningData{
		Name:        anomaliesDetection.TemplateName,
		IndexSuffix: anomaliesDetection.IndexSuffix,
		Template:    anomaliesDetection.Template,
	},
	common.VersioningData{
		Name:        odToRiEc2Reports.TemplateName,
		IndexSuffix: odToRiEc2Reports.IndexSuffix,
		Template:    odToRiEc2Reports.Template,
	},
}

// put the ElasticSearch index templates indices at startup.
func init() {
	for _, data := range versioningData {
		putTemplate(data.Name, data.Template)
	}
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
