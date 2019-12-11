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
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
)

const MonitorJobStsSessionName = "monitor-instance"

type (
	// JobReport is saved in ES to have all the information of an MediaConvert instance
	JobReport struct {
		utils.ReportBase
		Job Job `json:"job"`
	}

	VideoDetail struct {
		HeightInPx int64 `json:"heightInPx"`
		WidthInPx  int64 `json:"widthInPx"`
	}

	OutputDetail struct {
		DurationInMs int64 `json:"durationInMs"`
		VideoDetails VideoDetail
	}

	OutputGroupDetail struct {
		OutputDetails []OutputDetail `json:"outputDetails"`
	}

	Timing struct {
		FinishTime time.Time `json:"finishTime"`
		StartTime  time.Time `json:"startTime"`
		SubmitTime time.Time `json:"submitTime"`
	}

	Job struct {
		Arn                  string              `json:"arn"`
		Id                   string              `json:"id"`
		Region               string              `json:"region"`
		BillingTagsSource    string              `json:"billingTagsSource"`
		CreatedAt            time.Time           `json:"createdAt"`
		CurrentPhase         string              `json:"currentPhase"`
		ErrorCode            int64               `json:"errorCode"`
		ErrorMessage         string              `json:"errorMessage"`
		JobPercentComplete   int64               `json:"jobPercentComplete"`
		JobTemplate          string              `json:"jobTemplate"`
		OutputGroupDetails   []OutputGroupDetail `json:"outputGroupDetails"`
		Queue                string              `json:"queue"`
		RetryCount           int64               `json:"retryCount"`
		Role                 string              `json:"role"`
		Status               string              `json:"status"`
		StatusUpdateInterval string              `json:"statusUpdateInterval"`
		Timing               Timing
		UserMetadata         map[string]string `json:"userMetadata"`
		Cost                 float64           `json:"cost"`
	}
)

// importJobsToEs imports MediaConvert instances in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importJobsToEs(ctx context.Context, aa taws.AwsAccount, instances []JobReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating MediaConvert instances for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixMediaConvertReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, instance := range instances {
		id, err := generateId(instance)
		if err != nil {
			logger.Error("Error when marshaling instance var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, instance, TypeMediaConvertReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put MediaConvert instances in ES", err.Error())
		return err
	}
	logger.Info("MediaConvert instances put in ES", nil)
	return nil
}

func generateId(instance JobReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		instance.Account,
		instance.ReportDate,
		instance.Job.Id,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func merge(cs ...<-chan Job) <-chan Job {
	var wg sync.WaitGroup
	out := make(chan Job)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Job) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func getOutputGroupDetails(groupDetails []*mediaconvert.OutputGroupDetail) []OutputGroupDetail {
	var outputGroupDetail []OutputGroupDetail
	for _, groupDetail := range groupDetails {
		var outputDetail OutputGroupDetail
		for _, detail := range groupDetail.OutputDetails {
			outputDetail.OutputDetails = append(outputDetail.OutputDetails, OutputDetail{
				DurationInMs: aws.Int64Value(detail.DurationInMs),
				VideoDetails: VideoDetail{
					HeightInPx: aws.Int64Value(detail.VideoDetails.HeightInPx),
					WidthInPx:  aws.Int64Value(detail.VideoDetails.WidthInPx),
				},
			})
		}
		outputGroupDetail = append(outputGroupDetail, outputDetail)
	}
	return outputGroupDetail
}

func getUserMetadata(initialUserMetadata map[string]*string) map[string]string {
	UserMetadata := make(map[string]string, 0)
	for key, value := range initialUserMetadata {
		UserMetadata[key] = aws.StringValue(value)
	}
	return UserMetadata
}
