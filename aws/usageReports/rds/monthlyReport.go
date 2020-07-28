//   Copyright 2019 MSolution.IO
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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/rdsReports"
)

// fetchMonthlyInstancesList fetches instances based on billing data
func fetchMonthlyInstancesList(ctx context.Context, creds *credentials.Credentials, inst common.CostPerResource,
	region string, instanceChan chan rdsReports.Instance, startDate, endDate time.Time) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := rds.New(sess)
	desc := rds.DescribeDBInstancesInput{DBInstanceIdentifier: aws.String(inst.Resource)}
	instances, err := svc.DescribeDBInstances(&desc)
	if err != nil {
		return err
	}
	for _, DBInstance := range instances.DBInstances {
		tags := getInstanceTags(ctx, DBInstance, svc)
		stats := getInstanceStats(ctx, DBInstance, sess, startDate, endDate)
		costs := make(map[string]float64, 0)
		costs["instance"] = inst.Cost
		instanceChan <- rdsReports.Instance{
			InstanceBase: rdsReports.InstanceBase{
				DBInstanceIdentifier: aws.StringValue(DBInstance.DBInstanceIdentifier),
				AvailabilityZone:     aws.StringValue(DBInstance.AvailabilityZone),
				DBInstanceClass:      aws.StringValue(DBInstance.DBInstanceClass),
				Engine:               aws.StringValue(DBInstance.Engine),
				AllocatedStorage:     aws.Int64Value(DBInstance.AllocatedStorage),
				MultiAZ:              aws.BoolValue(DBInstance.MultiAZ),
			},
			Tags:  tags,
			Costs: costs,
			Stats: stats,
		}
	}
	return nil
}

// getRdsMetrics gets credentials, accounts and region to fetch RDS instances stats
func getRdsMetrics(ctx context.Context, instancesList []common.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) ([]rdsReports.InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, RDSStsSessionName)
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
	instanceChans := make([]<-chan rdsReports.Instance, 0, len(regions))
	for _, instance := range instancesList {
		for _, region := range regions {
			if strings.Contains(instance.Region, region) {
				instanceChan := make(chan rdsReports.Instance)
				go fetchMonthlyInstancesList(ctx, creds, instance, region, instanceChan, startDate, endDate)
				instanceChans = append(instanceChans, instanceChan)
			}
		}
	}
	instances := make([]rdsReports.InstanceReport, 0)
	for instance := range merge(instanceChans...) {
		instances = append(instances, rdsReports.InstanceReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Instance: instance,
		})
	}
	return instances, nil
}

// filterRdsInstances filters cost per instance to get only costs associated to a RDS instance
func filterRdsInstances(rdsCost []common.CostPerResource) []common.CostPerResource {
	costInstances := []common.CostPerResource{}
	for _, instance := range rdsCost {
		// format in billing data for an RDS instance is: "arn:aws:rds:us-west-2:394125495069:db:instancename"
		// so i get the 7th element of the split by ":"
		split := strings.Split(instance.Resource, ":")
		if len(split) == 7 && split[2] == "rds" {
			costInstances = append(costInstances, common.CostPerResource{split[6], instance.Cost, instance.Region})
		}
	}
	return costInstances
}

// PutRdsMonthlyReport puts a monthly report of RDS in ES
func PutRdsMonthlyReport(ctx context.Context, rdsCost []common.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting RDS monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance := filterRdsInstances(rdsCost)
	if len(costInstance) == 0 {
		logger.Info("No RDS instances found in billing data.", nil)
		return false, nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, rdsReports.IndexSuffix)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an RDS monthly report", nil)
		return false, nil
	}
	instances, err := getRdsMetrics(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importInstancesToEs(ctx, aa, instances)
	if err != nil {
		return false, err
	}
	return true, nil
}
