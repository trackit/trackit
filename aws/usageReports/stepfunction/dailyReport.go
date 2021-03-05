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

package stepfunction

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchDailyStepFunctionsList sends in stepFunctionInfoChan the stepFunctions fetched from ListStateMachines
// and filled by DescribeStepFunctions and getStepFunctionsStates.
func fetchDailyStepFunctionsList(ctx context.Context, creds *credentials.Credentials, region string, stepFunctionChan chan Step) error {
	defer close(stepFunctionChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := sfn.New(sess)
	stepFunctions, err := svc.ListStateMachines(nil)
	if err != nil {
		logger.Error("Error when describing stepFunctions", err.Error())
		return err
	}
	for _, stepFunction := range stepFunctions.StateMachines {
		stepFunctionChan <- Step{
			StepBase: StepBase{
				Name:         aws.StringValue(stepFunction.Name),
				Arn:          aws.StringValue(stepFunction.StateMachineArn),
				CreationDate: aws.TimeValue(stepFunction.CreationDate),
				Region:       region,
			},
			Tags: getStepFunctionTags(ctx, stepFunction, svc),
		}
	}
	return nil
}

// FetchDailyStepFunctionsStats fetches the stats of the StepFunctions of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchDailyStepFunctionsStats should be called every hour.
func FetchDailyStepFunctionsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching StepFunction stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorStepFunctionStsSessionName)
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
	stepChans := make([]<-chan Step, 0, len(regions))
	for _, region := range regions {
		stepChan := make(chan Step)
		go fetchDailyStepFunctionsList(ctx, creds, region, stepChan)
		stepChans = append(stepChans, stepChan)
	}
	stepFunctions := make([]StepReport, 0)
	for stepFunction := range merge(stepChans...) {
		stepFunctions = append(stepFunctions, StepReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Step: stepFunction,
		})
	}
	return importStepFunctionToEs(ctx, awsAccount, stepFunctions)
}
