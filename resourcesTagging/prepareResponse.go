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

package resourcesTagging

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/tagging/utils"

	terrors "github.com/trackit/trackit/errors"
)

type (

	// Structure that allow to parse ES response for resources tagging
	ResponseResources struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time      string `json:"key_as_string"`
						Resources struct {
							Hits struct {
								Hits []struct {
									Resource utils.TaggingReportDocument `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"resources"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}
)

// prepareResponseResources parses the results from elasticsearch and returns an array of resources report
func prepareResponseResources(ctx context.Context, resResources *elastic.SearchResult) ([]utils.TaggingReportDocument, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseResources
	resources := make([]utils.TaggingReportDocument, 0)
	err := json.Unmarshal(*resResources.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES resources response", err)
		return nil, terrors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, date := range account.Dates.Buckets {
			for _, resource := range date.Resources.Hits.Hits {
				resources = append(resources, resource.Resource)
			}
		}

	}
	return resources, nil
}