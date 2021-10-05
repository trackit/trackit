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

package lambda

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/lambdaReports"
)

// fetchDailyFunctionsList sends in functionInfoChan the functions fetched from DescribeFunctions
// and filled by DescribeFunctions and getFunctionStats.
func fetchDailyFunctionsList(ctx context.Context, creds *credentials.Credentials, region string, functionChan chan lambdaReports.Function) error {
	defer close(functionChan)
	start, end := utils.GetCurrentCheckedDay()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := lambda.New(sess)
	functions, err := svc.ListFunctions(nil)
	if err != nil {
		logger.Error("Error when describing functions", err.Error())
		return err
	}
	for _, function := range functions.Functions {
		stats := getFunctionStats(ctx, function, sess, start, end)
		functionChan <- lambdaReports.Function{
			FunctionBase: lambdaReports.FunctionBase{
				Name:         aws.StringValue(function.FunctionName),
				Description:  aws.StringValue(function.Description),
				Version:      aws.StringValue(function.Version),
				LastModified: aws.StringValue(function.LastModified),
				Runtime:      aws.StringValue(function.Runtime),
				Size:         aws.Int64Value(function.CodeSize),
				Memory:       aws.Int64Value(function.MemorySize),
				Region:       region,
			},
			Tags:  getFunctionTags(ctx, function, svc),
			Stats: stats,
		}
	}
	return nil
}

// FetchDailyFunctionsStats fetches the stats of the Lambda functions of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchFunctionsStats should be called every hour.
func FetchDailyFunctionsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching Lambda function stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorFunctionStsSessionName)
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
	functionChans := make([]<-chan lambdaReports.Function, 0, len(regions))
	for _, region := range regions {
		functionChan := make(chan lambdaReports.Function)
		go fetchDailyFunctionsList(ctx, creds, region, functionChan)
		functionChans = append(functionChans, functionChan)
	}
	functions := make([]lambdaReports.FunctionReport, 0)
	for function := range merge(functionChans...) {
		functions = append(functions, lambdaReports.FunctionReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Function: function,
		})
	}
	return importFunctionsToEs(ctx, awsAccount, functions)
}
