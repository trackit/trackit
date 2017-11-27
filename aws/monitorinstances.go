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

package aws

import (
	"context"
	"reflect"
	"strings"
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

	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/es"
)

const (
	defaultRegion      = "us-west-2"
	elasticsearchIndex = "awshourlyinstanceusage"
	elasticsearchType  = "ec2instance"
)

var (
	// tagsAvailable is used to reflect the structs
	// InstanceTag and InstanceInfo.
	tagsAvailable = []string{"type", "index"}
)

type (
	// InstanceTag represents the tag fields that will
	// be imported in ElasticSearch.
	InstanceTag struct {
		Key   string `json:"key" type:"string" index:"not_analyzed"`
		Value string `json:"value" type:"string" index:"not_analyzed"`
	}

	// InstanceInfo represents all the informations of an EC2 instance.
	// It will be imported in ElasticSearch thanks to the struct tags.
	InstanceInfo struct {
		Account    string        `json:"account" type:"string" index:"not_analyzed"`
		Service    string        `json:"service" type:"string" index:"not_analyzed"`
		ID         string        `json:"id" type:"string" index:"not_analyzed"`
		Region     string        `json:"region" type:"string" index:"not_analyzed"`
		Hour       time.Time     `json:"hour" type:"date" format:"strict_date_optional_time||epoch_millis"`
		CPUAverage float64       `json:"cpuaverage" type:"double"`
		CPUPeak    float64       `json:"cpupeak" type:"double"`
		KeyPair    string        `json:"keypair" type:"string" index:"not_analyzed"`
		Type       string        `json:"type" type:"string" index:"not_analyzed"`
		Tags       []InstanceTag `json:"tags" type:"nested"`
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

// getAccountID allows to get the AWS Account ID (12 digits) of the
// credentials given.
func getAccountID(creds *credentials.Credentials, logger jsonlog.Logger) string {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(defaultRegion),
	}))
	svc := sts.New(sess)
	res, err := svc.GetCallerIdentity(nil)
	if err != nil {
		logger.Error("Get caller identity failed", err)
		return ""
	}
	return aws.StringValue(res.Account)
}

// getInstanceTag formats an array of *ec2.Tag to an array
// of InstanceTag.
func getInstanceTag(tags []*ec2.Tag) []InstanceTag {
	res := make([]InstanceTag, 0, 5)
	for _, tag := range tags {
		res = append(res, InstanceTag{aws.StringValue(tag.Key), aws.StringValue(tag.Value)})
	}
	return res
}

// getInstanceStats gets the CPU average and the CPU peak from CloudWatch
func getInstanceStats(instanceId string, sess *session.Session, region string, logger jsonlog.Logger) (float64, float64) {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, 0, 0, 0, now.Location())
	start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-2, 0, 0, 0, now.Location())
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
		logger.Error("Fetch CPU stats from CloudWatch failed", err)
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
func fetchInstancesList(creds *credentials.Credentials, region string, account string,
	logger jsonlog.Logger, instanceInfoChan chan InstanceInfo) {
	defer close(instanceInfoChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	instances, err := svc.DescribeInstances(nil)
	if err != nil {
		logger.Error("Describe instances failed", err)
		return
	}
	now := time.Now()
	hour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, 0, 0, 0, now.Location())
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			CPUAverage, CPUPeak := getInstanceStats(aws.StringValue(instance.InstanceId), sess, region, logger)
			instanceInfoChan <- InstanceInfo{
				Account:    account,
				Service:    "EC2",
				ID:         aws.StringValue(instance.InstanceId),
				Region:     aws.StringValue(instance.Placement.AvailabilityZone),
				Hour:       hour,
				KeyPair:    aws.StringValue(instance.KeyName),
				Tags:       getInstanceTag(instance.Tags),
				Type:       aws.StringValue(instance.InstanceType),
				CPUAverage: CPUAverage,
				CPUPeak:    CPUPeak,
			}
		}
	}
}

// fetchRegionsList fetchs the regions list from AWS and returns an array of their name.
func fetchRegionsList(creds *credentials.Credentials, logger jsonlog.Logger) []string {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(defaultRegion),
	}))
	svc := ec2.New(sess)
	regions, err := svc.DescribeRegions(nil)
	if err != nil {
		logger.Error("Describe regions failed", err)
		return []string{}
	}
	res := make([]string, 0, 10)
	for _, region := range regions.Regions {
		res = append(res, aws.StringValue(region.RegionName))
	}
	return res
}

// getPropertiesMapping uses the reflection of the struct InstanceInfo and
// InstanceTag to get the mapping of the ElasticSearch documents.
func getPropertiesMapping(field interface{}) map[string]interface{} {
	properties := make(map[string]interface{})
	rtype := reflect.TypeOf(field)
	for i := 0; i < rtype.NumField(); i++ {
		name := strings.ToLower(rtype.Field(i).Name)
		if value, ok := rtype.Field(i).Tag.Lookup("json"); ok {
			name = value
		}
		properties[name] = make(map[string]interface{})
		for _, tagAvailable := range tagsAvailable {
			if value, ok := rtype.Field(i).Tag.Lookup(tagAvailable); ok {
				properties[name].(map[string]interface{})[tagAvailable] = value
			}
		}
	}
	if _, ok := rtype.FieldByName("Tags"); ok {
		properties["tags"] = map[string]interface{}{
			"properties": getPropertiesMapping(InstanceTag{}),
		}
	}
	return properties
}

// createIndexES creates the ElasticSearch index with the properties
// from getPropertiesMapping.
func createIndexES(client *elastic.Client, logger jsonlog.Logger) {
	properties := getPropertiesMapping(InstanceInfo{})
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			elasticsearchType: map[string]interface{}{
				"properties": properties,
			},
		},
	}
	if _, err := client.CreateIndex(elasticsearchIndex).BodyJson(mapping).Do(context.Background()); err != nil {
		logger.Error("ElasticSearch create index failed", err)
		return
	}
}

// importInstancesToES imports an array of InstanceInfo in ElasticSearch.
// It calls createIndexES if the index doesn't exist.
func importInstancesToES(instances []InstanceInfo, logger jsonlog.Logger) {
	// TODO: Get credentials
	client, err := es.NewSignedElasticClient(config.EsAddress, credentials.NewSharedCredentials("", "default"))
	if err != nil {
		logger.Error("Get ElasticSearch client failed", err)
		return
	}
	if exists, err := client.IndexExists(elasticsearchIndex).Do(context.Background()); err != nil {
		logger.Error("ElasticSearch check if index exists failed", err)
		return
	} else if !exists {
		createIndexES(client, logger)
	}
	for _, instance := range instances {
		if res, err := client.Index().Index(elasticsearchIndex).Type(elasticsearchType).BodyJson(instance).Do(context.Background()); err != nil {
			logger.Error("Instance put in ES failed", err)
		} else {
			logger.Info("Instance put in ES", *res)
		}
	}
}

// FetchInstancesStats fetchs the stats of the EC2 instances of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchInstancesStats should be called every hour.
func FetchInstancesStats(awsAccount AwsAccount) {
	logger := jsonlog.DefaultLogger
	tcreds, err := GetTemporaryCredentials(context.Background(), awsAccount, "fetchInstanceStats")
	creds := credentials.NewStaticCredentials(
		aws.StringValue(tcreds.Credentials.AccessKeyId),
		aws.StringValue(tcreds.Credentials.SecretAccessKey),
		aws.StringValue(tcreds.Credentials.SessionToken),
	)
	account := getAccountID(creds, logger)
	if err != nil {
		logger.Error("Get temporary credentials failed", err)
		return
	}
	regions := fetchRegionsList(creds, logger)
	instanceInfoChans := make([]<-chan InstanceInfo, 0, 15)
	for _, region := range regions {
		instanceInfoChan := make(chan InstanceInfo)
		go fetchInstancesList(creds, region, account, logger, instanceInfoChan)
		instanceInfoChans = append(instanceInfoChans, instanceInfoChan)
	}
	instances := make([]InstanceInfo, 0, 20)
	for instance := range merge(instanceInfoChans...) {
		instances = append(instances, instance)
	}
	importInstancesToES(instances, logger)
}
