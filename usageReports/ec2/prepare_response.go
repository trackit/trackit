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
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/usageReports/ec2"
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

	// Structure that allow to parse ES response for EC2 Monthly instances
	ResponseEc2Monthly struct {
		Accounts struct {
			Buckets []struct {
				Instances struct {
					Hits struct {
						Hits []struct {
							Instance ec2.InstanceReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"instances"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for EC2 Daily instances
	ResponseEc2Daily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time      string `json:"key_as_string"`
						Instances struct {
							Hits struct {
								Hits []struct {
									Instance ec2.InstanceReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"instances"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// InstanceReport has all the information of an EC2 instance
	InstanceReport struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		ReportType string    `json:"reportType"`
		Instance   Instance  `json:"instance"`
	}

	// Instance contains the information of an EC2 instance
	Instance struct {
		Id         string             `json:"id"`
		Region     string             `json:"region"`
		State      string             `json:"state"`
		Purchasing string             `json:"purchasing"`
		KeyPair    string             `json:"keyPair"`
		Type       string             `json:"type"`
		Tags       map[string]string  `json:"tags"`
		Costs      map[string]float64 `json:"costs"`
		Stats      Stats              `json:"stats"`
	}

	// Stats contains statistics of an instance get on CloudWatch
	Stats struct {
		Cpu     Cpu     `json:"cpu"`
		Network Network `json:"network"`
		Volumes Volumes `json:"volumes"`
	}

	// Cpu contains cpu statistics of an instance
	Cpu struct {
		Average float64 `json:"average"`
		Peak    float64 `json:"peak"`
	}

	// Network contains network statistics of an instance
	Network struct {
		In  float64 `json:"in"`
		Out float64 `json:"out"`
	}

	// Volume contains information about EBS volumes
	Volumes struct {
		Read  map[string]float64 `json:"read"`
		Write map[string]float64 `json:"write"`
	}
)

func getEc2InstanceReportResponse(oldInstance ec2.InstanceReport) InstanceReport {
	tags := make(map[string]string, 0)
	for _, tag := range oldInstance.Instance.Tags {
		tags[tag.Key] = tag.Value
	}
	read := make(map[string]float64, 0)
	write := make(map[string]float64, 0)
	for _, volume := range oldInstance.Instance.Stats.Volumes {
		read[volume.Id] = volume.Read
		write[volume.Id] = volume.Write
	}
	newInstance := InstanceReport{
		Account:    oldInstance.Account,
		ReportType: oldInstance.ReportType,
		ReportDate: oldInstance.ReportDate,
		Instance: Instance{
			Id:         oldInstance.Instance.Id,
			Region:     oldInstance.Instance.Region,
			State:      oldInstance.Instance.State,
			Purchasing: oldInstance.Instance.Purchasing,
			KeyPair:    oldInstance.Instance.KeyPair,
			Type:       oldInstance.Instance.Type,
			Tags:       tags,
			Costs:      oldInstance.Instance.Costs,
			Stats: Stats{
				Cpu: Cpu{
					Average: oldInstance.Instance.Stats.Cpu.Average,
					Peak:    oldInstance.Instance.Stats.Cpu.Peak,
				},
				Network: Network{
					In:  oldInstance.Instance.Stats.Network.In,
					Out: oldInstance.Instance.Stats.Network.Out,
				},
				Volumes: Volumes{
					Read:  read,
					Write: write,
				},
			},
		},
	}
	return newInstance
}

// addCostToInstance adds a cost for an instance based on billing data
func addCostToInstance(instance ec2.InstanceReport, costs ResponseCost) ec2.InstanceReport {
	if instance.Instance.Costs == nil {
		instance.Instance.Costs = make(map[string]float64, 0)
	}
	for _, accounts := range costs.Accounts.Buckets {
		if accounts.Key != instance.Account {
			continue
		}
		for _, instanceCost := range accounts.Instances.Buckets {
			if strings.Contains(instanceCost.Key, instance.Instance.Id) {
				if len(instanceCost.Key) == 19 && strings.HasPrefix(instanceCost.Key, "i-") {
					instance.Instance.Costs["instance"] += instanceCost.Cost.Value
				} else {
					instance.Instance.Costs["cloudwatch"] += instanceCost.Cost.Value
				}
			}
			for _, volume := range instance.Instance.Stats.Volumes {
				if volume.Id == instanceCost.Key {
					instance.Instance.Costs[volume.Id] += instanceCost.Cost.Value
				}
			}
		}
		return instance
	}
	return instance
}

// prepareResponseEc2Daily parses the results from elasticsearch and returns an array of EC2 daily instances report
func prepareResponseEc2Daily(ctx context.Context, resEc2 *elastic.SearchResult, resCost *elastic.SearchResult) ([]InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedEc2 ResponseEc2Daily
	var parsedCost ResponseCost
	instances := make([]InstanceReport, 0)
	err := json.Unmarshal(*resEc2.Aggregations["accounts"], &parsedEc2.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES EC2 response", err)
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
				for _, instance := range date.Instances.Hits.Hits {
					instance.Instance = addCostToInstance(instance.Instance, parsedCost)
					instances = append(instances, getEc2InstanceReportResponse(instance.Instance))
				}
			}
		}
	}
	return instances, nil
}

// prepareResponseEc2Monthly parses the results from elasticsearch and returns an array of EC2 monthly instances report
func prepareResponseEc2Monthly(ctx context.Context, resEc2 *elastic.SearchResult) ([]InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseEc2Monthly
	instances := make([]InstanceReport, 0)
	err := json.Unmarshal(*resEc2.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES EC2 response", err)
		return nil, err
	}
	for _, account := range response.Accounts.Buckets {
		for _, instance := range account.Instances.Hits.Hits {
			instances = append(instances, getEc2InstanceReportResponse(instance.Instance))
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

// prepareResponseEc2Unused filter reports to get the unused instances sorted by cost
func prepareResponseEc2Unused(params Ec2UnusedQueryParams, instances []InstanceReport) (int, []InstanceReport, error) {
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
	if params.count >= 0 && params.count <= len(unusedInstances) {
		return http.StatusOK, unusedInstances[0:params.count], nil
	}
	return http.StatusOK, unusedInstances, nil
}
