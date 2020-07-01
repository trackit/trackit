package tagging

import (
	"context"
	"fmt"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
	ec2 "github.com/trackit/trackit/tagging/ec2"
	"github.com/trackit/trackit/tagging/utils"
)

type process func(ctx context.Context, account int, awsAccount string) ([]utils.TaggingReportDocument, error)

type processor struct {
	Name string
	Run  process
}

const destIndexName = "tagging"

// UpdateTagsForAccount updates tags in ES for the specified AWS account
func UpdateTagsForAccount(ctx context.Context, account int, awsAccount string) error {
	client := es.Client
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var documents []utils.TaggingReportDocument

	processors := []processor{
		processor{
			Name: "ec2",
			Run:  ec2.ProcessEc2,
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

	reportDate := time.Now().UTC()

	destIndexName := es.IndexNameForUserId(account, destIndexName)
	for _, document := range documents {
		document.ReportDate = reportDate

		_, err := client.Index().Index(destIndexName).Type("tagging").BodyJson(document).Do(ctx)
		if err != nil {
			logger.Error("Could not insert a document.", err)
		}
	}

	return nil
}
