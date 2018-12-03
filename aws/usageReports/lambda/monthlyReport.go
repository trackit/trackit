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

package lambda

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/trackit/jsonlog"
	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/config"
)

// fetchMonthlyInstancesList sends in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances and getInstanceStats.
func fetchMonthlyInstancesList(ctx context.Context, creds *credentials.Credentials, inst utils.CostPerResource,
	region string, instanceChan chan Instance) error {
	defer close(instanceChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := lambda.New(sess)
	input := lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(inst.Resource),
	}
	functions, err := svc.ListVersionsByFunction(&input)
	if err != nil {
		logger.Error("Error when describing instances", err.Error())
		return err
	}
	for _, function := range functions.Versions {
		costs := make(map[string]float64, 0)
		costs["function"] = inst.Cost
		instanceChan <- Instance{
			InstanceBase: InstanceBase{
				Name:        aws.StringValue(function.FunctionName),
				Description: aws.StringValue(function.Description),
				Size:        aws.Int64Value(function.CodeSize),
				Memory:      aws.Int64Value(function.MemorySize),
			},
			Tags:  getFunctionTags(ctx, function, svc),
			Costs: costs,
		}
	}
	return nil
}

// fetchMonthlyInstancesStats gets credentials, accounts and region to fetch Lambda instances stats
func fetchMonthlyInstancesStats(ctx context.Context, instances []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) ([]InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorInstanceStsSessionName)
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
	instanceChans := make([]<-chan Instance, 0, len(regions))
	for _, instance := range instances {
		for _, region := range regions {
			if strings.Contains(instance.Region, region) {
				instanceChan := make(chan Instance)
				go fetchMonthlyInstancesList(ctx, creds, instance, region, instanceChan)
				instanceChans = append(instanceChans, instanceChan)
			}
		}
	}
	instancesList := make([]InstanceReport, 0)
	for instance := range merge(instanceChans...) {
		instancesList = append(instancesList, InstanceReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Instance: instance,
		})
	}
	return instancesList, nil
}

// filterLambdaInstances filters cost per instance to get only costs associated to a Lambda instance
func filterLambdaInstances(lambdaCost []utils.CostPerResource) []utils.CostPerResource {
	costInstances := []utils.CostPerResource{}
	for _, instance := range lambdaCost {
		// format in billing data for a Lambda instance is: "arn:aws:lambda:(region):(aws_id):function:(lambda name)"
		// so i get the 7th element of the split by ":"
		split := strings.Split(instance.Resource, ":")
		if len(split) == 7 && split[2] == "lambda" {
			costInstances = append(costInstances, utils.CostPerResource{split[6], instance.Cost, instance.Region})
		}
	}
	return costInstances
}

// PutLambdaMonthlyReport puts a monthly report of Lambda instance in ES
func PutLambdaMonthlyReport(ctx context.Context, costs []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting Lambda monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costs = filterLambdaInstances(costs)
	if len(costs) == 0 {
		logger.Info("No Lambda instances found in billing data.", nil)
		return false, nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixLambdaReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an Lambda monthly report", nil)
		return false, nil
	}
	instances, err := fetchMonthlyInstancesStats(ctx, costs, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importInstancesToEs(ctx, aa, instances)
	if err != nil {
		return false, err
	}
	return true, nil
}
