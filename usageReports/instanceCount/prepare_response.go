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

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/instanceCountReports"
)

type (

	// ResponseCost allows us to parse an ES response for costs
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

	// ResponseInstanceCountMonthly allows us to parse an ES response for InstanceCount Monthly
	ResponseInstanceCountMonthly struct {
		Accounts struct {
			Buckets []struct {
				InstanceCount struct {
					Hits struct {
						Hits []struct {
							InstanceCount instanceCountReports.InstanceCountReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"reports"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// ResponseInstanceCountDaily allows us to parse an ES response for InstanceCount Daily
	ResponseInstanceCountDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time string `json:"key_as_string"`
						In   struct {
							Hits struct {
								Hits []struct {
									InstanceCount instanceCountReports.InstanceCountReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"instanceCount"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// InstanceCountReport is saved in ES to have all the information of an InstanceCount
	InstanceCountReport struct {
		common.ReportBase
		InstanceCount InstanceCount `json:"instanceCount"`
	}

	// InstanceCount contains all the information of an InstanceCount
	InstanceCount struct {
		Type   string               `json:"instanceType"`
		Region string               `json:"region"`
		Hours  []InstanceCountHours `json:"hours"`
	}

	InstanceCountHours struct {
		Hour  time.Time `json:"hour"`
		Count float64   `json:"count"`
	}
)

func getInstanceCountSnapshotReportResponse(oldInstanceCount instanceCountReports.InstanceCountReport) InstanceCountReport {
	hours := make([]InstanceCountHours, 0)
	for _, value := range oldInstanceCount.InstanceCount.Hours {
		hours = append(hours, InstanceCountHours{
			Hour:  value.Hour,
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

// There is a (likely out of date) prepareResponseInstanceCountDaily method left in source control but it's unused, so I've removed it

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
