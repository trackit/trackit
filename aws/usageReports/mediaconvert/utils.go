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

	// JobBase contains basics information of an MediaConvert instance
	JobBase struct {
		Arn    string `json:"arn"`
		Region string `json:"region"`
		Id     string `json:"id"`
	}

	VideoDetail struct {
		HeightInPx int64 `locationName:"heightInPx" type:"integer"`
		WidthInPx int64 `locationName:"widthInPx" type:"integer"`
	}

	OutputDetail struct {
		DurationInMs int64 `locationName:"durationInMs" type:"integer"`
		VideoDetails VideoDetail `locationName:"videoDetails" type:"structure"`
	}

	OutputGroupDetail struct {
		OutputDetails []OutputDetail `locationName:"outputDetails" type:"list"`
	}

	Timing struct {
		FinishTime time.Time `locationName:"finishTime" type:"timestamp" timestampFormat:"unixTimestamp"`
		StartTime time.Time `locationName:"startTime" type:"timestamp" timestampFormat:"unixTimestamp"`
		SubmitTime time.Time `locationName:"submitTime" type:"timestamp" timestampFormat:"unixTimestamp"`
	}

	Job struct {
		JobBase JobBase
		AccelerationStatus string `locationName:"accelerationStatus" type:"string" enum:"AccelerationStatus"`
		Arn string `locationName:"arn" type:"string"`
		BillingTagsSource string `locationName:"billingTagsSource" type:"string" enum:"BillingTagsSource"`
		CreatedAt time.Time `locationName:"createdAt" type:"timestamp" timestampFormat:"unixTimestamp"`
		CurrentPhase string `locationName:"currentPhase" type:"string" enum:"JobPhase"`
		ErrorCode int64 `locationName:"errorCode" type:"integer"`
		ErrorMessage string `locationName:"errorMessage" type:"string"`
		Id string `locationName:"id" type:"string"`
		JobPercentComplete int64 `locationName:"jobPercentComplete" type:"integer"`
		JobTemplate string `locationName:"jobTemplate" type:"string"`
		OutputGroupDetails []OutputGroupDetail `locationName:"outputGroupDetails" type:"list"`
		Queue string `locationName:"queue" type:"string"`
		RetryCount int64 `locationName:"retryCount" type:"integer"`
		Role string `locationName:"role" type:"string" required:"true"`
		Status string `locationName:"status" type:"string" enum:"JobStatus"`
		StatusUpdateInterval string `locationName:"statusUpdateInterval" type:"string" enum:"StatusUpdateInterval"`
		Timing Timing `locationName:"timing" type:"structure"`
		UserMetadata map[string]string `locationName:"userMetadata" type:"map"`
		Costs map[time.Time]float64 `json:"costs"`
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

// getPlatformName normalizes the platform name
func getPlatformName(platform string) string {
	if platform == "" {
		return "Linux/UNIX"
	}
	return platform
}
