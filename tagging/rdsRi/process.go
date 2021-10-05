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

package taggingrdsri

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es/indexes/rdsRiReports"
	"github.com/trackit/trackit/es/indexes/taggingReports"
	"github.com/trackit/trackit/tagging/utils"
)

const urlFormat = "https://console.aws.amazon.com/rds/home?region=%s#reserved-db-instance:ids=%s"

// Process generates tagging reports from RDS reserved instances reports
func Process(ctx context.Context, userId int, resourceTypeString string) ([]taggingReports.TaggingReportDocument, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Processing reports.", map[string]interface{}{
		"type": resourceTypeString,
	})

	hits, err := fetchReports(ctx, userId)
	if err != nil {
		return nil, err
	}

	var documents []taggingReports.TaggingReportDocument
	for _, hit := range hits {
		document, success := processHit(ctx, hit, resourceTypeString)
		if success {
			documents = append(documents, document)
		}
	}

	logger.Info("Reports processed.", map[string]interface{}{
		"type":  resourceTypeString,
		"count": len(documents),
	})
	return documents, nil
}

// processHit converts an elasticSearch hit into a TaggingReportDocument
// Second argument is true if operation is a success
func processHit(ctx context.Context, hit *elastic.SearchHit, resourceTypeString string) (taggingReports.TaggingReportDocument, bool) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var source rdsRiReports.InstanceReport
	err := json.Unmarshal(*hit.Source, &source)
	if err != nil {
		logger.Error("Could not process report.", map[string]interface{}{
			"type": resourceTypeString,
		})
		return taggingReports.TaggingReportDocument{}, false
	}

	regionForURL := utils.GetRegionForURL(source.Instance.AvailabilityZone)

	document := taggingReports.TaggingReportDocument{
		Account:      source.Account,
		ResourceID:   source.Instance.DBInstanceIdentifier,
		ResourceType: resourceTypeString,
		Region:       source.Instance.AvailabilityZone,
		URL:          fmt.Sprintf(urlFormat, regionForURL, source.Instance.DBInstanceIdentifier),
		Tags:         source.Instance.Tags,
	}

	return document, true
}
