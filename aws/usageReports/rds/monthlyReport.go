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
	"time"
	"strings"
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

// fetchMonthlyInstancesList fetches instances based on billing data
func fetchMonthlyInstancesList(ctx context.Context, creds *credentials.Credentials, instList []utils.CostPerInstance,
	region string, instanceChan chan Instance, startDate, endDate time.Time) error {
	defer close(instanceChan)
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
			stats := getInstanceStats(ctx, DBInstance, sess, startDate, endDate)
			detail := make(map[string]float64, 0)
			detail["instance"] = inst.Cost
			instanceChan <- Instance{
				DBInstanceIdentifier: aws.StringValue(DBInstance.DBInstanceIdentifier),
				DBInstanceClass:      aws.StringValue(DBInstance.DBInstanceClass),
				AllocatedStorage:     aws.Int64Value(DBInstance.AllocatedStorage),
				Engine:               aws.StringValue(DBInstance.Engine),
				AvailabilityZone:     aws.StringValue(DBInstance.AvailabilityZone),
				MultiAZ:              aws.BoolValue(DBInstance.MultiAZ),
				Cost:                 inst.Cost,
				CpuAverage:           stats.CpuAverage,
				CpuPeak:              stats.CpuPeak,
				FreeSpaceMin:         stats.FreeSpaceMin,
				FreeSpaceMax:         stats.FreeSpaceMax,
				FreeSpaceAve:         stats.FreeSpaceAve,
				CostDetail:           detail,
			}
		}
	}
	return nil
}

// getRdsMetrics gets credentials, accounts and region to fetch RDS instances stats
func getRdsMetrics(ctx context.Context, instances []utils.CostPerInstance, aa taws.AwsAccount, startDate, endDate time.Time) (Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, RDSStsSessionName)
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

// filterRdsInstances filters instances
func filterRdsInstances(rdsCost []utils.CostPerInstance) ([]utils.CostPerInstance) {
	costInstances := []utils.CostPerInstance{}
	for _, instance := range rdsCost {
		split := strings.Split(instance.Instance, ":")
		if len(split) == 7 || split[2] != "rds" {
			costInstances = append(costInstances, utils.CostPerInstance{split[6], instance.Cost})
		}
	}
	return costInstances
}

// GetRdsMonthlyReport puts a monthly report of RDS in ES
func GetRdsMonthlyReport(ctx context.Context, rdsCost []utils.CostPerInstance, aa taws.AwsAccount, startDate, endDate time.Time) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting RDS utils report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance := filterRdsInstances(rdsCost)
	if len(costInstance) == 0 {
		logger.Info("No RDS instances found in billing data.", nil)
		return nil
	}
	if already, err := utils.CheckAlreadyHistory(ctx, startDate, aa, IndexPrefixRDSReport); already || err != nil {
		logger.Info("There is already an RDS monthly report", err)
		return err
	}
	report, err := getRdsMetrics(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return err
	}
	return importReportToEs(ctx, aa, report)
}
