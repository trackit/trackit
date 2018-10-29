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

package es

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/es"
)

const ESStsSessionName = "fetch-es"
const MonitorDomainStsSessionName = "monitor-domain"

type (
	// Report represents the report with all the information for ES domains.
	Report struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		ReportType string    `json:"reportType"`
		Domains    []Domain  `json:"instances"`
	}

	// Domain represents all the informations of an ES domain.
	// It will be imported in ElasticSearch thanks to the struct tags.
	Domain struct {
		Arn                      string             `json:"arn"`
		Region                   string             `json:"region"`
		DomainID                 string             `json:"domainId"`
		DomainName               string             `json:"domainName"`
		CPUUtilizationAverage    float64            `json:"cpuUtilizationAverage"`
		CPUUtiliztionPeak        float64            `json:"cpuUtilizationPeak"`
		FreeStorageSpace         float64            `json:"freeStorageSpace"`
		TotalStorageSpace        int64              `json:"totalStorageSpace"`
		JVMMemoryPressureAverage float64            `json:"jvmMemoryPressureAverage"`
		JVMMemoryPressurePeak    float64            `json:"jvmMemoryPressurePeak"`
		InstanceType             string             `json:"instanceType"`
		InstanceCount            int64              `json:"instanceCount"`
		Tags                     map[string]string  `json:"tags"`
		Cost                     float64            `json:"cost"`
		CostDetail               map[string]float64 `json:"costDetail"`
	}

	// DomainStats contains statistic of a domain gotten by CloudWatch
	DomainStats struct {
		CPUUtilizationAverage    float64
		CPUUtiliztionPeak        float64
		FreeStorageSpace         float64
		JVMMemoryPressureAverage float64
		JVMMemoryPressurePeak    float64
	}
)

func transformDomainsListToString(domainNames []*elasticsearchservice.DomainInfo) []*string {
	res := make([]*string, 0)
	for _, domain := range domainNames {
		res = append(res, domain.DomainName)
	}
	return res
}

// importReportToEs imports an Report in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importReportToEs(ctx context.Context, aa taws.AwsAccount, report Report) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating ES domains for AWS account.", map[string]interface{}{
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
		logger.Error("Error when marshaling domain var", err.Error())
		return err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixESReport)
	if res, err := client.
		Index().
		Index(index).
		Type(TypeESReport).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting Domain in ES", err.Error())
	} else {
		logger.Info("Domain put in ES", *res)
	}
	return nil
}

// getDomainTag formats []*elasticsearchservice.Tag to map[string]string
func getDomainTag(tags []*elasticsearchservice.Tag) map[string]string {
	res := make(map[string]string)
	for _, tag := range tags {
		res[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
	}
	return res
}

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func merge(cs ...<-chan Domain) <-chan Domain {
	var wg sync.WaitGroup
	out := make(chan Domain)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Domain) {
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
