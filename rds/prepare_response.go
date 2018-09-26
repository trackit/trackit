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

package rds

import (
	"strings"
	"context"
	"encoding/json"

	"gopkg.in/olivere/elastic.v5"
)

type (

	// Structure that allow to parse ES response for costs
	ResponseCost struct {
		Accounts struct {
			Buckets []struct {
				Key string `json:"key"`
				Instances struct {
					Buckets []struct {
						Key string `json:"key"`
						Cost struct {
							Value float64 `json:"value"`
						} `json:"cost"`
					} `json:"buckets"`
				} `json:"instances"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for RDS usage report
	ResponseRds struct {
		TopReports struct {
			Buckets []struct {
				TopReportsHits struct {
					Hits struct {
						Hits []struct {
							Source Report `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"top_reports_hits"`
			} `json:"buckets"`
		} `json:"top_reports"`
	}

	// Report format for EC2 usage
	Report struct {
		Account    string `json:"account"`
		ReportDate string `json:"reportDate"`
		ReportType string `json:"reportType"`
		Instances  []struct {
			DBInstanceIdentifier string  `json:"dbInstanceIdentifier"`
			DBInstanceClass      string  `json:"dbInstanceClass"`
			AllocatedStorage     int64   `json:"allocatedStorage"`
			Engine               string  `json:"engine"`
			AvailabilityZone     string  `json:"availabilityZone"`
			MultiAZ              bool    `json:"multiAZ"`
			Cost                 float64 `json:"cost"`
		} `json:"instances"`
	}
)

// addCostToReport adds cost for each instance based on billing data
func addCostToReport(report Report, costs ResponseCost) (Report) {
	for _, accounts := range costs.Accounts.Buckets {
		if accounts.Key != report.Account {
			continue
		}
		for _, instance := range accounts.Instances.Buckets {
			for i := range report.Instances {
				split := strings.Split(instance.Key, ":")

				if len(split) == 7 && split[6] == report.Instances[i].DBInstanceIdentifier {
					report.Instances[i].Cost += instance.Cost.Value
				}
			}
		}
	}
	return report
}

// prepareResponse parses the results from elasticsearch and returns the RDS report
func prepareResponse(ctx context.Context, resRds *elastic.SearchResult, resCost *elastic.SearchResult) (interface{}, error) {
	ec2  := ResponseRds{}
	cost := ResponseCost{}
	reports := []Report{}
	err := json.Unmarshal(*resRds.Aggregations["top_reports"], &ec2.TopReports)
	if err != nil {
		return nil, err
	}
	if resCost != nil {
		json.Unmarshal(*resCost.Aggregations["accounts"], &cost.Accounts)
	}
	for _, account := range ec2.TopReports.Buckets {
		if len(account.TopReportsHits.Hits.Hits) > 0 {
			report := account.TopReportsHits.Hits.Hits[0].Source
			report = addCostToReport(report, cost)
			reports = append(reports, report)
		}
	}
	return reports, nil
}
