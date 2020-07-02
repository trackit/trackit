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

package tagginglambda

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/tagging/utils"
)

type function struct {
	Name   string              `json:"name"`
	Region string              `json:"region"`
	Tags   []utils.TagDocument `json:"tags"`
}

type source struct {
	Function function `json:"function"`
}

const sourceIndexName = "lambda-reports"
const urlFormat = "https://%s.console.aws.amazon.com/lambda/home?region=%s#/functions/%s"

// Process generates tagging reports from Lambda reports
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

	regionForURL := utils.GetRegionForURL(source.Function.Region)

	document := utils.TaggingReportDocument{
		Account:      awsAccount,
		ResourceID:   source.Function.Name,
		ResourceType: resourceTypeString,
		Region:       source.Function.Region,
		URL:          fmt.Sprintf(urlFormat, regionForURL, regionForURL, source.Function.Name),
		Tags:         source.Function.Tags,
	}
	return document, true
}
