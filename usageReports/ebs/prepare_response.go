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

package ebs

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/aws/usageReports/ebs"
	"github.com/trackit/trackit/errors"
)

type (

	// ResponseCost allows us to parse an ES response for costs
	ResponseCost struct {
		Accounts struct {
			Buckets []struct {
				Key       string `json:"key"`
				Snapshots struct {
					Buckets []struct {
						Key  string `json:"key"`
						Cost struct {
							Value float64 `json:"value"`
						} `json:"cost"`
					} `json:"buckets"`
				} `json:"snapshots"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// ResponseEbsMonthly allows us to parse an ES response for EBS Monthly snapshots
	ResponseEbsMonthly struct {
		Accounts struct {
			Buckets []struct {
				Snapshots struct {
					Hits struct {
						Hits []struct {
							Snapshot ebs.SnapshotReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"snapshots"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// ResponseEbsDaily allows us to parse an ES response for EBS Daily snapshots
	ResponseEbsDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time      string `json:"key_as_string"`
						Snapshots struct {
							Hits struct {
								Hits []struct {
									Snapshot ebs.SnapshotReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"snapshots"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// SnapshotReport has all the information of an EBS snapshot report
	SnapshotReport struct {
		utils.ReportBase
		Snapshot Snapshot `json:"snapshot"`
	}

	// Snapshot contains the information of an EBS snapshot
	Snapshot struct {
		ebs.SnapshotBase
		Tags   map[string]string `json:"tags"`
		Volume Volume            `json:"volume"`
		Cost   float64           `json:"cost"`
	}

	// Volume contains information about EBS volumes
	Volume struct {
		Id   string `json:"id"`
		Size int64  `json:"size"`
	}
)

func getEbsSnapshotReportResponse(oldSnapshot ebs.SnapshotReport) SnapshotReport {
	tags := make(map[string]string, len(oldSnapshot.Snapshot.Tags))
	for _, tag := range oldSnapshot.Snapshot.Tags {
		tags[tag.Key] = tag.Value
	}
	newSnapshot := SnapshotReport{
		ReportBase: oldSnapshot.ReportBase,
		Snapshot: Snapshot{
			SnapshotBase: oldSnapshot.Snapshot.SnapshotBase,
			Tags:         tags,
			Cost:         oldSnapshot.Snapshot.Cost,
			Volume: Volume{
				Id:   oldSnapshot.Snapshot.Volume.Id,
				Size: oldSnapshot.Snapshot.Volume.Size,
			},
		},
	}
	return newSnapshot
}

// prepareResponseEbsDaily parses the results from elasticsearch and returns an array of EBS daily snapshots report
func prepareResponseEbsDaily(ctx context.Context, resEbs *elastic.SearchResult, resCost *elastic.SearchResult) ([]SnapshotReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedEbs ResponseEbsDaily
	var parsedCost ResponseCost
	snapshots := make([]SnapshotReport, 0)
	err := json.Unmarshal(*resEbs.Aggregations["accounts"], &parsedEbs.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES EBS response", err)
		return nil, err
	}
	if resCost != nil {
		err = json.Unmarshal(*resCost.Aggregations["accounts"], &parsedCost.Accounts)
		if err != nil {
			logger.Error("Error while unmarshaling ES cost response", err)
		}
	}
	for _, account := range parsedEbs.Accounts.Buckets {
		var lastDate = ""
		for _, date := range account.Dates.Buckets {
			if date.Time > lastDate {
				lastDate = date.Time
			}
		}
		for _, date := range account.Dates.Buckets {
			if date.Time == lastDate {
				for _, snapshot := range date.Snapshots.Hits.Hits {
					snapshots = append(snapshots, getEbsSnapshotReportResponse(snapshot.Snapshot))
				}
			}
		}
	}
	return snapshots, nil
}

// prepareResponseEbsMonthly parses the results from elasticsearch and returns an array of EBS monthly snapshots report
func prepareResponseEbsMonthly(ctx context.Context, resEbs *elastic.SearchResult) ([]SnapshotReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseEbsMonthly
	snapshots := make([]SnapshotReport, 0)
	err := json.Unmarshal(*resEbs.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES EBS response", err)
		return nil, errors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, snapshot := range account.Snapshots.Hits.Hits {
			snapshots = append(snapshots, getEbsSnapshotReportResponse(snapshot.Snapshot))
		}
	}
	return snapshots, nil
}
