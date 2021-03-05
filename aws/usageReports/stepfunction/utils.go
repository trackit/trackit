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

const MonitorStepFunctionStsSessionName = "monitor-stepfunction"

type (
	// StepReport is saved in ES to have all the information of a StepFunction
	StepReport struct {
		utils.ReportBase
		Step Step `json:"step"`
	}

	// StepBase contains basics information of a StepFunction
	StepBase struct {
		Name         string    `json:"name"`
		CreationDate time.Time `json:"creationDate"`
		Arn          string    `json:"arn"`
		Region       string    `json:"region"`
	}

	// Step contains all the information of an StepFunction
	Step struct {
		StepBase
		Tags   []utils.Tag `json:"tags"`
		Cost   float64     `json:"cost"`
	}
)

// importStepFunctionToEs imports StepFunction in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importStepFunctionToEs(ctx context.Context, aa taws.AwsAccount, stepFunctions []StepReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating StepFunctions for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixStepFunctionReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, stepFunction := range stepFunctions {
		id, err := generateId(stepFunction)
		if err != nil {
			logger.Error("Error when marshaling stepFunction var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, stepFunction, TypeStepFunctionReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put StepFunctions in ES", err.Error())
		return err
	}
	logger.Info("StepFunctions put in ES", nil)
	return nil
}

func generateId(stepFunction StepReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
		Type       string    `json:"reportType"`
	}{
		stepFunction.Account,
		stepFunction.ReportDate,
		stepFunction.Step.Name,
		stepFunction.ReportType,
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
func merge(cs ...<-chan Step) <-chan Step {
	var wg sync.WaitGroup
	out := make(chan Step)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Step) {
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
