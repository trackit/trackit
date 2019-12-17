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

package mediaconvert

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchDailyJobsList get the daily jobs list for a region
func fetchDailyJobsList(_ context.Context, creds *credentials.Credentials,
	region string, jobsChan chan Job) error {
	defer close(jobsChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := mediaconvert.New(sess)
	endpoints, err := svc.DescribeEndpoints(nil)
	if err != nil {
		return err
	} else if endpoints == nil {
		return nil
	}
	var nextToken *string
	for _, endpoint := range endpoints.Endpoints {
		subSvc := mediaconvert.New(sess, &aws.Config{Endpoint: endpoint.Url})
		nextToken = nil
		for nextToken, err = getJobsFromAWS(jobsChan, subSvc, region, nextToken); nextToken != nil; {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getJobsFromAWS(jobsChan chan Job, svc *mediaconvert.MediaConvert, region string, token *string) (*string, error) {
	listJob, err := svc.ListJobs(&mediaconvert.ListJobsInput{NextToken: token})
	if err != nil {
		return nil, err
	}
	for _, job := range listJob.Jobs {
		jobsChan <- Job{
			Id:                   aws.StringValue(job.Id),
			Arn:                  aws.StringValue(job.Arn),
			Region:               region,
			BillingTagsSource:    aws.StringValue(job.BillingTagsSource),
			CreatedAt:            aws.TimeValue(job.CreatedAt),
			CurrentPhase:         aws.StringValue(job.CurrentPhase),
			ErrorCode:            aws.Int64Value(job.ErrorCode),
			ErrorMessage:         aws.StringValue(job.ErrorMessage),
			JobPercentComplete:   aws.Int64Value(job.JobPercentComplete),
			JobTemplate:          aws.StringValue(job.JobTemplate),
			OutputGroupDetails:   getOutputGroupDetails(job.OutputGroupDetails),
			Queue:                aws.StringValue(job.Queue),
			RetryCount:           aws.Int64Value(job.RetryCount),
			Role:                 aws.StringValue(job.Role),
			Status:               aws.StringValue(job.Status),
			StatusUpdateInterval: aws.StringValue(job.StatusUpdateInterval),
			Timing: Timing{
				FinishTime: aws.TimeValue(job.Timing.FinishTime),
				StartTime:  aws.TimeValue(job.Timing.StartTime),
				SubmitTime: aws.TimeValue(job.Timing.SubmitTime),
			},
			UserMetadata: getUserMetadata(job.UserMetadata),
			Cost:         0,
		}
	}
	return listJob.NextToken, nil
}

// fetchDailyJobsStats gets credentials, accounts and region to fetch Daily MediaConvert Jobs stats
func fetchDailyJobsStats(ctx context.Context, aa taws.AwsAccount) ([]JobReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorJobStsSessionName)
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
	jobChans := make([]<-chan Job, 0, len(regions))
	for _, region := range regions {
		jobChan := make(chan Job)
		go fetchDailyJobsList(ctx, creds, region, jobChan)
		jobChans = append(jobChans, jobChan)
	}
	now := time.Now().UTC()
	jobsList := make([]JobReport, 0)
	for job := range merge(jobChans...) {
		jobsList = append(jobsList, JobReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Job: job,
		})
	}
	return jobsList, nil
}

// PutMediaConvertDailyReport puts a monthly report of MediaConvert Jobs in ES
func PutMediaConvertDailyReport(ctx context.Context, aa taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting MediaConvert daily report", map[string]interface{}{
		"awsAccountId": aa.Id,
	})
	jobs, err := fetchDailyJobsStats(ctx, aa)
	if err != nil {
		return err
	}
	err = importJobsToEs(ctx, aa, jobs)
	if err != nil {
		return err
	}
	return nil
}
