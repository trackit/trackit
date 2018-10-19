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

package es

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws/usageReports/es"
	"gopkg.in/olivere/elastic.v5"
)

type (

	// Structure that allow to parse ES response for costs
	ResponseCost struct {
		Accounts struct {
			Buckets []struct {
				Key     string `json:"key"`
				Domains struct {
					Buckets []struct {
						Key  string `json:"key"`
						Cost struct {
							Value float64 `json:"value"`
						} `json:"cost"`
					} `json:"buckets"`
				} `json:"domains"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for ES usage report
	ResponseEs struct {
		TopReports struct {
			Buckets []struct {
				TopReportsHits struct {
					Hits struct {
						Hits []struct {
							Source es.Report `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"top_reports_hits"`
			} `json:"buckets"`
		} `json:"top_reports"`
	}
)

// prepareResponseEsDaily parses the results from elasticsearch and returns the RDS report
func prepareResponseEsDaily(ctx context.Context, resEs *elastic.SearchResult, resCost *elastic.SearchResult) ([]es.Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	parsedEs := ResponseEs{}
	cost := ResponseCost{}
	reports := []es.Report{}
	err := json.Unmarshal(*resEs.Aggregations["top_reports"], &parsedEs.TopReports)
	if err != nil {
		logger.Error("Error while unmarshaling ES elasticsearch response", err)
		return nil, err
	}
	if resCost != nil {
		err := json.Unmarshal(*resCost.Aggregations["accounts"], &cost.Accounts)
		if err != nil {
			logger.Error("Error while unmarshaling ES cost response", err)
			return nil, err
		}
	}
	for _, account := range parsedEs.TopReports.Buckets {
		if len(account.TopReportsHits.Hits.Hits) > 0 {
			report := account.TopReportsHits.Hits.Hits[0].Source
			report = addCostToReport(report, cost)
			reports = append(reports, report)
		}
	}
	return reports, nil
}

// prepareResponseEsMonthly parses the results from elasticsearch and returns the RDS usage report
func prepareResponseEsMonthly(ctx context.Context, resEs *elastic.SearchResult) ([]es.Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	response := ResponseEs{}
	reports := []es.Report{}
	err := json.Unmarshal(*resEs.Aggregations["top_reports"], &response.TopReports)
	if err != nil {
		logger.Error("Error while unmarshaling ES ES response", err)
		return nil, err
	}
	for _, account := range response.TopReports.Buckets {
		if len(account.TopReportsHits.Hits.Hits) > 0 {
			reports = append(reports, account.TopReportsHits.Hits.Hits[0].Source)
		}
	}
	return reports, nil
}

// addCostToReport adds cost for each instance based on billing data
func addCostToReport(report es.Report, costs ResponseCost) es.Report {
	for _, accounts := range costs.Accounts.Buckets {
		if accounts.Key != report.Account {
			continue
		}
		for _, domainCost := range accounts.Domains.Buckets {
			for i := range report.Domains {
				if report.Domains[i].CostDetail == nil {
					report.Domains[i].CostDetail = make(map[string]float64, 0)
				}
				if domainCost.Key == report.Domains[i].Arn {
					report.Domains[i].Cost += domainCost.Cost.Value
					report.Domains[i].CostDetail["domain"] += domainCost.Cost.Value
				}
			}
		}
	}
	return report
}

func isDomainUnused(domain es.Domain) bool {
	average := domain.CPUUtilizationAverage
	peak := domain.CPUUtiliztionPeak
	if peak >= 60.0 {
		return false
	} else if average >= 10 {
		return false
	}
	return true
}

func prepareResponseEsUnused(params esUnusedQueryParams, reports []es.Report) (int, []es.Domain, error) {
	domains := []es.Domain{}
	for _, report := range reports {
		for _, domain := range report.Domains {
			if isDomainUnused(domain) {
				domains = append(domains, domain)
			}
		}
	}
	sort.SliceStable(domains, func(i, j int) bool {
		return domains[i].Cost > domains[j].Cost
	})
	if params.count >= 0 && params.count <= len(domains) {
		return http.StatusOK, domains[0:params.count], nil
	}
	return http.StatusOK, domains, nil
}
