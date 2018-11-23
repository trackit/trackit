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

package reservedInstances

import (
	"context"
	"encoding/json"
	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/aws/usageReports/reservedInstances"
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
							Instance reservedInstances.InstanceReport `json:"_source"`
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
									Instance reservedInstances.InstanceReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"instances"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// InstanceReport has all the information of an EC2 instance report
	InstanceReport struct {
		utils.ReportBase
		Instance Instance `json:"instance"`
	}

	// Instance contains the information of an EC2 instance
	Instance struct {
		reservedInstances.InstanceBase
		Tags  map[string]string  `json:"tags"`
	}
)

func getEc2InstanceReportResponse(oldInstance reservedInstances.InstanceReport) InstanceReport {
	tags := make(map[string]string, 0)
	for _, tag := range oldInstance.Instance.Tags {
		tags[tag.Key] = tag.Value
	}
	newInstance := InstanceReport{
		ReportBase: oldInstance.ReportBase,
		Instance: Instance{
			InstanceBase: oldInstance.Instance.InstanceBase,
			Tags:         tags,
		},
	}
	return newInstance
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
		return nil, errors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, instance := range account.Instances.Hits.Hits {
			instances = append(instances, getEc2InstanceReportResponse(instance.Instance))
		}
	}
	return instances, nil
}
