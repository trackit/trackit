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

package ec2

import (
	"strings"
	"context"
	"net/http"
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

	// Structure that allow to parse ES response for EC2 usage report
	ResponseEc2 struct {
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
		Account    string     `json:"account"`
		ReportDate string     `json:"reportDate"`
		ReportType string     `json:"reportType"`
		Instances  []Instance `json:"instances"`
	}

	// Instance contains stats of an EC2 instance
	Instance struct {
		Id         string             `json:"id"`
		Region     string             `json:"region"`
		State      string             `json:"state"`
		Purchasing string             `json:"purchasing"`
		Cost       float64            `json:"cost"`
		CpuAverage float64            `json:"cpuAverage"`
		CpuPeak    float64            `json:"cpuPeak"`
		NetworkIn  int64              `json:"networkIn"`
		NetworkOut int64              `json:"networkOut"`
		IORead     map[string]float64 `json:"ioRead"`
		IOWrite    map[string]float64 `json:"ioWrite"`
		KeyPair    string             `json:"keyPair"`
		Type       string             `json:"type"`
		Tags       map[string]string  `json:"tags"`
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
				if strings.Contains(instance.Key, report.Instances[i].Id) {
					report.Instances[i].Cost += instance.Cost.Value
				}
				for volume := range report.Instances[i].IOWrite {
					if volume == instance.Key {
						report.Instances[i].Cost += instance.Cost.Value
					}
				}
				for volume := range report.Instances[i].IORead {
					if volume == instance.Key {
						report.Instances[i].Cost += instance.Cost.Value
					}
				}
			}
		}
	}
	return report
}

// prepareResponseEc2 parses the results from elasticsearch and returns the EC2 usage report
func prepareResponseEc2(ctx context.Context, resEc2 *elastic.SearchResult, resCost *elastic.SearchResult) ([]Report, error) {
	ec2  := ResponseEc2{}
	cost := ResponseCost{}
	reports := []Report{}
	err := json.Unmarshal(*resEc2.Aggregations["top_reports"], &ec2.TopReports)
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

// prepareResponseEc2Monthly parses the results from elasticsearch and returns the EC2 usage report
func prepareResponseEc2Monthly(ctx context.Context, resEc2 *elastic.SearchResult) ([]Report, error) {
	response := ResponseEc2{}
	reports := []Report{}
	err := json.Unmarshal(*resEc2.Aggregations["top_reports"], &response.TopReports)
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

func isInstanceUnused(instance Instance) bool {
	average := instance.CpuAverage
	peak    := instance.CpuPeak
	if peak >= 60.0 {
		return false
	} else if average >= 10.0 {
		return false
	}
	return true
}

// prepareResponseEc2Unused filter reports to get the top instances
func prepareResponseEc2Unused(params ec2UnusedQueryParams, reports []Report) (int, []Instance, error) {
	instances := []Instance{}
	for _, report := range reports {
		for _, instance := range report.Instances {
			if isInstanceUnused(instance) {
				instances = append(instances, instance)
			}
		}
	}
	for i := 0; i < len(instances) - 1; i++ {
		if instances[i].Cost < instances[i + 1].Cost {
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
