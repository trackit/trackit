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
	"fmt"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	bulk "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
	ebs "github.com/trackit/trackit/tagging/ebs"
	ec2 "github.com/trackit/trackit/tagging/ec2"
	ec2Ri "github.com/trackit/trackit/tagging/ec2Ri"
	elasticache "github.com/trackit/trackit/tagging/elasticache"
	esProc "github.com/trackit/trackit/tagging/es"
	lambda "github.com/trackit/trackit/tagging/lambda"
	rds "github.com/trackit/trackit/tagging/rds"
	rdsRi "github.com/trackit/trackit/tagging/rdsRi"
	"github.com/trackit/trackit/tagging/utils"
)

type process func(ctx context.Context, awsAccount aws.AwsAccount, resourceTypeString string) ([]utils.TaggingReportDocument, error)

type processor struct {
	Name string
	Run  process
}

const destIndexName = "tagging-reports"
const destTypeName = "tagging-reports"

var processors = []processor{
	processor{
		Name: "ebs",
		Run:  ebs.Process,
	},
	processor{
		Name: "ec2",
		Run:  ec2.Process,
	},
	processor{
		Name: "ec2-ri",
		Run:  ec2Ri.Process,
	},
	processor{
		Name: "elasticache",
		Run:  elasticache.Process,
	},
	processor{
		Name: "es",
		Run:  esProc.Process,
	},
	processor{
		Name: "lambda",
		Run:  lambda.Process,
	},
	processor{
		Name: "rds",
		Run:  rds.Process,
	},
	processor{
		Name: "rds-ri",
		Run:  rdsRi.Process,
	},
}

// UpdateTagsForAccount updates tags in ES for the specified AWS account
func UpdateTagsForAccount(ctx context.Context, awsAccount aws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var documents []utils.TaggingReportDocument

	for _, processor := range processors {
		newDocuments, err := processor.Run(ctx, awsAccount, processor.Name)
		if err == nil {
			documents = append(documents, newDocuments...)
		} else {
			logger.Error(fmt.Sprintf("Generation of tagging reports for resources of type '%s' failed: %s", processor.Name, err.Error()), nil)
		}
	}

	return pushToEs(ctx, documents, awsAccount.UserId)
}

func pushToEs(ctx context.Context, documents []utils.TaggingReportDocument, account int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	reportDate := time.Now().UTC()
	logger.Info("Pushing generated tagging reports to ES.", map[string]interface{}{
		"reportDate": reportDate.String(),
	})
	destIndexName := es.IndexNameForUserId(account, destIndexName)
	bulkProcessor, err := bulk.GetBulkProcessor(ctx)
	if err != nil {
		return err
	}

	for _, document := range documents {
		document.ReportDate = reportDate

		documentID, err := utils.GenerateBulkID(document)
		if err != nil {
			logger.Error("Could not add a tagging report to bulk processor.", err.Error())
			continue
		}

		bulkProcessor = bulk.AddDocToBulkProcessor(bulkProcessor, document, destTypeName, destIndexName, documentID)
	}

	bulkProcessor.Flush()
	err = bulkProcessor.Close()
	if err != nil {
		logger.Error("Failed to put tagging reports in ES", err.Error())
		return err
	}

	return nil
}
