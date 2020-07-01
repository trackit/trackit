package taggingec2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/tagging/utils"
)

type instance struct {
	ID     string              `json:"id"`
	Region string              `json:"region"`
	Tags   []utils.TagDocument `json:"tags"`
}

type source struct {
	Instance instance `json:"instance"`
}

const sourceIndexName = "ec2-reports"
const urlFormat = "https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#Instances:instanceId=%s"

// ProcessEc2 generates tagging reports from EC2 reports
func ProcessEc2(ctx context.Context, account int, awsAccount string) ([]utils.TaggingReportDocument, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Processing EC2 reports.", nil)

	hits, err := fetchEc2Reports(ctx, account)
	if err != nil {
		return nil, err
	}

	var documents []utils.TaggingReportDocument
	for _, hit := range hits {
		document, success := processEc2Hit(ctx, hit, awsAccount)
		if success {
			documents = append(documents, document)
		}
	}

	logger.Info(fmt.Sprintf("%d EC2 reports processed.", len(documents)), nil)
	return documents, nil
}

// processEc2Hit converts an elasticSearch hit into a TaggingReportDocument
// Second argument is true if operation is a success
func processEc2Hit(ctx context.Context, hit *elastic.SearchHit, awsAccount string) (utils.TaggingReportDocument, bool) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var source source
	err := json.Unmarshal(*hit.Source, &source)
	if err != nil {
		logger.Error("Could not process a EC2 report.", nil)
		return utils.TaggingReportDocument{}, false
	}

	regionForURL := utils.GetRegionForURL(source.Instance.Region)

	document := utils.TaggingReportDocument{
		Account:      awsAccount,
		ResourceID:   source.Instance.ID,
		ResourceType: "ec2",
		Region:       source.Instance.Region,
		URL:          fmt.Sprintf(urlFormat, regionForURL, regionForURL, source.Instance.ID),
		Tags:         source.Instance.Tags,
	}
	return document, true
}
