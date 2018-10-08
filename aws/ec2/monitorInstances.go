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

package ec2

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/trackit/jsonlog"
	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/es"
)

const MonitorInstanceStsSessionName = "monitor-instance"

type (
	TagName        string
	TagValue       string
	VolumeName     string
	VolumeValue    float64

	// ReportInfo represents the report with all the informations for EC2 instances.
	// It will be imported in ElasticSearch thanks to the struct tags.
	ReportInfo struct {
		Account    string         `json:"account"`
		ReportDate time.Time      `json:"reportDate"`
		ReportType string         `json:"reportType"`
		Instances  []InstanceInfo `json:"instances"`
	}

	// InstanceInfo represents all the informations of an EC2 instance.
	// It will be imported in ElasticSearch thanks to the struct tags.
	InstanceInfo struct {
		Id         string                     `json:"id"`
		Region     string                     `json:"region"`
		State      string                     `json:"state"`
		Purchasing string                     `json:"purchasing"`
		CpuAverage float64                    `json:"cpuAverage"`
		CpuPeak    float64                    `json:"cpuPeak"`
		NetworkIn  float64                    `json:"networkIn"`
		NetworkOut float64                    `json:"networkOut"`
		IORead     map[VolumeName]VolumeValue `json:"ioRead"`
		IOWrite    map[VolumeName]VolumeValue `json:"ioWrite"`
		KeyPair    string                     `json:"keyPair"`
		Type       string                     `json:"type"`
		Tags       map[TagName]TagValue       `json:"tags"`
		Cost       float64                    `json:"cost"`
	}

	InstanceStats struct {
		CpuAverage float64
		CpuPeak    float64
		NetworkIn  float64
		NetworkOut float64
		IORead     map[VolumeName]VolumeValue
		IOWrite    map[VolumeName]VolumeValue
	}
)

// Merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func Merge(cs ...<-chan InstanceInfo) <-chan InstanceInfo {
	var wg sync.WaitGroup
	out := make(chan InstanceInfo)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan InstanceInfo) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// GetAccountId gets the AWS Account ID for the given credentials
func GetAccountId(ctx context.Context, sess *session.Session) (string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := sts.New(sess)
	res, err := svc.GetCallerIdentity(nil)
	if err != nil {
		logger.Error("Error when getting caller identity", err.Error())
		return "", err
	}
	return aws.StringValue(res.Account), nil
}

// GetInstanceTag formats []*ec2.Tag to map[TagName]TagValue
func GetInstanceTag(tags []*ec2.Tag) map[TagName]TagValue {
	res := make(map[TagName]TagValue)
	for _, tag := range tags {
		res[TagName(aws.StringValue(tag.Key))] = TagValue(aws.StringValue(tag.Value))
	}
	return res
}

func getCurrentCheckedDay() (start time.Time, end time.Time) {
	now := time.Now()
	end = time.Date(now.Year(), now.Month(), now.Day()-1, 24, 0, 0, 0, now.Location())
	start = time.Date(now.Year(), now.Month(), now.Day()-31, 0, 0, 0, 0, now.Location())
	return start, end
}

// getInstanceCPUStats gets the CPU average and the CPU peak from CloudWatch
func getInstanceCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension) (float64, float64, error) {
	start, end := getCurrentCheckedDay()
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("CPUUtilization"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
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
func getInstanceNetworkStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension) (float64, float64, error) {
	start, end := getCurrentCheckedDay()
	statsIn, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("NetworkIn"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	}
	statsOut, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("NetworkOut"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
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
func getInstanceInternalIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension) (float64, float64, error) {
	start, end := getCurrentCheckedDay()
	statsRead, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("DiskReadBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	}
	statsWrite, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("DiskWriteBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
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
func getInstanceELBIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension) (float64, float64, error) {
	start, end := getCurrentCheckedDay()
	statsRead, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EBS"),
		MetricName: aws.String("VolumeReadBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, err
	}
	statsWrite, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EBS"),
		MetricName: aws.String("VolumeWriteBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
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
func getInstanceIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, volumes []string) (map[VolumeName]VolumeValue,
	map[VolumeName]VolumeValue, error) {
	statsRead := make(map[VolumeName]VolumeValue, 0)
	statsWrite := make(map[VolumeName]VolumeValue, 0)
	internalRead, internalWrite, err := getInstanceInternalIOStats(svc, dimensions)
	if err != nil {
		return nil, nil, err
	}
	statsRead["internal"] = VolumeValue(internalRead)
	statsWrite["internal"] = VolumeValue(internalWrite)
	for _, volume := range volumes {
		dimensionsEBS := []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("VolumeId"),
				Value: aws.String(volume),
			},
		}
		read, write, err := getInstanceELBIOStats(svc, dimensionsEBS)
		if err != nil {
			return nil, nil, err
		}
		statsRead[VolumeName(volume)] = VolumeValue(read)
		statsWrite[VolumeName(volume)] = VolumeValue(write)
	}
	return statsRead, statsWrite, nil
}

// getInstanceStats gets the instance stats from CloudWatch
func getInstanceStats(ctx context.Context, instance *ec2.Instance, sess *session.Session) (InstanceStats) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{
		&cloudwatch.Dimension{
			Name:  aws.String("InstanceId"),
			Value: aws.String(aws.StringValue(instance.InstanceId)),
		},
	}
	CpuAverage, CpuPeak, err := getInstanceCPUStats(svc, dimensions)
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
		return InstanceStats{}
	}
	NetworkIn, NetworkOut, err := getInstanceNetworkStats(svc, dimensions)
	if err != nil {
		logger.Error("Error when fetching Network stats from CloudWatch", err.Error())
		return InstanceStats{}
	}
	volumes := make([]string, 0)
	for _, volume := range instance.BlockDeviceMappings {
		volumes = append(volumes, aws.StringValue(volume.Ebs.VolumeId))
	}
	IORead, IOWrite, err := getInstanceIOStats(svc, dimensions, volumes)
	if err != nil {
		logger.Error("Error when fetching IO stats from CloudWatch", err.Error())
		return InstanceStats{}
	}
	return InstanceStats{
		CpuAverage,
		CpuPeak,
		NetworkIn,
		NetworkOut,
		IORead,
		IOWrite,
	}
}

// GetPurchasingOption returns a string that describes how the instance given as parameter have been purchased
func GetPurchasingOption(instance *ec2.Instance) (string) {
	var purchasing string
	lifeCycle := aws.StringValue(instance.InstanceLifecycle)
	tenancy   := aws.StringValue(instance.Placement.Tenancy)
	if tenancy == "" || tenancy == "default" {
		if lifeCycle == "" {
			purchasing = "on demand"
		} else {
			purchasing = lifeCycle
		}
	} else {
		purchasing = tenancy
	}
	return purchasing
}

// fetchInstancesList sent in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances, getAccountID and getInstanceStats.
func fetchInstancesList(ctx context.Context, creds *credentials.Credentials,
	region string, instanceInfoChan chan InstanceInfo) error {
	defer close(instanceInfoChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	instances, err := svc.DescribeInstances(nil)
	if err != nil {
		logger.Error("Error when describing instances", err.Error())
		return err
	}
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			stats := getInstanceStats(ctx, instance, sess)
			instanceInfoChan <- InstanceInfo{
				Id:         aws.StringValue(instance.InstanceId),
				Region:     aws.StringValue(instance.Placement.AvailabilityZone),
				KeyPair:    aws.StringValue(instance.KeyName),
				Tags:       GetInstanceTag(instance.Tags),
				Type:       aws.StringValue(instance.InstanceType),
				State:      aws.StringValue(instance.State.Name),
				Purchasing: GetPurchasingOption(instance),
				CpuAverage: stats.CpuAverage,
				CpuPeak:    stats.CpuPeak,
				NetworkIn:  stats.NetworkIn,
				NetworkOut: stats.NetworkOut,
				IORead:     stats.IORead,
				IOWrite:    stats.IOWrite,
				Cost:       0,
			}
		}
	}
	return nil
}

// FetchRegionsList fetchs the regions list from AWS and returns an array of their name.
func FetchRegionsList(ctx context.Context, sess *session.Session) ([]string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := ec2.New(sess)
	regions, err := svc.DescribeRegions(nil)
	if err != nil {
		logger.Error("Error when describing regions", err.Error())
		return []string{}, err
	}
	res := make([]string, 0)
	for _, region := range regions.Regions {
		res = append(res, aws.StringValue(region.RegionName))
	}
	return res, nil
}

// importInstancesToEs imports an array of InstanceInfo in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importInstancesToEs(ctx context.Context, aa taws.AwsAccount, report ReportInfo) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating EC2 instances for AWS account.", map[string]interface{}{
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
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixEC2Report)
	if res, err := client.
		Index().
		Index(index).
		Type(TypeEC2Report).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting InstanceInfo in ES", err.Error())
		return err
	} else {
		logger.Info("Instance put in ES", *res)
		return nil
	}
}

// FetchInstancesStats fetchs the stats of the EC2 instances of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchInstancesStats should be called every hour.
func FetchInstancesStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching EC2 instance stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorInstanceStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return err
	}
	report := ReportInfo{
		account,
		time.Now().UTC(),
		"daily",
		make([]InstanceInfo, 0),
	}
	regions, err := FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	instanceInfoChans := make([]<-chan InstanceInfo, 0, len(regions))
	for _, region := range regions {
		instanceInfoChan := make(chan InstanceInfo)
		go fetchInstancesList(ctx, creds, region, instanceInfoChan)
		instanceInfoChans = append(instanceInfoChans, instanceInfoChan)
	}
	for instance := range Merge(instanceInfoChans...) {
		report.Instances = append(report.Instances, instance)
	}
	return importInstancesToEs(ctx, awsAccount, report)
}
