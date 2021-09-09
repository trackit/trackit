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

const MonitorSQSStsSessionName = "monitor-sqs"

type (
	// QueueReport is saved in ES to have all the information of a SQS Queue
	QueueReport struct {
		utils.ReportBase
		Queue Queue `json:"queue"`
	}

	// QueueBase contains basics information of a SQS Queue
	QueueBase struct {
		Name   string `json:"name"`
		Url    string `json:"url"`
		Region string `json:"region"`
	}

	// Queue contains all the information of an SQS Queue
	Queue struct {
		QueueBase
		Tags []utils.Tag `json:"tags"`
	}
)

// importSQSToEs imports SQS Queues in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importSQSToEs(ctx context.Context, aa taws.AwsAccount, queues []QueueReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating SQS for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixSQSReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, queue := range queues {
		id, err := generateId(queue)
		if err != nil {
			logger.Error("Error when marshaling queue var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, queue, TypeSQSReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put SQS in ES", err.Error())
		return err
	}
	logger.Info("SQS put in ES", nil)
	return nil
}

func generateId(queue QueueReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
		Type       string    `json:"reportType"`
	}{
		queue.Account,
		queue.ReportDate,
		queue.Queue.Name,
		queue.ReportType,
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
func merge(cs ...<-chan Queue) <-chan Queue {
	var wg sync.WaitGroup
	out := make(chan Queue)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Queue) {
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
