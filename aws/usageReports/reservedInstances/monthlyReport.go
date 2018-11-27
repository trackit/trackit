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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/config"
)


// fetchMonthlyInstancesList sends in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances and getInstanceStats.
func fetchMonthlyInstancesList(ctx context.Context, creds *credentials.Credentials, inst utils.CostPerResource,
	region string, instanceChan chan Instance, startDate, endDate time.Time) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	instances, err := svc.DescribeReservedInstances(nil)
	if err != nil {
		return err
	}
	for _, reservation := range instances.ReservedInstances {
		costs := make(map[string]float64, 0)
		instanceChan <- Instance{
			InstanceBase: InstanceBase{
				Id:              aws.StringValue(reservation.ReservedInstancesId),
				Region:          aws.StringValue(reservation.AvailabilityZone),
				Type:            aws.StringValue(reservation.InstanceType),
				FixedPrice:      aws.Float64Value(reservation.FixedPrice),
				UsagePrice:      aws.Float64Value(reservation.UsagePrice),
				Duration:        aws.Int64Value(reservation.Duration),
				Start:           aws.TimeValue(reservation.Start),
				End:             aws.TimeValue(reservation.End),
				InstanceCount:   aws.Int64Value(reservation.InstanceCount),
				InstanceTenancy: aws.StringValue(reservation.InstanceTenancy),
			},
			Tags:  getInstanceTag(reservation.Tags),
			Costs: costs,
		}
	}
	return nil
}

// getReservedInstancesMetrics gets credentials, accounts and region to fetch ReservedInstances instances stats
func fetchMonthlyInstancesStats(ctx context.Context, instances []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) ([]InstanceReport, error) {
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
	instanceChans := make([]<-chan Instance, 0, len(regions))
	for _, instance := range instances {
		for _, region := range regions {
			if strings.Contains(instance.Region, region) {
				instanceChan := make(chan Instance)
				go fetchMonthlyInstancesList(ctx, creds, instance, region, instanceChan, startDate, endDate)
				instanceChans = append(instanceChans, instanceChan)
			}
		}
	}
	instancesList := make([]InstanceReport, 0)
	for instance := range merge(instanceChans...) {
		instancesList = append(instancesList, InstanceReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Instance: instance,
		})
	}
	return instancesList, nil
}

// filterInstancesCosts filters instances, cloudwatch and volumes of ReservedInstances instances costs
func filterInstancesCosts(reservedInstancesCost, cloudwatchCost []utils.CostPerResource) ([]utils.CostPerResource, []utils.CostPerResource, []utils.CostPerResource) {
	newInstance := make([]utils.CostPerResource, 0)
	newVolume := make([]utils.CostPerResource, 0)
	newCloudWatch := make([]utils.CostPerResource, 0)
	for _, instance := range reservedInstancesCost {
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

func addCostToInstances(instances []InstanceReport, costVolume, costCloudWatch []utils.CostPerResource) []InstanceReport {
	for i, instance := range instances {
		//for _, volume := range instance.Instance.Stats.Volumes {
		//	for _, costPerVolume := range costVolume {
		//		if volume.Id == costPerVolume.Resource {
		//			instances[i].Instance.Costs[volume.Id] += costPerVolume.Cost
		//		}
		//	}
		//}
		for _, cloudWatch := range costCloudWatch {
			if strings.Contains(cloudWatch.Resource, instance.Instance.Id) {
				instances[i].Instance.Costs["cloudwatch"] += cloudWatch.Cost
			}
		}
	}
	return instances
}

// PutReservedInstancesMonthlyReport puts a monthly report of ReservedInstances instance in ES
func PutReservedInstancesMonthlyReport(ctx context.Context, reservedInstancesCost, cloudWatchCost []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting ReservedInstances monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance, costVolume, costCloudWatch := filterInstancesCosts(reservedInstancesCost, cloudWatchCost)
	if len(costInstance) == 0 {
		logger.Info("No ReservedInstances instances found in billing data.", nil)
		return false, nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixReservedInstancesReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an ReservedInstances monthly report", nil)
		return false, nil
	}
	instances, err := fetchMonthlyInstancesStats(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	instances = addCostToInstances(instances, costVolume, costCloudWatch)
	err = importInstancesToEs(ctx, aa, instances)
	if err != nil {
		return false, err
	}
	return true, nil
}
