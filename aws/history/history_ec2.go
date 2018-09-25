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

package history

import (
	"time"
	"strings"
	"context"
	"crypto/md5"
	"encoding/json"
	"encoding/base64"

	"github.com/trackit/jsonlog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/config"
	taws "github.com/trackit/trackit-server/aws"
	tec2 "github.com/trackit/trackit-server/aws/ec2"
)

// getInstanceCPUStats gets the CPU average and the CPU peak from CloudWatch
func getInstanceCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, startDate, endDate time.Time) (float64, float64, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("CPUUtilization"),
		StartTime:  aws.Time(startDate),
		EndTime:    aws.Time(endDate),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Average"), aws.String("Maximum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	} else if len(stats.Datapoints) > 0 {
		return aws.Float64Value(stats.Datapoints[0].Average), aws.Float64Value(stats.Datapoints[0].Maximum), nil
	} else {
		return 0, 0, nil
	}
}

// getInstanceNetworkStats gets the network in and out stats from CloudWatch
func getInstanceNetworkStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, startDate, endDate time.Time) (float64, float64, error) {
	statsIn, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("NetworkIn"),
		StartTime:  aws.Time(startDate),
		EndTime:    aws.Time(endDate),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	}
	statsOut, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("NetworkOut"),
		StartTime:  aws.Time(startDate),
		EndTime:    aws.Time(endDate),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	} else if len(statsIn.Datapoints) > 0 && len(statsOut.Datapoints) > 0 {
		return aws.Float64Value(statsIn.Datapoints[0].Sum), aws.Float64Value(statsOut.Datapoints[0].Sum), nil
	} else {
		return 0, 0, nil
	}
}

// getInstanceInternalIOStats gets the IO read and write stats from CloudWatch
func getInstanceInternalIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, startDate, endDate time.Time) (float64, float64, error) {
	statsRead, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("DiskReadBytes"),
		StartTime:  aws.Time(startDate),
		EndTime:    aws.Time(endDate),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	}
	statsWrite, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("DiskWriteBytes"),
		StartTime:  aws.Time(startDate),
		EndTime:    aws.Time(endDate),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	} else if len(statsRead.Datapoints) > 0 && len(statsWrite.Datapoints) > 0 {
		return aws.Float64Value(statsRead.Datapoints[0].Sum), aws.Float64Value(statsWrite.Datapoints[0].Sum), nil
	} else {
		return 0, 0, nil
	}
}

// getInstanceELBIOStats gets the IO read and write stats from CloudWatch
func getInstanceELBIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, startDate, endDate time.Time) (float64, float64, error) {
	statsRead, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EBS"),
		MetricName: aws.String("VolumeReadBytes"),
		StartTime:  aws.Time(startDate),
		EndTime:    aws.Time(endDate),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	}
	statsWrite, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EBS"),
		MetricName: aws.String("VolumeWriteBytes"),
		StartTime:  aws.Time(startDate),
		EndTime:    aws.Time(endDate),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	} else if len(statsRead.Datapoints) > 0 && len(statsWrite.Datapoints) > 0 {
		return aws.Float64Value(statsRead.Datapoints[0].Sum), aws.Float64Value(statsWrite.Datapoints[0].Sum), nil
	} else {
		return 0, 0, nil
	}
}

// getInstanceIOStats gets the IO read and write stats from CloudWatch
func getInstanceIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, volumes []string, startDate, endDate time.Time) (map[tec2.VolumeName]tec2.VolumeValue,
	map[tec2.VolumeName]tec2.VolumeValue, error) {
	statsRead := make(map[tec2.VolumeName]tec2.VolumeValue, 0)
	statsWrite := make(map[tec2.VolumeName]tec2.VolumeValue, 0)
	internalRead, internalWrite, err := getInstanceInternalIOStats(svc, dimensions, startDate, endDate)
	if err != nil {
		return nil, nil, err
	}
	statsRead["internal"] = tec2.VolumeValue(internalRead)
	statsWrite["internal"] = tec2.VolumeValue(internalWrite)
	for _, volume := range volumes {
		dimensionsEBS := []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("VolumeId"),
				Value: aws.String(volume),
			},
		}
		read, write, err := getInstanceELBIOStats(svc, dimensionsEBS, startDate, endDate)
		if err != nil {
			return nil, nil, err
		}
		statsRead[tec2.VolumeName(volume)] = tec2.VolumeValue(read)
		statsWrite[tec2.VolumeName(volume)] = tec2.VolumeValue(write)
	}
	return statsRead, statsWrite, nil
}

// getInstanceStats gets the instance stats from CloudWatch
func getInstanceStats(ctx context.Context, instance *ec2.Instance, sess *session.Session, startDate, endDate time.Time) (tec2.InstanceStats) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{
		&cloudwatch.Dimension{
			Name:  aws.String("InstanceId"),
			Value: aws.String(aws.StringValue(instance.InstanceId)),
		},
	}
	CpuAverage, CpuPeak, err := getInstanceCPUStats(svc, dimensions, startDate, endDate)
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
		return tec2.InstanceStats{}
	}
	NetworkIn, NetworkOut, err := getInstanceNetworkStats(svc, dimensions, startDate, endDate)
	if err != nil {
		logger.Error("Error when fetching Network stats from CloudWatch", err.Error())
		return tec2.InstanceStats{}
	}
	volumes := make([]string, 0)
	for _, volume := range instance.BlockDeviceMappings {
		volumes = append(volumes, aws.StringValue(volume.Ebs.VolumeId))
	}
	IORead, IOWrite, err := getInstanceIOStats(svc, dimensions, volumes, startDate, endDate)
	if err != nil {
		logger.Error("Error when fetching IO stats from CloudWatch", err.Error())
		return tec2.InstanceStats{}
	}
	return tec2.InstanceStats{
		CpuAverage,
		CpuPeak,
		NetworkIn,
		NetworkOut,
		IORead,
		IOWrite,
	}
}

// fetchEc2InstancesList sent in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances, getAccountID and getInstanceStats.
func fetchEc2InstancesList(ctx context.Context, creds *credentials.Credentials, instList []CostPerInstance,
	region string, instanceInfoChan chan tec2.InstanceInfo, startDate, endDate time.Time) error {
	defer close(instanceInfoChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	for _, inst := range instList {
		desc := ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice([]string{inst.Instance})}
		instances, err := svc.DescribeInstances(&desc)
		if err != nil {
			continue
		}
		for _, reservation := range instances.Reservations {
			for _, instance := range reservation.Instances {
				stats := getInstanceStats(ctx, instance, sess, startDate, endDate)
				instanceInfoChan <- tec2.InstanceInfo{
					Id:         aws.StringValue(instance.InstanceId),
					Region:     aws.StringValue(instance.Placement.AvailabilityZone),
					KeyPair:    aws.StringValue(instance.KeyName),
					Tags:       tec2.GetInstanceTag(instance.Tags),
					Type:       aws.StringValue(instance.InstanceType),
					State:      aws.StringValue(instance.State.Name),
					Purchasing: tec2.GetPurchasingOption(instance),
					CpuAverage: stats.CpuAverage,
					CpuPeak:    stats.CpuPeak,
					NetworkIn:  stats.NetworkIn,
					NetworkOut: stats.NetworkOut,
					IORead:     stats.IORead,
					IOWrite:    stats.IOWrite,
					Cost:       inst.Cost,
				}
			}
		}
	}
	return nil
}

// getEc2Metrics get credentials, accounts and region to fetch EC2 instances stats
func getEc2Metrics(ctx context.Context, instances []CostPerInstance, aa taws.AwsAccount, startDate, endDate time.Time) (tec2.ReportInfo, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, tec2.MonitorInstanceStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return tec2.ReportInfo{}, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := tec2.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return tec2.ReportInfo{}, err
	}
	report := tec2.ReportInfo{
		Account:    account,
		ReportDate: startDate,
		ReportType: "monthly",
		Instances:  make([]tec2.InstanceInfo, 0),
	}
	regions, err := tec2.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return tec2.ReportInfo{}, err
	}
	instanceInfoChans := make([]<-chan tec2.InstanceInfo, 0, len(regions))
	for _, region := range regions {
		instanceInfoChan := make(chan tec2.InstanceInfo)
		go fetchEc2InstancesList(ctx, creds, instances, region, instanceInfoChan, startDate, endDate)
		instanceInfoChans = append(instanceInfoChans, instanceInfoChan)
	}
	for instance := range tec2.Merge(instanceInfoChans...) {
		report.Instances = append(report.Instances, instance)
	}
	return report, nil
}

// putEc2ReportInEs puts the history EC2 report in elasticsearch
func putEc2ReportInEs(ctx context.Context, report tec2.ReportInfo, aa taws.AwsAccount) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating EC2 history instances for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	client := es.Client
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
	}{
		report.Account,
		report.ReportDate,
	})
	if err != nil {
		logger.Error("Error when marshaling instance var", err.Error())
		return err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	index := es.IndexNameForUserId(aa.UserId, tec2.IndexPrefixEC2Report)
	if res, err := client.
		Index().
		Index(index).
		Type(tec2.TypeEC2Report).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting InstanceInfo history in ES", err.Error())
	} else {
		logger.Info("Instance put in ES", *res)
	}
	return nil
}

// filterEc2Instances filter instances and volumes of EC2 instances
func filterEc2Instances(ec2Cost, cloudwatchCost []CostPerInstance) ([]CostPerInstance, []CostPerInstance) {
	newInstance := []CostPerInstance{}
	newVolume   := []CostPerInstance{}
	for _, instance := range ec2Cost {
		if len(instance.Instance) == 19 && strings.HasPrefix(instance.Instance, "i-") {
			newInstance = append(newInstance, instance)
		}
		if len(instance.Instance) == 21 && strings.HasPrefix(instance.Instance, "vol-") {
			newVolume = append(newVolume, instance)
		}
	}
	for _, instance := range cloudwatchCost {
		for i, cost := range newInstance {
			if strings.Contains(instance.Instance, cost.Instance) {
				newInstance[i].Cost += instance.Cost
			}
		}
	}
	return newInstance, newVolume
}

// getEc2HistoryReport puts a monthly report of EC2 instance in ES
func getEc2HistoryReport(ctx context.Context, ec2Cost []CostPerInstance, cloudwatchCost []CostPerInstance, aa taws.AwsAccount, startDate, endDate time.Time) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting EC2 history report for " + string(aa.Id) + " (" + aa.Pretty + ")", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance, costVolume := filterEc2Instances(ec2Cost, cloudwatchCost)
	if len(costInstance) == 0 {
		logger.Info("No EC2 instances found in billing data.", aa)
		return nil
	}
	if already, err := checkAlreadyHistory(ctx, startDate, aa, tec2.IndexPrefixEC2Report); already || err != nil {
		logger.Info("There is already an EC2 history report", aa)
		return err
	}
	report, err := getEc2Metrics(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return err
	}
	for i, instance := range report.Instances {
		for volume := range instance.IORead {
			for _, costPerVolume := range costVolume {
				if string(volume) == costPerVolume.Instance {
					report.Instances[i].Cost += costPerVolume.Cost
				}
			}
		}
		for volume := range instance.IOWrite {
			for _, costPerVolume := range costVolume {
				if string(volume) == costPerVolume.Instance {
					report.Instances[i].Cost += costPerVolume.Cost
				}
			}
		}
	}
	err = putEc2ReportInEs(ctx, report, aa)
	return err
}