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
// and filled by DescribeInstances, getAccountID and getInstanceStats.
func fetchMonthlyInstancesList(ctx context.Context, creds *credentials.Credentials, instList []utils.CostPerResource,
	region string, instanceChan chan Instance, startDate, endDate time.Time) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	for _, inst := range instList {
		desc := ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice([]string{inst.Resource})}
		instances, err := svc.DescribeInstances(&desc)
		if err != nil {
			continue
		}
		for _, reservation := range instances.Reservations {
			for _, instance := range reservation.Instances {
				stats := getInstanceStats(ctx, instance, sess, startDate, endDate)
				detail := make(map[string]float64, 0)
				detail["instance"] = inst.Cost
				instanceChan <- Instance{
					Id:         aws.StringValue(instance.InstanceId),
					Region:     aws.StringValue(instance.Placement.AvailabilityZone),
					KeyPair:    aws.StringValue(instance.KeyName),
					Tags:       getInstanceTag(instance.Tags),
					Type:       aws.StringValue(instance.InstanceType),
					State:      aws.StringValue(instance.State.Name),
					Purchasing: getPurchasingOption(instance),
					CpuAverage: stats.CpuAverage,
					CpuPeak:    stats.CpuPeak,
					NetworkIn:  stats.NetworkIn,
					NetworkOut: stats.NetworkOut,
					IORead:     stats.IORead,
					IOWrite:    stats.IOWrite,
					Cost:       inst.Cost,
					CostDetail: detail,
				}
			}
		}
	}
	return nil
}

// getEc2Metrics gets credentials, accounts and region to fetch EC2 instances stats
func fetchMonthlyInstancesStats(ctx context.Context, instances []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorInstanceStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return Report{}, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return Report{}, err
	}
	report := Report{
		Account:    account,
		ReportDate: startDate,
		ReportType: "monthly",
		Instances:  make([]Instance, 0),
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return Report{}, err
	}
	instanceChans := make([]<-chan Instance, 0, len(regions))
	for _, region := range regions {
		instanceChan := make(chan Instance)
		go fetchMonthlyInstancesList(ctx, creds, instances, region, instanceChan, startDate, endDate)
		instanceChans = append(instanceChans, instanceChan)
	}
	for instance := range merge(instanceChans...) {
		report.Instances = append(report.Instances, instance)
	}
	return report, nil
}

// filterInstancesCosts filters instances, cloudwatch and volumes of EC2 instances costs
func filterInstancesCosts(ec2Cost, cloudwatchCost []utils.CostPerResource) ([]utils.CostPerResource, []utils.CostPerResource, []utils.CostPerResource) {
	newInstance := make([]utils.CostPerResource, 0)
	newVolume := make([]utils.CostPerResource, 0)
	newCloudWatch := make([]utils.CostPerResource, 0)
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

func addCostToInstances(report Report, costVolume, costCloudWatch []utils.CostPerResource) Report {
	for i, instance := range report.Instances {
		for volume := range instance.IORead {
			for _, costPerVolume := range costVolume {
				if volume == costPerVolume.Resource {
					report.Instances[i].Cost += costPerVolume.Cost
					report.Instances[i].CostDetail[volume] += costPerVolume.Cost
				}
			}
		}
		for volume := range instance.IOWrite {
			for _, costPerVolume := range costVolume {
				if volume == costPerVolume.Resource {
					report.Instances[i].Cost += costPerVolume.Cost
					report.Instances[i].CostDetail[volume] += costPerVolume.Cost
				}
			}
		}
		for _, cloudWatch := range costCloudWatch {
			if strings.Contains(cloudWatch.Resource, instance.Id) {
				report.Instances[i].Cost += cloudWatch.Cost
				report.Instances[i].CostDetail["cloudWatch"] += cloudWatch.Cost
			}
		}
	}
	return report
}

// PutEc2MonthlyReport puts a monthly report of EC2 instance in ES
func PutEc2MonthlyReport(ctx context.Context, ec2Cost, cloudWatchCost []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting EC2 monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance, costVolume, costCloudWatch := filterInstancesCosts(ec2Cost, cloudWatchCost)
	if len(costInstance) == 0 {
		logger.Info("No EC2 instances found in billing data.", nil)
		return nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixEC2Report)
	if err != nil {
		return err
	} else if already {
		logger.Info("There is already an EC2 monthly report", nil)
		return nil
	}
	report, err := fetchMonthlyInstancesStats(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return err
	}
	report = addCostToInstances(report, costVolume, costCloudWatch)
	return importReportToEs(ctx, aa, report)
}
