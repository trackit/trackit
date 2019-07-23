//   Copyright 2019 MSolution.IO
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

package lambda

import (
	"context"
	"encoding/json"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/aws/usageReports/lambda"
)

type (
	// Structure that allow to parse ES response for Lambda Daily functions
	ResponseLambdaDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time      string `json:"key_as_string"`
						Functions struct {
							Hits struct {
								Hits []struct {
									Function lambda.FunctionReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"functions"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// FunctionReport has all the information of an Lambda function report
	FunctionReport struct {
		utils.ReportBase
		Function Function `json:"function"`
	}

	// Function contains the information of an Lambda function
	Function struct {
		lambda.FunctionBase
		Tags  map[string]string `json:"tags"`
		Stats lambda.Stats      `json:"stats"`
	}
)

func getLambdaFunctionReportResponse(oldFunction lambda.FunctionReport) FunctionReport {
	tags := make(map[string]string, 0)
	for _, tag := range oldFunction.Function.Tags {
		tags[tag.Key] = tag.Value
	}
	newFunction := FunctionReport{
		ReportBase: oldFunction.ReportBase,
		Function: Function{
			FunctionBase: oldFunction.Function.FunctionBase,
			Tags:         tags,
			Stats:        oldFunction.Function.Stats,
		},
	}
	return newFunction
}

// prepareResponseLambdaDaily parses the results from elasticsearch and returns an array of Lambda daily functions report
func prepareResponseLambdaDaily(ctx context.Context, resLambda *elastic.SearchResult) ([]FunctionReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedLambda ResponseLambdaDaily
	functions := make([]FunctionReport, 0)
	err := json.Unmarshal(*resLambda.Aggregations["accounts"], &parsedLambda.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES Lambda response", err)
		return nil, err
	}
	for _, account := range parsedLambda.Accounts.Buckets {
		var lastDate = ""
		for _, date := range account.Dates.Buckets {
			if date.Time > lastDate {
				lastDate = date.Time
			}
		}
		for _, date := range account.Dates.Buckets {
			if date.Time == lastDate {
				for _, function := range date.Functions.Hits.Hits {
					functions = append(functions, getLambdaFunctionReportResponse(function.Function))
				}
			}
		}
	}
	return functions, nil
}
