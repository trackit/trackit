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
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/util"
	"github.com/trackit/trackit-server/config"
	taws "github.com/trackit/trackit-server/aws"
	trds "github.com/trackit/trackit-server/aws/rds"
)

// getRdsInstanceCPUStats gets the CPU average and the CPU peak from CloudWatch
func getRdsInstanceCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/RDS"),
		MetricName: aws.String("CPUUtilization"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
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

// getRdsInstanceFreeSpaceStats gets the free space stats from CloudWatch
func getRdsInstanceFreeSpaceStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, float64, error) {
	freeSpace, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/RDS"),
		MetricName: aws.String("FreeStorageSpace"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
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

// getRdsInstanceStats gets the instance stats from CloudWatch
func getRdsInstanceStats(ctx context.Context, instance *rds.DBInstance, sess *session.Session, start, end time.Time) (trds.InstanceStats) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{
		&cloudwatch.Dimension{
			Name:  aws.String("DBInstanceIdentifier"),
			Value: aws.String(aws.StringValue(instance.DBInstanceIdentifier)),
		},
	}
	CpuAverage, CpuPeak, err := getRdsInstanceCPUStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
		return trds.InstanceStats{}
	}
	freeSpaceMin, freeSpaceMax, freeSpaceAve, err := getRdsInstanceFreeSpaceStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching IO stats from CloudWatch", err.Error())
		return trds.InstanceStats{}
	}
	return trds.InstanceStats{
		CpuAverage,
		CpuPeak,
		freeSpaceMin,
		freeSpaceMax,
		freeSpaceAve,
	}
}

// fetchRDSInstancesList fetches the list of instances for a specific region
func fetchRdsInstancesList(ctx context.Context, creds *credentials.Credentials, instList []CostPerInstance,
	region string, instanceInfoChan chan trds.RDSInstance, startDate, endDate time.Time) error {
	defer close(instanceInfoChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := rds.New(sess)
	for _, inst := range instList {
		desc := rds.DescribeDBInstancesInput{DBInstanceIdentifier: aws.String(inst.Instance)}
		instances, err := svc.DescribeDBInstances(&desc)
		if err != nil {
			continue
		}
		for _, DBInstance := range instances.DBInstances {
			stats := getRdsInstanceStats(ctx, DBInstance, sess, startDate, endDate)
			instanceInfoChan <- trds.RDSInstance{
				DBInstanceIdentifier: util.SafeStringFromPtr(DBInstance.DBInstanceIdentifier),
				DBInstanceClass:      util.SafeStringFromPtr(DBInstance.DBInstanceClass),
				AllocatedStorage:     util.SafeInt64FromPtr(DBInstance.AllocatedStorage),
				Engine:               util.SafeStringFromPtr(DBInstance.Engine),
				AvailabilityZone:     util.SafeStringFromPtr(DBInstance.AvailabilityZone),
				MultiAZ:              util.SafeBoolFromPtr(DBInstance.MultiAZ),
				Cost:                 inst.Cost,
				CpuAverage:           stats.CpuAverage,
				CpuPeak:              stats.CpuPeak,
				FreeSpaceMin:         stats.FreeSpaceMin,
				FreeSpaceMax:         stats.FreeSpaceMax,
				FreeSpaceAve:         stats.FreeSpaceAve,
			}
		}
	}
	return nil
}

// getRdsMetrics gets credentials, accounts and region to fetch RDS instances stats
func getRdsMetrics(ctx context.Context, instances []CostPerInstance, aa taws.AwsAccount, startDate, endDate time.Time) (trds.RDSReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, trds.RDSStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return trds.RDSReport{}, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := trds.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return trds.RDSReport{}, err
	}
	report := trds.RDSReport{
		Account:    account,
		ReportDate: startDate,
		ReportType: "monthly",
		Instances:  make([]trds.RDSInstance, 0),
	}
	regions, err := trds.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return trds.RDSReport{}, err
	}
	instanceInfoChans := make([]<-chan trds.RDSInstance, 0, len(regions))
	for _, region := range regions {
		instanceInfoChan := make(chan trds.RDSInstance)
		go fetchRdsInstancesList(ctx, creds, instances, region, instanceInfoChan, startDate, endDate)
		instanceInfoChans = append(instanceInfoChans, instanceInfoChan)
	}
	for instance := range trds.Merge(instanceInfoChans...) {
		report.Instances = append(report.Instances, instance)
	}
	return report, nil
}

// putRdsReportInEs puts the history RDS report in elasticsearch
func putRdsReportInEs(ctx context.Context, report trds.RDSReport, aa taws.AwsAccount) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating RDS history instances for AWS account.", map[string]interface{}{
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
	index := es.IndexNameForUserId(aa.UserId, trds.IndexPrefixRDSReport)
	if res, err := client.
		Index().
		Index(index).
		Type(trds.TypeRDSReport).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting InstanceInfo history in ES", err.Error())
		return err
	} else {
		logger.Info("Instance put in ES", *res)
		return nil
	}
}

// filterRdsInstances filters instances
func filterRdsInstances(rdsCost []CostPerInstance) ([]CostPerInstance) {
	costInstances := []CostPerInstance{}
	for _, instance := range rdsCost {
		split := strings.Split(instance.Instance, ":")
		if len(split) == 7 || split[2] != "rds" {
			costInstances = append(costInstances, CostPerInstance{split[6], instance.Cost})
		}
	}
	return costInstances
}

// getRdsHistoryReport puts a monthly report of RDS in ES
func getRdsHistoryReport(ctx context.Context, rdsCost []CostPerInstance, aa taws.AwsAccount, startDate, endDate time.Time) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting RDS history report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance := filterRdsInstances(rdsCost)
	if len(costInstance) == 0 {
		logger.Info("No RDS instances found in billing data.", nil)
		return nil
	}
	if already, err := checkAlreadyHistory(ctx, startDate, aa, trds.IndexPrefixRDSReport); already || err != nil {
		logger.Info("There is already an RDS history report", err)
		return err
	}
	report, err := getRdsMetrics(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return err
	}
	return putRdsReportInEs(ctx, report, aa)
}
