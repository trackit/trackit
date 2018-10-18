//   Copyright 2018 MSolution.IO
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
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws/usageReports/rds"
	"sort"
)

type (

	// Structure that allow to parse ES response for costs
	ResponseCost struct {
		Accounts struct {
			Buckets []struct {
				Key       string `json:"key"`
				Instances struct {
					Buckets []struct {
						Key  string `json:"key"`
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
							Source rds.Report `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"top_reports_hits"`
			} `json:"buckets"`
		} `json:"top_reports"`
	}
)

// addCostToReport adds cost for each instance based on billing data
func addCostToReport(report rds.Report, costs ResponseCost) rds.Report {
	for _, accounts := range costs.Accounts.Buckets {
		if accounts.Key != report.Account {
			continue
		}
		for _, instance := range accounts.Instances.Buckets {
			for i := range report.Instances {
				if report.Instances[i].CostDetail == nil {
					report.Instances[i].CostDetail = make(map[string]float64, 0)
				}
				// format in billing data for an RDS instance is: "arn:aws:rds:us-west-2:394125495069:db:instancename"
				// so i get the 7th element of the split by ":"
				split := strings.Split(instance.Key, ":")
				if len(split) == 7 && split[6] == report.Instances[i].DBInstanceIdentifier {
					report.Instances[i].Cost += instance.Cost.Value
					report.Instances[i].CostDetail["instance"] += instance.Cost.Value
				}
			}
		}
	}
	return report
}

// prepareResponseRdsDaily parses the results from elasticsearch and returns the RDS report
func prepareResponseRdsDaily(ctx context.Context, resRds *elastic.SearchResult, resCost *elastic.SearchResult) ([]rds.Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	parsedRds := ResponseRds{}
	cost := ResponseCost{}
	reports := []rds.Report{}
	err := json.Unmarshal(*resRds.Aggregations["top_reports"], &parsedRds.TopReports)
	if err != nil {
		logger.Error("Error while unmarshaling ES RDS response", err)
		return nil, err
	}
	if resCost != nil {
		err = json.Unmarshal(*resCost.Aggregations["accounts"], &cost.Accounts)
		if err != nil {
			logger.Error("Error while unmarshaling ES cost response", err)
		}
	}
	for _, account := range parsedRds.TopReports.Buckets {
		if len(account.TopReportsHits.Hits.Hits) > 0 {
			report := account.TopReportsHits.Hits.Hits[0].Source
			report = addCostToReport(report, cost)
			reports = append(reports, report)
		}
	}
	return reports, nil
}

// prepareResponseRdsMonthly parses the results from elasticsearch and returns the RDS usage report
func prepareResponseRdsMonthly(ctx context.Context, resRds *elastic.SearchResult) ([]rds.Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	response := ResponseRds{}
	reports := []rds.Report{}
	err := json.Unmarshal(*resRds.Aggregations["top_reports"], &response.TopReports)
	if err != nil {
		logger.Error("Error while unmarshaling ES RDS response", err)
		return nil, err
	}
	for _, account := range response.TopReports.Buckets {
		if len(account.TopReportsHits.Hits.Hits) > 0 {
			reports = append(reports, account.TopReportsHits.Hits.Hits[0].Source)
		}
	}
	return reports, nil
}

func isInstanceUnused(instance rds.Instance) bool {
	average := instance.CpuAverage
	peak := instance.CpuPeak
	if peak >= 60.0 {
		return false
	} else if average >= 10.0 {
		return false
	}
	return true
}

// prepareResponseRdsUnused filter reports to get the top instances
func prepareResponseRdsUnused(params rdsUnusedQueryParams, reports []rds.Report) (int, []rds.Instance, error) {
	instances := []rds.Instance{}
	for _, report := range reports {
		for _, instance := range report.Instances {
			if isInstanceUnused(instance) {
				instances = append(instances, instance)
			}
		}
	}
	sort.SliceStable(instances, func(i, j int) bool { return instances[i].Cost > instances[j].Cost })
	if params.count >= 0 && params.count <= len(instances) {
		return http.StatusOK, instances[0:params.count], nil
	}
	return http.StatusOK, instances, nil
}
