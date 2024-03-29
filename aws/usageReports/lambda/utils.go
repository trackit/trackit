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
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
)

const MonitorFunctionStsSessionName = "monitor-function"

type (
	// FunctionReport is saved in ES to have all the information of an Lambda function
	FunctionReport struct {
		utils.ReportBase
		Function Function `json:"function"`
	}

	// FunctionBase contains basics information of an Lambda function
	FunctionBase struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		Version      string `json:"version"`
		LastModified string `json:"lastModified"`
		Runtime      string `json:"runtime"`
		Size         int64  `json:"size"`
		Memory       int64  `json:"memory"`
		Region       string `json:"region"`
	}

	// Function contains all the information of an Lambda function
	Function struct {
		FunctionBase
		Tags  []utils.Tag `json:"tags"`
		Stats Stats       `json:"stats"`
	}

	// Stats contains statistics of an Lambda function
	Stats struct {
		Invocations Invocations `json:"invocations"`
		Duration    Duration    `json:"duration"`
	}

	Invocations struct {
		Total  float64 `json:"total"`
		Failed float64 `json:"failed"`
	}

	Duration struct {
		Average float64 `json:"average"`
		Maximum float64 `json:"maximum"`
	}
)

// importFunctionsToEs imports Lambda functions in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importFunctionsToEs(ctx context.Context, aa taws.AwsAccount, functions []FunctionReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating Lambda functions for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixLambdaReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, function := range functions {
		id, err := generateId(function)
		if err != nil {
			logger.Error("Error when marshaling function var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, function, TypeLambdaReport, index, id)
	}
	err = bp.Flush()
	if closeErr := bp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		logger.Error("Fail to put Lambda functions in ES", err.Error())
		return err
	}
	logger.Info("Lambda functions put in ES", nil)
	return nil
}

func generateId(function FunctionReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
		Type       string    `json:"reportType"`
	}{
		function.Account,
		function.ReportDate,
		function.Function.Name,
		function.ReportType,
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
func merge(cs ...<-chan Function) <-chan Function {
	var wg sync.WaitGroup
	out := make(chan Function)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Function) {
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
