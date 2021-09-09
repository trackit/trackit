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

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/aws/usageReports/es"
	"github.com/trackit/trackit/errors"
)

type (

	// ResponseCost allows us to parse an ES response for costs
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

	// ResponseEsMonthly allows us to parse an ES response for ES monthly domains
	ResponseEsMonthly struct {
		Accounts struct {
			Buckets []struct {
				Domains struct {
					Hits struct {
						Hits []struct {
							Domain es.DomainReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"domains"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// ResponseEsDaily allows us to parse an ES response for ES daily domains
	ResponseEsDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time    string `json:"key_as_string"`
						Domains struct {
							Hits struct {
								Hits []struct {
									Domain es.DomainReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"domains"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// DomainReport represents the report with all the information for an ES domain report
	DomainReport struct {
		utils.ReportBase
		Domain Domain `json:"domain"`
	}

	// Domain represents all the informations of an ES domain.
	Domain struct {
		es.DomainBase
		Tags  map[string]string  `json:"tags"`
		Costs map[string]float64 `json:"costs"`
		Stats es.Stats           `json:"stats"`
	}
)

func getEsDomainReportResponse(oldDomain es.DomainReport) DomainReport {
	tags := make(map[string]string)
	for _, tag := range oldDomain.Domain.Tags {
		tags[tag.Key] = tag.Value
	}
	newDomain := DomainReport{
		ReportBase: oldDomain.ReportBase,
		Domain: Domain{
			DomainBase: oldDomain.Domain.DomainBase,
			Tags:       tags,
			Costs:      oldDomain.Domain.Costs,
			Stats:      oldDomain.Domain.Stats,
		},
	}
	return newDomain
}

// addCostToDomain adds cost for each domain based on billing data
func addCostToDomain(domain es.DomainReport, costs ResponseCost) es.DomainReport {
	domain.Domain.Costs = make(map[string]float64, 1)
	for _, accounts := range costs.Accounts.Buckets {
		if accounts.Key != domain.Account {
			continue
		}
		for _, domainCost := range accounts.Domains.Buckets {
			if domainCost.Key == domain.Domain.Arn {
				domain.Domain.Costs["domain"] += domainCost.Cost.Value
			}
		}
	}
	return domain
}

// prepareResponseEsDaily parses the results from elasticsearch and returns an array of ES daily domains report
func prepareResponseEsDaily(ctx context.Context, resEc2 *elastic.SearchResult, resCost *elastic.SearchResult) ([]DomainReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedEc2 ResponseEsDaily
	var parsedCost ResponseCost
	domains := make([]DomainReport, 0)
	err := json.Unmarshal(*resEc2.Aggregations["accounts"], &parsedEc2.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES daily response", err)
		return nil, err
	}
	if resCost != nil {
		err = json.Unmarshal(*resCost.Aggregations["accounts"], &parsedCost.Accounts)
		if err != nil {
			logger.Error("Error while unmarshaling ES cost response", err)
		}
	}
	for _, account := range parsedEc2.Accounts.Buckets {
		var lastDate = ""
		for _, date := range account.Dates.Buckets {
			if date.Time > lastDate {
				lastDate = date.Time
			}
		}
		for _, date := range account.Dates.Buckets {
			if date.Time == lastDate {
				for _, domain := range date.Domains.Hits.Hits {
					domain.Domain = addCostToDomain(domain.Domain, parsedCost)
					domains = append(domains, getEsDomainReportResponse(domain.Domain))
				}
			}
		}
	}
	return domains, nil
}

// prepareResponseEsMonthly parses the results from elasticsearch and returns an array of ES monthly domains report
func prepareResponseEsMonthly(ctx context.Context, resEc2 *elastic.SearchResult) ([]DomainReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseEsMonthly
	domains := make([]DomainReport, 0)
	err := json.Unmarshal(*resEc2.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES EC2 response", err)
		return nil, errors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, domain := range account.Domains.Hits.Hits {
			domains = append(domains, getEsDomainReportResponse(domain.Domain))
		}
	}
	return domains, nil
}

func isDomainUnused(domain Domain) bool {
	average := domain.Stats.Cpu.Average
	peak := domain.Stats.Cpu.Peak
	if average == -1 || peak == -1 {
		return false
	} else if peak >= 60.0 {
		return false
	} else if average >= 10 {
		return false
	}
	return true
}

func prepareResponseEsUnused(params EsUnusedQueryParams, domains []DomainReport) (int, []DomainReport, error) {
	unusedDomains := make([]DomainReport, 0)
	for _, domain := range domains {
		if isDomainUnused(domain.Domain) {
			unusedDomains = append(unusedDomains, domain)
		}
	}
	sort.SliceStable(unusedDomains, func(i, j int) bool {
		var cost1, cost2 float64
		for _, cost := range unusedDomains[i].Domain.Costs {
			cost1 += cost
		}
		for _, cost := range unusedDomains[j].Domain.Costs {
			cost2 += cost
		}
		return cost1 > cost2
	})
	if params.Count >= 0 && params.Count <= len(domains) {
		return http.StatusOK, unusedDomains[0:params.Count], nil
	}
	return http.StatusOK, unusedDomains, nil
}
