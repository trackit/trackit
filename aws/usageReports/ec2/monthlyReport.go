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

package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/ec2Reports"
)

// getElasticSearchEc2Instance prepares and run the request to retrieve the a report of an instance
// It will return the data and an error.
func getElasticSearchEc2Instance(ctx context.Context, account, instance string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("instance.id", instance))
	search := client.Search().Index(index).Size(1).Query(query)
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, errors.GetErrorMessage(ctx, err)
	}
	return res, nil
}

// getInstanceInfoFromEs gets information about an instance from previous report to put it in the new report
func getInstanceInfoFromES(ctx context.Context, instance common.CostPerResource, account string, userId int) ec2Reports.Instance {
	var docType ec2Reports.InstanceReport
	var inst = ec2Reports.Instance{
		InstanceBase: ec2Reports.InstanceBase{
			Id:         instance.Resource,
			Region:     "N/A",
			State:      "N/A",
			Purchasing: "N/A",
			KeyPair:    "",
			Type:       "N/A",
			Platform:   "Linux/UNIX",
		},
		Tags:  make([]common.Tag, 0),
		Costs: make(map[string]float64),
		Stats: ec2Reports.Stats{
			Cpu: ec2Reports.Cpu{
				Average: -1,
				Peak:    -1,
			},
			Network: ec2Reports.Network{
				In:  -1,
				Out: -1,
			},
			Volumes: make([]ec2Reports.Volume, 0),
		},
	}
	inst.Costs["instance"] = instance.Cost
	res, err := getElasticSearchEc2Instance(ctx, account, instance.Resource,
		es.Client, es.IndexNameForUserId(userId, ec2Reports.Model.IndexSuffix))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			inst.Region = docType.Instance.Region
			inst.Purchasing = docType.Instance.Purchasing
			inst.KeyPair = docType.Instance.KeyPair
			inst.Type = docType.Instance.Type
			inst.Platform = docType.Instance.Platform
			inst.Tags = docType.Instance.Tags
		}
	}
	return inst
}

// fetchMonthlyInstancesList sends in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances and getInstanceStats.
func fetchMonthlyInstancesList(ctx context.Context, creds *credentials.Credentials, inst common.CostPerResource,
	account, region string, instanceChan chan ec2Reports.Instance, startDate, endDate time.Time, userId int) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	desc := ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice([]string{inst.Resource})}
	instances, err := svc.DescribeInstances(&desc)
	if err != nil {
		instanceChan <- getInstanceInfoFromES(ctx, inst, account, userId)
		return err
	}
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			stats := getInstanceStats(ctx, instance, sess, startDate, endDate)
			costs := make(map[string]float64, 1)
			costs["instance"] = inst.Cost
			instanceChan <- ec2Reports.Instance{
				InstanceBase: ec2Reports.InstanceBase{
					Id:         aws.StringValue(instance.InstanceId),
					Region:     aws.StringValue(instance.Placement.AvailabilityZone),
					State:      aws.StringValue(instance.State.Name),
					Purchasing: getPurchasingOption(instance),
					KeyPair:    aws.StringValue(instance.KeyName),
					Type:       aws.StringValue(instance.InstanceType),
					Platform:   getPlatformName(aws.StringValue(instance.Platform)),
				},
				Tags:  getInstanceTag(instance.Tags),
				Costs: costs,
				Stats: stats,
			}
		}
	}
	return nil
}

// getEc2Metrics gets credentials, accounts and region to fetch EC2 instances stats
func fetchMonthlyInstancesStats(ctx context.Context, instances []common.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) ([]ec2Reports.InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorInstanceStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return nil, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return nil, err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return nil, err
	}
	instanceChans := make([]<-chan ec2Reports.Instance, 0, len(regions))
	for _, instance := range instances {
		for _, region := range regions {
			if strings.Contains(instance.Region, region) {
				instanceChan := make(chan ec2Reports.Instance)
				go fetchMonthlyInstancesList(ctx, creds, instance, account, region, instanceChan, startDate, endDate, aa.UserId)
				instanceChans = append(instanceChans, instanceChan)
			}
		}
	}
	instancesList := make([]ec2Reports.InstanceReport, 0)
	for instance := range merge(instanceChans...) {
		instancesList = append(instancesList, ec2Reports.InstanceReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Instance: instance,
		})
	}
	return instancesList, nil
}

// filterInstancesCosts filters instances, cloudwatch and volumes of EC2 instances costs
func filterInstancesCosts(ec2Cost, cloudwatchCost []common.CostPerResource) ([]common.CostPerResource, []common.CostPerResource, []common.CostPerResource) {
	newInstance := make([]common.CostPerResource, 0)
	newVolume := make([]common.CostPerResource, 0)
	newCloudWatch := make([]common.CostPerResource, 0)
	for _, instance := range ec2Cost {
		if len(instance.Resource) == 19 && strings.HasPrefix(instance.Resource, "i-") {
			newInstance = append(newInstance, instance)
		}
		if len(instance.Resource) == 21 && strings.HasPrefix(instance.Resource, "vol-") {
			newVolume = append(newVolume, instance)
		}
	}
	for _, instance := range cloudwatchCost {
		for _, cost := range newInstance {
			if strings.Contains(instance.Resource, cost.Resource) {
				newCloudWatch = append(newCloudWatch, instance)
			}
		}
	}
	return newInstance, newVolume, newCloudWatch
}

func addCostToInstances(instances []ec2Reports.InstanceReport, costVolume, costCloudWatch []common.CostPerResource) []ec2Reports.InstanceReport {
	for i, instance := range instances {
		for _, volume := range instance.Instance.Stats.Volumes {
			for _, costPerVolume := range costVolume {
				if volume.Id == costPerVolume.Resource {
					instances[i].Instance.Costs[volume.Id] += costPerVolume.Cost
				}
			}
		}
		for _, cloudWatch := range costCloudWatch {
			if strings.Contains(cloudWatch.Resource, instance.Instance.Id) {
				instances[i].Instance.Costs["cloudwatch"] += cloudWatch.Cost
			}
		}
	}
	return instances
}

func getEc2Recommendations(instances []ec2Reports.InstanceReport) []ec2Reports.InstanceReport {
	for idx, report := range instances {
		instances[idx].Instance.Recommendation = getEC2RecommendationTypeReason(report.Instance)
	}
	return instances
}

// PutEc2MonthlyReport puts a monthly report of EC2 instance in ES
func PutEc2MonthlyReport(ctx context.Context, ec2Cost, cloudWatchCost []common.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting EC2 monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance, costVolume, costCloudWatch := filterInstancesCosts(ec2Cost, cloudWatchCost)
	if len(costInstance) == 0 {
		logger.Info("No EC2 instances found in billing data.", nil)
		return false, nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, ec2Reports.Model.IndexSuffix)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an EC2 monthly report", nil)
		return false, nil
	}
	instances, err := fetchMonthlyInstancesStats(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	instances = addCostToInstances(instances, costVolume, costCloudWatch)
	instances = getEc2Recommendations(instances)
	err = importInstancesToEs(ctx, aa, instances)
	if err != nil {
		return false, err
	}
	return true, nil
}
