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

package ec2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/ec2Reports"
)

// fetchDailyInstancesList sends in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances and getInstanceStats.
func fetchDailyInstancesList(ctx context.Context, creds *credentials.Credentials, region string, instanceChan chan ec2Reports.Instance) error {
	defer close(instanceChan)
	start, end := utils.GetCurrentCheckedDay()
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
			stats := getInstanceStats(ctx, instance, sess, start, end)
			costs := make(map[string]float64, 0)
			instanceChan <- ec2Reports.Instance{
				InstanceBase: ec2Reports.InstanceBase{
					Id:         aws.StringValue(instance.InstanceId),
					Region:     aws.StringValue(instance.Placement.AvailabilityZone),
					State:      aws.StringValue(instance.State.Name),
					Purchasing: getPurchasingOption(instance),
					KeyPair:    aws.StringValue(instance.KeyName),
					Type:       aws.StringValue(instance.InstanceType),
					Platform:   getPlatformName(aws.StringValue(instance.Platform)),
				},
				Tags:  getInstanceTag(instance.Tags),
				Costs: costs,
				Stats: stats,
			}
		}
	}
	return nil
}

// FetchDailyInstancesStats fetches the stats of the EC2 instances of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchInstancesStats should be called every hour.
func FetchDailyInstancesStats(ctx context.Context, awsAccount taws.AwsAccount) error {
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
	now := time.Now().UTC()
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
	instanceChans := make([]<-chan ec2Reports.Instance, 0, len(regions))
	for _, region := range regions {
		instanceChan := make(chan ec2Reports.Instance)
		go fetchDailyInstancesList(ctx, creds, region, instanceChan)
		instanceChans = append(instanceChans, instanceChan)
	}
	instances := make([]ec2Reports.InstanceReport, 0)
	for instance := range merge(instanceChans...) {
		instances = append(instances, ec2Reports.InstanceReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Instance: instance,
		})
	}
	return importInstancesToEs(ctx, awsAccount, instances)
}
