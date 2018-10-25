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
	"sort"
	"strings"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/usageReports/rds"
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

	// Structure that allow to parse ES response for RDS Monthly instances
	ResponseRdsMonthly struct {
		Accounts struct {
			Buckets []struct {
				Instances struct {
					Hits struct {
						Hits []struct {
							Instance rds.InstanceReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"instances"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for RDS Daily instances
	ResponseRdsDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time      string `json:"key_as_string"`
						Instances struct {
							Hits struct {
								Hits []struct {
									Instance rds.InstanceReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"instances"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}
)

// addCostToInstance adds cost for an instance based on billing data
func addCostToInstance(instance rds.InstanceReport, costs ResponseCost) rds.InstanceReport {
	if instance.Instance.Costs == nil {
		instance.Instance.Costs = make(map[string]float64, 0)
	}
	for _, accounts := range costs.Accounts.Buckets {
		if accounts.Key != instance.Account {
			continue
		}
		for _, instanceCost := range accounts.Instances.Buckets {
			// format in billing data for an RDS instance is: "arn:aws:rds:us-west-2:394125495069:db:instancename"
			// so i get the 7th element of the split by ":"
			split := strings.Split(instanceCost.Key, ":")
			if len(split) == 7 && split[6] == instance.Instance.DBInstanceIdentifier {
				instance.Instance.Costs["instance"] += instanceCost.Cost.Value
			}
		}
		return instance
	}
	return instance
}

// prepareResponseRdsMonthly parses the results from elasticsearch and returns an array of RDS daily instances report
func prepareResponseRdsDaily(ctx context.Context, resRds *elastic.SearchResult, resCost *elastic.SearchResult) ([]rds.InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedRds ResponseRdsDaily
	var parsedCost ResponseCost
	instances := make([]rds.InstanceReport, 0)
	err := json.Unmarshal(*resRds.Aggregations["accounts"], &parsedRds.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES RDS response", err)
		return nil, err
	}
	if resCost != nil {
		err = json.Unmarshal(*resCost.Aggregations["accounts"], &parsedCost.Accounts)
		if err != nil {
			logger.Error("Error while unmarshaling ES cost response", err)
		}
	}
	for _, account := range parsedRds.Accounts.Buckets {
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
					instances = append(instances, instance.Instance)
				}
			}
		}
	}
	return instances, nil
}

// prepareResponseRdsMonthly parses the results from elasticsearch and returns an array of RDS monthly instances report
func prepareResponseRdsMonthly(ctx context.Context, resRds *elastic.SearchResult) ([]rds.InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseRdsMonthly
	instances := make([]rds.InstanceReport, 0)
	err := json.Unmarshal(*resRds.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES RDS response", err)
		return nil, err
	}
	for _, account := range response.Accounts.Buckets {
		for _, instance := range account.Instances.Hits.Hits {
			instances = append(instances, instance.Instance)
		}
	}
	return instances, nil
}

func isInstanceUnused(instance rds.Instance) bool {
	average := instance.Stats.Cpu.Average
	peak := instance.Stats.Cpu.Peak
	if peak >= 60.0 {
		return false
	} else if average >= 10.0 {
		return false
	}
	return true
}

// prepareResponseRdsUnused filter reports to get the unused instances sorted by cost
func prepareResponseRdsUnused(params RdsUnusedQueryParams, instances []rds.InstanceReport) (int, []rds.InstanceReport, error) {
	unusedInstances := make([]rds.InstanceReport, 0)
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
	if params.count >= 0 && params.count <= len(unusedInstances) {
		return http.StatusOK, unusedInstances[0:params.count], nil
	}
	return http.StatusOK, unusedInstances, nil
}
