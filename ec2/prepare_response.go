//   Copyright 2017 MSolution.IO
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

package ec2

import (
	"context"
	"encoding/json"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws/ec2"
)

// parseESResult parses an *elastic.SearchResult according to it's resultType
func parseESResult(ctx context.Context, res *elastic.SearchResult) (*ec2.ReportInfo, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	reports := make([]ec2.ReportInfo, 0)
	for _, hit := range res.Hits.Hits {
		var parsedDocument ec2.ReportInfo
		err := json.Unmarshal(*hit.Source, &parsedDocument)
		if err != nil {
			logger.Error("Failed to parse elasticsearch document.", err.Error())
			return nil, err
		}
		reports = append(reports, parsedDocument)
	}
	return &reports[0], nil
}

// prepareResponse parses the results from elasticsearch and returns a list of EC2 reports with stats of each instances
func prepareResponse(ctx context.Context, rawReport *elastic.SearchResult) (interface{}, error) {
	report, err := parseESResult(ctx, rawReport)
	if err != nil {
		return nil, err
	}
	return report, nil
}
