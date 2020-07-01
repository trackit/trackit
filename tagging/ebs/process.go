package taggingebs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/tagging/utils"
)

type snapshot struct {
	ID     string              `json:"id"`
	Region string              `json:"region"`
	Tags   []utils.TagDocument `json:"tags"`
}

type source struct {
	Snapshot snapshot `json:"snapshot"`
}

const sourceIndexName = "ebs-reports"
const urlFormat = "https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#Snapshots:all;search=%s"

// Process generates tagging reports from EBS reports
func Process(ctx context.Context, account int, awsAccount string, resourceTypeString string) ([]utils.TaggingReportDocument, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info(fmt.Sprintf("Processing %s reports.", resourceTypeString), nil)

	hits, err := fetchReports(ctx, account)
	if err != nil {
		return nil, err
	}

	var documents []utils.TaggingReportDocument
	for _, hit := range hits {
		document, success := processHit(ctx, hit, awsAccount, resourceTypeString)
		if success {
			documents = append(documents, document)
		}
	}

	logger.Info(fmt.Sprintf("%d %s reports processed.", len(documents), resourceTypeString), nil)
	return documents, nil
}

// processHit converts an elasticSearch hit into a TaggingReportDocument
// Second argument is true if operation is a success
func processHit(ctx context.Context, hit *elastic.SearchHit, awsAccount string, resourceTypeString string) (utils.TaggingReportDocument, bool) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var source source
	err := json.Unmarshal(*hit.Source, &source)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not process a %s report.", resourceTypeString), nil)
		return utils.TaggingReportDocument{}, false
	}

	regionForURL := utils.GetRegionForURL(source.Snapshot.Region)

	document := utils.TaggingReportDocument{
		Account:      awsAccount,
		ResourceID:   source.Snapshot.ID,
		ResourceType: resourceTypeString,
		Region:       source.Snapshot.Region,
		URL:          fmt.Sprintf(urlFormat, regionForURL, regionForURL, source.Snapshot.ID),
		Tags:         source.Snapshot.Tags,
	}
	return document, true
}
