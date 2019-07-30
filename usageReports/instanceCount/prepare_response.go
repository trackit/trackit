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

package instanceCount

import (
	"context"
	"encoding/json"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/aws/usageReports/instanceCount"
	"github.com/trackit/trackit/errors"
)

type (

	// Structure that allow to parse ES response for costs
	ResponseCost struct {
		Accounts struct {
			Buckets []struct {
				Key           string `json:"key"`
				InstanceCount struct {
					Buckets []struct {
						Key  string `json:"key"`
						Cost struct {
							Value float64 `json:"value"`
						} `json:"cost"`
					} `json:"buckets"`
				} `json:"instanceCount"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for InstanceCount Monthly
	ResponseInstanceCountMonthly struct {
		Accounts struct {
			Buckets []struct {
				InstanceCount struct {
					Hits struct {
						Hits []struct {
							InstanceCount instanceCount.InstanceCountReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"reports"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for InstanceCount Daily
	ResponseInstanceCountDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time string `json:"key_as_string"`
						In   struct {
							Hits struct {
								Hits []struct {
									InstanceCount instanceCount.InstanceCountReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"instanceCount"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// InstanceCount is saved in ES to have all the information of an InstanceCount
	InstanceCountReport struct {
		utils.ReportBase
		InstanceCount InstanceCount `json:"instanceCount"`
	}

	// InstanceCount contains all the information of an InstanceCount
	InstanceCount struct {
		Type   string              `json:"instanceType"`
		Region string              `json:"region"`
		Hours []InstanceCountHours `json:"hours"`
	}

	InstanceCountHours struct {
		Hour  time.Time `json:"hour"`
		Count float64   `json:"count"`
	}
)

func getInstanceCountSnapshotReportResponse(oldInstanceCount instanceCount.InstanceCountReport) InstanceCountReport {
	hours := make([]InstanceCountHours, 0)
	for _, value := range oldInstanceCount.InstanceCount.Hours {
		hours = append(hours, InstanceCountHours{
			Hour: value.Hour,
			Count: value.Count,
		})
	}
	newInstance := InstanceCountReport{
		ReportBase: oldInstanceCount.ReportBase,
		InstanceCount: InstanceCount{
			Type:   oldInstanceCount.InstanceCount.Type,
			Region: oldInstanceCount.InstanceCount.Region,
			Hours:  hours,
		},
	}
	return newInstance
}

// prepareResponseInstanceCountDaily parses the results from elasticsearch and returns an array of InstanceCount daily instanceCount report
func prepareResponseInstanceCountDaily(ctx context.Context, resInstanceCount *elastic.SearchResult, resCost *elastic.SearchResult) ([]InstanceCountReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedInstanceCount ResponseInstanceCountDaily
	reports := make([]InstanceCountReport, 0)
	err := json.Unmarshal(*resInstanceCount.Aggregations["accounts"], &parsedInstanceCount)
	if err != nil {
		logger.Error("Error while unmarshaling ES InstanceCount response", err)
		return nil, err
	}
	for _, account := range parsedInstanceCount.Accounts.Buckets {
		var lastDate = ""
		for _, date := range account.Dates.Buckets {
			if date.Time > lastDate {
				lastDate = date.Time
			}
		}
		for _, date := range account.Dates.Buckets {
			if date.Time == lastDate {
				for _, report := range date.In.Hits.Hits {
					reports = append(reports, getInstanceCountSnapshotReportResponse(report.InstanceCount))
				}
			}
		}
	}
	return reports, nil
}

// prepareResponseInstanceCountMonthly parses the results from elasticsearch and returns an array of InstanceCount monthly report
func prepareResponseInstanceCountMonthly(ctx context.Context, resInstanceCount *elastic.SearchResult) ([]InstanceCountReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseInstanceCountMonthly
	reports := make([]InstanceCountReport, 0)
	err := json.Unmarshal(*resInstanceCount.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES InstanceCount response", err)
		return nil, errors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, report := range account.InstanceCount.Hits.Hits {
			reports = append(reports, getInstanceCountSnapshotReportResponse(report.InstanceCount))
		}
	}
	return reports, nil
}
