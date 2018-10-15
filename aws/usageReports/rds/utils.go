//   Copyright 2018 MSolution.IO
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

package rds

import (
	"time"
	"sync"
	"context"
	"crypto/md5"
	"encoding/json"
	"encoding/base64"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/es"
	taws "github.com/trackit/trackit-server/aws"
)

const RDSStsSessionName = "fetch-rds"

type (

	Report struct {
		Account    string     `json:"account"`
		ReportDate time.Time  `json:"reportDate"`
		ReportType string     `json:"reportType"`
		Instances  []Instance `json:"instances"`
	}

	Instance struct {
		DBInstanceIdentifier string             `json:"dbInstanceIdentifier"`
		DBInstanceClass      string             `json:"dbInstanceClass"`
		AllocatedStorage     int64              `json:"allocatedStorage"`
		Engine               string             `json:"engine"`
		AvailabilityZone     string             `json:"availabilityZone"`
		MultiAZ              bool               `json:"multiAZ"`
		Cost                 float64            `json:"cost"`
		CostDetail           map[string]float64 `json:"costDetail"`
		CpuAverage           float64            `json:"cpuAverage"`
		CpuPeak              float64            `json:"cpuPeak"`
		FreeSpaceMin         float64            `json:"freeSpaceMinimum"`
		FreeSpaceMax         float64            `json:"freeSpaceMaximum"`
		FreeSpaceAve         float64            `json:"freeSpaceAverage"`
	}

	InstanceStats struct {
		CpuAverage   float64
		CpuPeak      float64
		FreeSpaceMin float64
		FreeSpaceMax float64
		FreeSpaceAve float64
	}
)

// importReportToEs saves a report into elasticsearch
func importReportToEs(ctx context.Context, aa taws.AwsAccount, report Report) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating RDS report for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	client := es.Client
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
	}{
		report.Account,
		report.ReportDate,
	})
	if err != nil {
		logger.Error("Error when marshaling instance var", err.Error())
		return err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixRDSReport)
	if res, err := client.
		Index().
		Index(index).
		Type(TypeRDSReport).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting RDSReport in ES", err.Error())
		return err
	} else {
		logger.Info("RDSReport put in ES", *res)
		return nil
	}
}

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func merge(cs ...<-chan Instance) <-chan Instance {
	var wg sync.WaitGroup
	out := make(chan Instance)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Instance) {
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
