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
	"github.com/trackit/trackit/tagging/utils"
)

type process func(ctx context.Context, account int, awsAccount string) ([]utils.TaggingReportDocument, error)

type processor struct {
	Name string
	Run  process
}

const destIndexName = "tagging"
const destTypeName = "tagging"

// UpdateTagsForAccount updates tags in ES for the specified AWS account
func UpdateTagsForAccount(ctx context.Context, account int, awsAccount string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var documents []utils.TaggingReportDocument

	processors := []processor{
		processor{
			Name: "ec2",
			Run:  ec2.ProcessEc2,
		},
		processor{
			Name: "ebs",
			Run:  ebs.ProcessEbs,
		},
	}

	for _, processor := range processors {
		newDocuments, err := processor.Run(ctx, account, awsAccount)
		if err == nil {
			documents = append(documents, newDocuments...)
		} else {
			logger.Error(fmt.Sprintf("Processor '%s' failed: %s", processor.Name, err.Error()), nil)
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
