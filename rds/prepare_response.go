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
	"net/http"
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

	// Report format for RDS usage
	Report struct {
		Account    string     `json:"account"`
		ReportDate string     `json:"reportDate"`
		ReportType string     `json:"reportType"`
		Instances  []Instance `json:"instances"`
	}

	// Instance contains stats of an RDS instance
	Instance  struct {
		DBInstanceIdentifier string  `json:"dbInstanceIdentifier"`
		DBInstanceClass      string  `json:"dbInstanceClass"`
		AllocatedStorage     int64   `json:"allocatedStorage"`
		Engine               string  `json:"engine"`
		AvailabilityZone     string  `json:"availabilityZone"`
		MultiAZ              bool    `json:"multiAZ"`
		Cost                 float64 `json:"cost"`
		CpuAverage           float64 `json:"cpuAverage"`
		CpuPeak              float64 `json:"cpuPeak"`
		FreeSpaceMin         float64 `json:"freeSpaceMinimum"`
		FreeSpaceMax         float64 `json:"freeSpaceMaximum"`
		FreeSpaceAve         float64 `json:"freeSpaceAverage"`
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

// prepareResponseRdsDaily parses the results from elasticsearch and returns the RDS report
func prepareResponseRdsDaily(ctx context.Context, resRds *elastic.SearchResult, resCost *elastic.SearchResult) ([]Report, error) {
	rds  := ResponseRds{}
	cost := ResponseCost{}
	reports := []Report{}
	err := json.Unmarshal(*resRds.Aggregations["top_reports"], &rds.TopReports)
	if err != nil {
		return nil, err
	}
	if resCost != nil {
		json.Unmarshal(*resCost.Aggregations["accounts"], &cost.Accounts)
	}
	for _, account := range rds.TopReports.Buckets {
		if len(account.TopReportsHits.Hits.Hits) > 0 {
			report := account.TopReportsHits.Hits.Hits[0].Source
			report = addCostToReport(report, cost)
			reports = append(reports, report)
		}
	}
	return reports, nil
}

// prepareResponseRdsMonthly parses the results from elasticsearch and returns the RDS usage report
func prepareResponseRdsMonthly(ctx context.Context, resRds *elastic.SearchResult) ([]Report, error) {
	response := ResponseRds{}
	reports := []Report{}
	err := json.Unmarshal(*resRds.Aggregations["top_reports"], &response.TopReports)
	if err != nil {
		return nil, err
	}
	for _, account := range response.TopReports.Buckets {
		if len(account.TopReportsHits.Hits.Hits) > 0 {
			reports = append(reports, account.TopReportsHits.Hits.Hits[0].Source)
		}
	}
	return reports, nil
}

// compareInstanceCpu compares the Cpu average of two rds instances
// If inst1 is lower it returns 1, if inst2 is lower it returns -1, if they are equals it returns 0
func compareInstanceCpu(inst1, inst2 Instance) int {
	if inst1.CpuAverage < inst2.CpuAverage {
		return 1
	} else if inst1.CpuAverage > inst2.CpuAverage {
		return -1
	} else {
		return 0
	}
}

// compareInstanceFreeSpace compares the free space minimum of two rds instances
// If inst1 is higher it returns 1, if inst2 is higher it returns -1, if they are equals it returns 0
func compareInstanceFreeSpace(inst1, inst2 Instance) int {
	if inst1.FreeSpaceMin > inst2.FreeSpaceMin {
		return 1
	} else if inst1.FreeSpaceMin < inst2.FreeSpaceMin {
		return -1
	} else {
		return 0
	}
}

// compareInstance compares two rds instances depending on "by" parameter
// If inst1 is lower it returns 1, if inst2 is lower it returns -1, if they are equals it returns 0
func compareInstance(inst1, inst2 Instance, by string) int {
	switch by {
	case "cpu":
		return compareInstanceCpu(inst1, inst2)
	case "freespace":
		return compareInstanceFreeSpace(inst1, inst2)
	default:
		return 0
	}
}

// prepareResponseRdsUnused filter reports to get the top instances
func prepareResponseRdsUnused(params rdsUnusedQueryParams, reports []Report) (int, []Instance, error) {
	instances := []Instance{}
	for _, report := range reports {
		for _, instance := range report.Instances {
			instances = append(instances, instance)
		}
	}
	for i := 0; i < len(instances) - 1; i++ {
		if diff := compareInstance(instances[i], instances[i + 1], params.by); diff == -1 {
			tmp := instances[i]
			instances[i] = instances[i + 1]
			instances[i + 1] = tmp
			i -= 2
			if i < -1 {
				i = -1
			}
		}
	}
	if params.count >= 0 && params.count <= len(instances) {
		return http.StatusOK, instances[0:params.count], nil
	}
	return http.StatusOK, instances, nil
}
