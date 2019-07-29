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

package elasticache

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/aws/usageReports/elasticache"
	"github.com/trackit/trackit/errors"
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

	// Structure that allow to parse ES response for ElastiCache Monthly instances
	ResponseElastiCacheMonthly struct {
		Accounts struct {
			Buckets []struct {
				Instances struct {
					Hits struct {
						Hits []struct {
							Instance elasticache.InstanceReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"instances"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for ElastiCache Daily instances
	ResponseElastiCacheDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time      string `json:"key_as_string"`
						Instances struct {
							Hits struct {
								Hits []struct {
									Instance elasticache.InstanceReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"instances"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// InstanceReport has all the information of an ElastiCache instance report
	InstanceReport struct {
		utils.ReportBase
		Instance Instance `json:"instance"`
	}

	// Instance contains the information of an ElastiCache instance
	Instance struct {
		elasticache.InstanceBase
		Tags  map[string]string  `json:"tags"`
		Costs map[string]float64 `json:"costs"`
		Stats elasticache.Stats  `json:"stats"`
	}
)

func getElastiCacheInstanceReportResponse(oldInstance elasticache.InstanceReport) InstanceReport {
	tags := make(map[string]string, 0)
	for _, tag := range oldInstance.Instance.Tags {
		tags[tag.Key] = tag.Value
	}
	newInstance := InstanceReport{
		ReportBase: oldInstance.ReportBase,
		Instance: Instance{
			InstanceBase: oldInstance.Instance.InstanceBase,
			Tags:         tags,
			Costs:        oldInstance.Instance.Costs,
			Stats:        oldInstance.Instance.Stats,
		},
	}
	return newInstance
}

// addCostToInstance adds a cost for an instance based on billing data
func addCostToInstance(instance elasticache.InstanceReport, costs ResponseCost) elasticache.InstanceReport {
	if instance.Instance.Costs == nil {
		instance.Instance.Costs = make(map[string]float64, 0)
	}
	for _, accounts := range costs.Accounts.Buckets {
		if accounts.Key != instance.Account {
			continue
		}
		for _, instanceCost := range accounts.Instances.Buckets {
			// format ARN for an ElastiCache instance is:
			// "arn:aws:elasticache:[region]:[aws_id]:cluster:[cluster name]"
			split := strings.Split(instanceCost.Key, ":")
			if len(split) == 7 && split[2] == "elasticache" && split[6] == instance.Instance.Id {
				instance.Instance.Costs[split[6]] += instanceCost.Cost.Value
			}
		}
		return instance
	}
	return instance
}

// prepareResponseElastiCacheDaily parses the results from elasticsearch and returns an array of ElastiCache daily instances report
func prepareResponseElastiCacheDaily(ctx context.Context, resElastiCache *elastic.SearchResult, resCost *elastic.SearchResult) ([]InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedElastiCache ResponseElastiCacheDaily
	var parsedCost ResponseCost
	instances := make([]InstanceReport, 0)
	err := json.Unmarshal(*resElastiCache.Aggregations["accounts"], &parsedElastiCache.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES ElastiCache response", err)
		return nil, err
	}
	if resCost != nil {
		err = json.Unmarshal(*resCost.Aggregations["accounts"], &parsedCost.Accounts)
		if err != nil {
			logger.Error("Error while unmarshaling ES cost response", err)
		}
	}
	for _, account := range parsedElastiCache.Accounts.Buckets {
		var lastDate = ""
		for _, date := range account.Dates.Buckets {
			if date.Time > lastDate {
				lastDate = date.Time
			}
		}
		for _, date := range account.Dates.Buckets {
			if date.Time == lastDate {
				for _, instance := range date.Instances.Hits.Hits {
					instance.Instance = addCostToInstance(instance.Instance, parsedCost)
					instances = append(instances, getElastiCacheInstanceReportResponse(instance.Instance))
				}
			}
		}
	}
	return instances, nil
}

// prepareResponseElastiCacheMonthly parses the results from elasticsearch and returns an array of ElastiCache monthly instances report
func prepareResponseElastiCacheMonthly(ctx context.Context, resElastiCache *elastic.SearchResult) ([]InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseElastiCacheMonthly
	instances := make([]InstanceReport, 0)
	err := json.Unmarshal(*resElastiCache.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES ElastiCache response", err)
		return nil, errors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, instance := range account.Instances.Hits.Hits {
			instances = append(instances, getElastiCacheInstanceReportResponse(instance.Instance))
		}
	}
	return instances, nil
}

func isInstanceUnused(instance Instance) bool {
	average := instance.Stats.Cpu.Average
	peak := instance.Stats.Cpu.Peak
	if peak >= 60.0 {
		return false
	} else if average >= 10.0 {
		return false
	}
	return true
}

// prepareResponseElastiCacheUnused filter reports to get the unused instances sorted by cost
func prepareResponseElastiCacheUnused(params ElastiCacheUnusedQueryParams, instances []InstanceReport) (int, []InstanceReport, error) {
	unusedInstances := make([]InstanceReport, 0)
	for _, instance := range instances {
		if isInstanceUnused(instance.Instance) {
			unusedInstances = append(unusedInstances, instance)
		}
	}
	sort.SliceStable(unusedInstances, func(i, j int) bool {
		var cost1, cost2 float64
		for _, cost := range unusedInstances[i].Instance.Costs {
			cost1 += cost
		}
		for _, cost := range unusedInstances[j].Instance.Costs {
			cost2 += cost
		}
		return cost1 > cost2
	})
	if params.Count >= 0 && params.Count <= len(unusedInstances) {
		return http.StatusOK, unusedInstances[0:params.Count], nil
	}
	return http.StatusOK, unusedInstances, nil
}
