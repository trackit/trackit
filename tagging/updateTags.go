package tagging

import (
	"context"
	"fmt"
	"time"

	"github.com/trackit/jsonlog"

	bulk "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
	ebs "github.com/trackit/trackit/tagging/ebs"
	ec2 "github.com/trackit/trackit/tagging/ec2"
	ec2Ri "github.com/trackit/trackit/tagging/ec2Ri"
	lambda "github.com/trackit/trackit/tagging/lambda"
	rds "github.com/trackit/trackit/tagging/rds"
	rdsRi "github.com/trackit/trackit/tagging/rdsRi"
	"github.com/trackit/trackit/tagging/utils"
)

type process func(ctx context.Context, account int, awsAccount string, resourceTypeString string) ([]utils.TaggingReportDocument, error)

type processor struct {
	Name string
	Run  process
}

const destIndexName = "tagging"
const destTypeName = "tagging"

var processors = []processor{
	processor{
		Name: "ec2",
		Run:  ec2.Process,
	},
	processor{
		Name: "ebs",
		Run:  ebs.Process,
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
	processor{
		Name: "ec2-ri",
		Run:  ec2Ri.Process,
	},
}

// UpdateTagsForAccount updates tags in ES for the specified AWS account
func UpdateTagsForAccount(ctx context.Context, account int, awsAccount string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var documents []utils.TaggingReportDocument

	for _, processor := range processors {
		newDocuments, err := processor.Run(ctx, account, awsAccount, processor.Name)
		if err == nil {
			documents = append(documents, newDocuments...)
		} else {
			logger.Error(fmt.Sprintf("Generation of tagging reports for resources of type '%s' failed: %s", processor.Name, err.Error()), nil)
		}
	}

	logger.Info("Pushing generated tagging reports to ES.", nil)
	reportDate := time.Now().UTC()
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
