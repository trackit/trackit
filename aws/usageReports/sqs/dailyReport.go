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

package sqs

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchDailySQSList sends in sqsChan the SQS fetched from ListQueues
func fetchDailySQSList(ctx context.Context, creds *credentials.Credentials, region string, sqsChan chan Queue) error {
	defer close(sqsChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := sqs.New(sess)
	queues, err := svc.ListQueues(nil)
	if err != nil {
		logger.Error("Error when describing SQS Queues", err.Error())
		return err
	}
	for _, queueUrl := range queues.QueueUrls {
		ss := strings.Split(aws.StringValue(queueUrl), "/")
		queueName := ss[len(ss)-1]
		sqsChan <- Queue{
			QueueBase: QueueBase{
				Name:   queueName,
				Url:    aws.StringValue(queueUrl),
				Region: region,
			},
			Tags: getSQSTags(ctx, queueUrl, svc),
		}
	}
	return nil
}

// FetchDailySQSStats fetches the stats of the SQS Queue of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchDailySQSStats should be called every hour.
func FetchDailySQSStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching SQS Queue stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorSQSStsSessionName)
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
	sqsChans := make([]<-chan Queue, 0, len(regions))
	for _, region := range regions {
		sqsChan := make(chan Queue)
		go fetchDailySQSList(ctx, creds, region, sqsChan)
		sqsChans = append(sqsChans, sqsChan)
	}
	queues := make([]QueueReport, 0)
	for queue := range merge(sqsChans...) {
		queues = append(queues, QueueReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Queue: queue,
		})
	}
	return importSQSToEs(ctx, awsAccount, queues)
}
