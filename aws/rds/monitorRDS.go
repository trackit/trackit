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

package rds

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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/es"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const (
	RDSStsSessionName = "fetch-rds"
)

type (
	RDSInstance struct {
		DBInstanceIdentifier string  `json:"dbInstanceIdentifier"`
		DBInstanceClass      string  `json:"dbInstanceClass"`
		AllocatedStorage     int64   `json:"allocatedStorage"`
		Engine               string  `json:"engine"`
		AvailabilityZone     string  `json:"availabilityZone"`
		MultiAZ              bool    `json:"multiAZ"`
		Cost                 float64 `json:"cost"`
		CpuAverage           float64 `json:"cpuAverage"`
		CpuPeak              float64 `json:"cpuPeak"`
		FreeSpaceMin         float64 `json:"freeSpaceMinimum"`
		FreeSpaceMax         float64 `json:"freeSpaceMaximum"`
		FreeSpaceAve         float64 `json:"freeSpaceAverage"`
	}

	RDSReport struct {
		Account    string        `json:"account"`
		ReportDate time.Time     `json:"reportDate"`
		ReportType string        `json:"reportType"`
		Instances  []RDSInstance `json:"instances"`
	}

	InstanceStats struct {
		CpuAverage   float64
		CpuPeak      float64
		FreeSpaceMin float64
		FreeSpaceMax float64
		FreeSpaceAve float64
	}
)

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func Merge(cs ...<-chan RDSInstance) <-chan RDSInstance {
	var wg sync.WaitGroup
	out := make(chan RDSInstance)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan RDSInstance) {
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

func getCurrentCheckedDay() (start time.Time, end time.Time) {
	now := time.Now()
	end = time.Date(now.Year(), now.Month(), now.Day()-1, 24, 0, 0, 0, now.Location())
	start = time.Date(now.Year(), now.Month(), now.Day()-31, 0, 0, 0, 0, now.Location())
	return start, end
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

// ingestRDSReport saves a report into elasticsearch
func ingestRDSReport(ctx context.Context, aa taws.AwsAccount, report RDSReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating RDS report for AWS account.", map[string]interface{}{
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
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixRDSReport)
	if res, err := client.
		Index().
		Index(index).
		Type(TypeRDSReport).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting RDSReport in ES", err.Error())
	} else {
		logger.Info("RDSReport put in ES", *res)
	}
	return nil
}

// fetchRegionsList fetchs the regions list from AWS and returns an array of their name.
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

// getInstanceCPUStats gets the CPU average and the CPU peak from CloudWatch
func getInstanceCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension) (float64, float64, error) {
	start, end := getCurrentCheckedDay()
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/RDS"),
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

// getInstanceFreeSpaceStats gets the free space stats from CloudWatch
func getInstanceFreeSpaceStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension) (float64, float64, float64, error) {
	start, end := getCurrentCheckedDay()
	freeSpace, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/RDS"),
		MetricName: aws.String("FreeStorageSpace"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30), // Period of thirty days expressed in seconds
		Statistics: []*string{aws.String("Minimum"), aws.String("Maximum"), aws.String("Average")},
		Dimensions: dimensions,
	})
	if err != nil {
		return 0, 0, 0, err
	} else if len(freeSpace.Datapoints) > 0 {
		return aws.Float64Value(freeSpace.Datapoints[0].Minimum),
		aws.Float64Value(freeSpace.Datapoints[0].Maximum),
		aws.Float64Value(freeSpace.Datapoints[0].Average), nil
	} else {
		return 0, 0, 0, nil
	}
}

// getInstanceStats gets the instance stats from CloudWatch
func getInstanceStats(ctx context.Context, instance *rds.DBInstance, sess *session.Session) (InstanceStats) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{
		&cloudwatch.Dimension{
			Name:  aws.String("DBInstanceIdentifier"),
			Value: aws.String(aws.StringValue(instance.DBInstanceIdentifier)),
		},
	}
	CpuAverage, CpuPeak, err := getInstanceCPUStats(svc, dimensions)
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
		return InstanceStats{}
	}
	freeSpaceMin, freeSpaceMax, freeSpaceAve, err := getInstanceFreeSpaceStats(svc, dimensions)
	if err != nil {
		logger.Error("Error when fetching IO stats from CloudWatch", err.Error())
		return InstanceStats{}
	}
	return InstanceStats{
		CpuAverage,
		CpuPeak,
		freeSpaceMin,
		freeSpaceMax,
		freeSpaceAve,
	}
}

// fetchRDSInstancesList fetches the list of instances for a specific region
func fetchRDSInstancesList(ctx context.Context, creds *credentials.Credentials, region string, RDSInstanceChan chan RDSInstance) error {
	defer close(RDSInstanceChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := rds.New(sess)
	desc := rds.DescribeDBInstancesInput{}
	instances, err := svc.DescribeDBInstances(&desc)
	if err != nil {
		logger.Error("Error when getting DB instances pages", err.Error())
		return err
	}
	for _, DBInstance := range instances.DBInstances {
		stats := getInstanceStats(ctx, DBInstance, sess)
		RDSInstanceChan <- RDSInstance{
			DBInstanceIdentifier: aws.StringValue(DBInstance.DBInstanceIdentifier),
			DBInstanceClass:      aws.StringValue(DBInstance.DBInstanceClass),
			AllocatedStorage:     aws.Int64Value(DBInstance.AllocatedStorage),
			Engine:               aws.StringValue(DBInstance.Engine),
			AvailabilityZone:     aws.StringValue(DBInstance.AvailabilityZone),
			MultiAZ:              aws.BoolValue(DBInstance.MultiAZ),
			Cost:                 0,
			CpuAverage:           stats.CpuAverage,
			CpuPeak:              stats.CpuPeak,
			FreeSpaceMin:         stats.FreeSpaceMin,
			FreeSpaceMax:         stats.FreeSpaceMax,
			FreeSpaceAve:         stats.FreeSpaceAve,
		}
	}
	return nil
}

// FetchRDSInfos retrieves RDS informations from the AWS API and generates a report
func FetchRDSInfos(ctx context.Context, aa taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	instances := []RDSInstance{}
	creds, err := taws.GetTemporaryCredentials(aa, RDSStsSessionName)
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
	regions, err := FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	report := RDSReport{
		Account:    account,
		ReportDate: time.Now().UTC(),
		ReportType: "daily",
		Instances:  instances,
	}
	RDSInstanceChans := make([]<-chan RDSInstance, 0, len(regions))
	for _, region := range regions {
		RDSInstanceChan := make(chan RDSInstance)
		go fetchRDSInstancesList(ctx, creds, region, RDSInstanceChan)
		RDSInstanceChans = append(RDSInstanceChans, RDSInstanceChan)
	}
	for instance := range Merge(RDSInstanceChans...) {
		report.Instances = append(report.Instances, instance)
	}
	return ingestRDSReport(ctx, aa, report)
}
