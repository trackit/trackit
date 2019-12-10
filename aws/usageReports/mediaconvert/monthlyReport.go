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
	"encoding/json"
	"fmt"
	"github.com/trackit/trackit/es"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/errors"
)

// getElasticSearchMediaConvertJob prepares and run the request to retrieve the a report of an instance
// It will return the data and an error.
func getElasticSearchMediaConvertJob(ctx context.Context, account, instance string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("instance.id", instance))
	search := client.Search().Index(index).Size(1).Query(query)
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, errors.GetErrorMessage(ctx, err)
	}
	return res, nil
}

// getJobInfoFromEs gets information about an instance from previous report to put it in the new report
func getJobInfoFromES(ctx context.Context, cost JobInformations, account string, userId int) Job {
	var docType JobReport
	var job = Job{
		Id:         "N/A",
		Region:     "N/A",
		Arn:      "N/A",
		Costs: make(map[time.Time]float64, 0),
	}
	res, err := getElasticSearchMediaConvertJob(ctx, account, cost.Arn,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixMediaConvertReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			job.Region = docType.Job.Region
			job.Id = docType.Job.Id
			job.Arn = docType.Job.Arn
			job.Costs = docType.Job.Costs
		}
	}
	return job
}

// fetchMonthlyJobsList sends in instanceInfoChan the instances fetched from DescribeJobs
// and filled by DescribeJobs and getJobStats.
func fetchMonthlyJobsList(ctx context.Context, creds *credentials.Credentials,
	account, region, jobId string, cost JobInformations, instanceChan chan Job, startDate, endDate time.Time, userId int) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := mediaconvert.New(sess)
	endpoints, err := svc.DescribeEndpoints(nil)
	if endpoints == nil || err != nil {
		instanceChan <- getJobInfoFromES(ctx, cost, account, userId)
		return err
	}
	var job *mediaconvert.GetJobOutput
	job = nil
	for _, endpoint := range endpoints.Endpoints {
		subSvc := mediaconvert.New(sess, &aws.Config{Endpoint: endpoint.Url})
		job, err = subSvc.GetJob(&mediaconvert.GetJobInput{Id: &jobId})
		if job != nil {
			break
		}
	}
	if job == nil {
		instanceChan <- getJobInfoFromES(ctx, cost, account, userId)
		return err
	}
	instanceChan <- Job{
		Id: aws.StringValue(job.Job.Id),
		Arn: aws.StringValue(job.Job.Arn),
		Region: cost.Region,
		BillingTagsSource: aws.StringValue(job.Job.BillingTagsSource),
		CreatedAt: aws.TimeValue(job.Job.CreatedAt),
		CurrentPhase: aws.StringValue(job.Job.CurrentPhase),
		ErrorCode: aws.Int64Value(job.Job.ErrorCode),
		ErrorMessage: aws.StringValue(job.Job.ErrorMessage),
		JobPercentComplete: aws.Int64Value(job.Job.JobPercentComplete),
		JobTemplate: aws.StringValue(job.Job.JobTemplate),
		OutputGroupDetails: getOutputGroupDetails(job.Job.OutputGroupDetails),
		Queue: aws.StringValue(job.Job.Queue),
		RetryCount: aws.Int64Value(job.Job.RetryCount),
		Role: aws.StringValue(job.Job.Role),
		Status: aws.StringValue(job.Job.Status),
		StatusUpdateInterval: aws.StringValue(job.Job.StatusUpdateInterval),
		Timing: Timing{
			FinishTime: aws.TimeValue(job.Job.Timing.FinishTime),
			StartTime: aws.TimeValue(job.Job.Timing.StartTime),
			SubmitTime: aws.TimeValue(job.Job.Timing.SubmitTime),
		},
		UserMetadata: getUserMetadata(job.Job.UserMetadata),
		Costs:   cost.Cost,
	}
	return nil
}

// getMediaConvertMetrics gets credentials, accounts and region to fetch MediaConvert instances stats
func fetchMonthlyJobsStats(ctx context.Context, aa taws.AwsAccount, costs []JobInformations, startDate, endDate time.Time) ([]JobReport, error) {
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
	for _, cost := range costs {
		jobRegion := getJobRegion(cost.Arn)
		jobId := getJobId(cost.Arn)
		for _, region := range regions {
			if strings.Contains(region, jobRegion) {
				jobChan := make(chan Job)
				go fetchMonthlyJobsList(ctx, creds, account, region, jobId, cost, jobChan, startDate, endDate, aa.UserId)
				jobChans = append(jobChans, jobChan)
			}
		}
	}
	jobsList := make([]JobReport, 0)
	for job := range merge(jobChans...) {
		jobsList = append(jobsList, JobReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Job: job,
		})
	}
	return jobsList, nil
}

// PutMediaConvertMonthlyReport puts a monthly report of MediaConvert instance in ES
func PutMediaConvertMonthlyReport(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting MediaConvert monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costs := getMediaConvertJobCosts(ctx, aa, startDate, endDate)
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixMediaConvertReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an MediaConvert monthly report", nil)
		return false, nil
	}
	jobs, err := fetchMonthlyJobsStats(ctx, aa, costs, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importJobsToEs(ctx, aa, jobs)
	if err != nil {
		return false, err
	}
	return true, nil
}
