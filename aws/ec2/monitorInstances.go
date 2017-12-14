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
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit2/aws"
	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/es"
)

const MonitorInstanceStsSessionName = "monitor-instance"

type (
	TagName  string
	TagValue string

	// InstanceInfo represents all the informations of an EC2 instance.
	// It will be imported in ElasticSearch thanks to the struct tags.
	InstanceInfo struct {
		Account    string               `json:"account"`
		Service    string               `json:"service"`
		Id         string               `json:"id"`
		Region     string               `json:"region"`
		StartDate  time.Time            `json:"startDate"`
		EndDate    time.Time            `json:"endDate"`
		CpuAverage float64              `json:"cpuAverage"`
		CpuPeak    float64              `json:"cpuPeak"`
		KeyPair    string               `json:"keyPair"`
		Type       string               `json:"type"`
		Tags       map[TagName]TagValue `json:"tags"`
	}
)

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func merge(cs ...<-chan InstanceInfo) <-chan InstanceInfo {
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

// getAccountId gets the AWS Account ID for the given credentials
func getAccountId(ctx context.Context, sess *session.Session) (string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := sts.New(sess)
	res, err := svc.GetCallerIdentity(nil)
	if err != nil {
		logger.Error("Error when getting caller identity", err.Error())
		return "", err
	}
	return aws.StringValue(res.Account), nil
}

// getInstanceTag formats []*ec2.Tag to map[TagName]TagValue
func getInstanceTag(tags []*ec2.Tag) map[TagName]TagValue {
	res := make(map[TagName]TagValue)
	for _, tag := range tags {
		res[TagName(aws.StringValue(tag.Key))] = TagValue(aws.StringValue(tag.Value))
	}
	return res
}

func getCurrentCheckedHour() (start time.Time, end time.Time) {
	now := time.Now()
	end = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, 0, 0, 0, now.Location())
	start = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-2, 0, 0, 0, now.Location())
	return start, end
}

// getInstanceStats gets the CPU average and the CPU peak from CloudWatch
func getInstanceStats(ctx context.Context, instanceId string,
	sess *session.Session, region string) (float64, float64) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	start, end := getCurrentCheckedHour()
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{
		&cloudwatch.Dimension{
			Name:  aws.String("InstanceId"),
			Value: aws.String(instanceId),
		},
	}
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("CPUUtilization"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(86400 * 30)),
		Statistics: []*string{aws.String("Average"), aws.String("Maximum")},
		Dimensions: dimensions,
	})
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
		return 0, 0
	} else if len(stats.Datapoints) > 0 {
		return aws.Float64Value(stats.Datapoints[0].Average), aws.Float64Value(stats.Datapoints[0].Maximum)
	} else {
		return 0, 0
	}
}

// fetchInstancesList, firstly, calls DescribeInstances from AWS to list all the instances in a region.
// Secondly, it generates an InstanceInfo filled by DescribeInstances, getAccountID and getInstanceStats.
// Finally, it sends the InstanceInfo generated in a chan.
func fetchInstancesList(ctx context.Context, creds *credentials.Credentials,
	region string, account string, instanceInfoChan chan InstanceInfo) {
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
		return
	}
	start, end := getCurrentCheckedHour()
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			CpuAverage, CpuPeak := getInstanceStats(ctx, aws.StringValue(instance.InstanceId), sess, region)
			instanceInfoChan <- InstanceInfo{
				Account:    account,
				Service:    "EC2",
				Id:         aws.StringValue(instance.InstanceId),
				Region:     aws.StringValue(instance.Placement.AvailabilityZone),
				StartDate:  start,
				EndDate:    end,
				KeyPair:    aws.StringValue(instance.KeyName),
				Tags:       getInstanceTag(instance.Tags),
				Type:       aws.StringValue(instance.InstanceType),
				CpuAverage: CpuAverage,
				CpuPeak:    CpuPeak,
			}
		}
	}
}

// fetchRegionsList fetchs the regions list from AWS and returns an array of their name.
func fetchRegionsList(ctx context.Context, sess *session.Session) ([]string, error) {
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

// createIndexEs creates the ElasticSearch index with the properties
// from getPropertiesMapping.
func createIndexEs(ctx context.Context, client *elastic.Client) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	mapping := "{\"mappings\":{\"" + DocumentEc2 + "\":" + mappingEc2Instance + "}}"
	if _, err := client.CreateIndex(IndexHourlyInstanceUsage).BodyString(mapping).Do(context.Background()); err != nil {
		logger.Error("Error when creating ElasticSearch index", err.Error())
	}
}

// importInstancesToEs imports an array of InstanceInfo in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importInstancesToEs(ctx context.Context, instances []InstanceInfo) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	client := es.Client
	if exists, err := client.IndexExists(IndexHourlyInstanceUsage).Do(context.Background()); err != nil {
		logger.Error("Error when checking if the ElasticSearch index exists", err.Error())
		return
	} else if !exists {
		createIndexEs(ctx, client)
	}
	for _, instance := range instances {
		ji, err := json.Marshal(instance)
		if err != nil {
			logger.Error("Error when marshaling instance var", err.Error())
			continue
		}
		hash := md5.Sum(ji)
		hash64 := base64.URLEncoding.EncodeToString(hash[:])
		if res, err := client.
			Index().
			Index(IndexHourlyInstanceUsage).
			Type(DocumentEc2).
			BodyJson(instance).
			Id(hash64).
			Do(context.Background()); err != nil {
			logger.Error("Error when putting InstanceInfo in ES", err.Error())
		} else {
			logger.Info("Instance put in ES", *res)
		}
	}
}

// FetchInstancesStats fetchs the stats of the EC2 instances of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchInstancesStats should be called every hour.
func FetchInstancesStats(awsAccount taws.AwsAccount) error {
	ctx := context.Background()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching instance stats for "+string(awsAccount.Id)+" ("+awsAccount.Pretty+")", nil)
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorInstanceStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := getAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return err
	}
	regions, err := fetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	instanceInfoChans := make([]<-chan InstanceInfo, 0, 15)
	for _, region := range regions {
		instanceInfoChan := make(chan InstanceInfo)
		go fetchInstancesList(ctx, creds, region, account, instanceInfoChan)
		instanceInfoChans = append(instanceInfoChans, instanceInfoChan)
	}
	instances := make([]InstanceInfo, 0, 20)
	for instance := range merge(instanceInfoChans...) {
		instances = append(instances, instance)
	}
	importInstancesToEs(ctx, instances)
	return nil
}
