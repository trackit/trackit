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

package ec2

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/es"
)

const MonitorInstanceStsSessionName = "monitor-instance"

type (
	// Report represents the report with all the information for EC2 instances.
	Report struct {
		Account    string     `json:"account"`
		ReportDate time.Time  `json:"reportDate"`
		ReportType string     `json:"reportType"`
		Instances  []Instance `json:"instances"`
	}

	// Instance represents all the information of an EC2 instance.
	Instance struct {
		Id         string             `json:"id"`
		Region     string             `json:"region"`
		State      string             `json:"state"`
		Purchasing string             `json:"purchasing"`
		CpuAverage float64            `json:"cpuAverage"`
		CpuPeak    float64            `json:"cpuPeak"`
		NetworkIn  float64            `json:"networkIn"`
		NetworkOut float64            `json:"networkOut"`
		IORead     map[string]float64 `json:"ioRead"`
		IOWrite    map[string]float64 `json:"ioWrite"`
		KeyPair    string             `json:"keyPair"`
		Type       string             `json:"type"`
		Tags       map[string]string  `json:"tags"`
		Cost       float64            `json:"cost"`
		CostDetail map[string]float64 `json:"costDetail"`
	}

	// instanceStats contains statistics of an instance get by CloudWatch
	instanceStats struct {
		CpuAverage float64
		CpuPeak    float64
		NetworkIn  float64
		NetworkOut float64
		IORead     map[string]float64
		IOWrite    map[string]float64
	}
)

// importReportToEs imports an EC2 report in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importReportToEs(ctx context.Context, aa taws.AwsAccount, report Report) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating EC2 instances for AWS account.", map[string]interface{}{
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
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixEC2Report)
	if res, err := client.
		Index().
		Index(index).
		Type(TypeEC2Report).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting InstanceInfo in ES", err.Error())
		return err
	} else {
		logger.Info("Instance put in ES", *res)
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
