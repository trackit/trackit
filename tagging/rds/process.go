package taggingrds

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
	Region string              `json:"availabilityZone"`
	Tags   []utils.TagDocument `json:"tags"`
}

type source struct {
	Instance instance `json:"instance"`
}

const sourceIndexName = "rds-reports"
const urlFormat = "https://%s.console.aws.amazon.com/rds/home?region=%s#database:id=%s;is-cluster=false"

// Process generates tagging reports from RDS reports
func Process(ctx context.Context, account int, awsAccount string) ([]utils.TaggingReportDocument, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Processing RDS reports.", nil)

	hits, err := fetchReports(ctx, account)
	if err != nil {
		return nil, err
	}

	var documents []utils.TaggingReportDocument
	for _, hit := range hits {
		document, success := processHit(ctx, hit, awsAccount)
		if success {
			documents = append(documents, document)
		}
	}

	logger.Info(fmt.Sprintf("%d RDS reports processed.", len(documents)), nil)
	return documents, nil
}

// processHit converts an elasticSearch hit into a TaggingReportDocument
// Second argument is true if operation is a success
func processHit(ctx context.Context, hit *elastic.SearchHit, awsAccount string) (utils.TaggingReportDocument, bool) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var source source
	err := json.Unmarshal(*hit.Source, &source)
	if err != nil {
		logger.Error("Could not process a RDS report.", nil)
		return utils.TaggingReportDocument{}, false
	}

	regionForURL := utils.GetRegionForURL(source.Instance.Region)

	document := utils.TaggingReportDocument{
		Account:      awsAccount,
		ResourceID:   source.Instance.ID,
		ResourceType: "rds",
		Region:       source.Instance.Region,
		URL:          fmt.Sprintf(urlFormat, regionForURL, regionForURL, source.Instance.ID),
		Tags:         source.Instance.Tags,
	}

	return document, true
}
