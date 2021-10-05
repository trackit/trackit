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

package cloudformation

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es/indexes/common"
)

// fetchDailyCloudFormationList sends in stackInfoChan the stacks fetched from ListStacks (get only Stacks with a Create Complete Status)
func fetchDailyCloudFormationList(ctx context.Context, creds *credentials.Credentials, region string, stackInfoChan chan Stack) error {
	defer close(stackInfoChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := cloudformation.New(sess)
	stacks, err := svc.ListStacks(&cloudformation.ListStacksInput{
		StackStatusFilter: aws.StringSlice([]string{cloudformation.StackStatusCreateComplete}),
	})
	if err != nil {
		logger.Error("Error when describing CloudFormation stacks", err.Error())
		return err
	}
	for _, stack := range stacks.StackSummaries {
		stackInfoChan <- Stack{
			StackBase: StackBase{
				Name:         aws.StringValue(stack.StackName),
				Id:           aws.StringValue(stack.StackId),
				CreationDate: aws.TimeValue(stack.CreationTime),
				Region:       region,
			},
			Tags: getCloudFormationTags(ctx, stack, svc),
		}
	}
	return nil
}

// FetchDailyCloudFormationStats fetches the stats of the CloudFormation Stacks of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchDailyCloudFormationStats should be called every hour.
func FetchDailyCloudFormationStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching CloudFormation stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorCloudFormationStsSessionName)
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
	stackChans := make([]<-chan Stack, 0, len(regions))
	for _, region := range regions {
		stackChan := make(chan Stack)
		go fetchDailyCloudFormationList(ctx, creds, region, stackChan)
		stackChans = append(stackChans, stackChan)
	}
	stacks := make([]StackReport, 0)
	for stack := range merge(stackChans...) {
		stacks = append(stacks, StackReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Stack: stack,
		})
	}
	return importStacksToEs(ctx, awsAccount, stacks)
}
