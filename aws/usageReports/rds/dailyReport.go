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

package rds

import (
	"time"
	"context"

	"github.com/trackit/jsonlog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/trackit/trackit-server/config"
	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
)

// fetchDailyInstancesList fetches the list of instances for a specific region
func fetchDailyInstancesList(ctx context.Context, creds *credentials.Credentials, region string, InstanceChan chan Instance) error {
	defer close(InstanceChan)
	start, end := utils.GetCurrentCheckedDay()
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
		stats := getInstanceStats(ctx, DBInstance, sess, start, end)
		InstanceChan <- Instance{
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
			CostDetail:           make(map[string]float64, 0),
		}
	}
	return nil
}

// FetchDailyInstanceStats retrieves RDS information from the AWS API and generates a report
func FetchDailyInstancesStats(ctx context.Context, aa taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching RDS instance stats", map[string]interface{}{"awsAccountId": aa.Id})
	instances := []Instance{}
	creds, err := taws.GetTemporaryCredentials(aa, RDSStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	report := Report{
		Account:    account,
		ReportDate: time.Now().UTC(),
		ReportType: "daily",
		Instances:  instances,
	}
	InstanceChans := make([]<-chan Instance, 0, len(regions))
	for _, region := range regions {
		InstanceChan := make(chan Instance)
		go fetchDailyInstancesList(ctx, creds, region, InstanceChan)
		InstanceChans = append(InstanceChans, InstanceChan)
	}
	for instance := range merge(InstanceChans...) {
		report.Instances = append(report.Instances, instance)
	}
	return importReportToEs(ctx, aa, report)
}
