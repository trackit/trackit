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
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/es"
)

const ESStsSessionName = "fetch-es"
const MonitorDomainStsSessionName = "monitor-domain"

type (
	// DomainReport represents the report with all the information for ES domains.
	DomainReport struct {
		utils.ReportBase
		Domain Domain `json:"domain"`
	}

	// DomainBase contains basics information of an ES domain
	DomainBase struct {
		Arn               string `json:"arn"`
		Region            string `json:"region"`
		DomainID          string `json:"domainId"`
		DomainName        string `json:"domainName"`
		InstanceType      string `json:"instanceType"`
		InstanceCount     int64  `json:"instanceCount"`
		TotalStorageSpace int64  `json:"totalStorageSpace"`
	}

	// Domain contains all information of an ES domain that will be save in ES
	Domain struct {
		DomainBase
		Tags  []utils.Tag        `json:"tags"`
		Costs map[string]float64 `json:"costs"`
		Stats Stats              `json:"stats"`
	}

	// Stats contains statistics of a domain get on CloudWatch
	Stats struct {
		Cpu               Cpu               `json:"cpu"`
		FreeSpace         float64           `json:"freeSpace"`
		JVMMemoryPressure JVMMemoryPressure `json:"JVMMemoryPressure"`
	}

	// Cpu contains cpu statistics of a domain
	Cpu struct {
		Average float64 `json:"average"`
		Peak    float64 `json:"peak"`
	}

	// JVMMemoryPressure contains JVMMemoryPressure statistics of a domain
	JVMMemoryPressure struct {
		Average float64 `json:"average"`
		Peak    float64 `json:"peak"`
	}
)

func transformDomainsListToString(domainNames []*elasticsearchservice.DomainInfo) []*string {
	res := make([]*string, 0)
	for _, domain := range domainNames {
		res = append(res, domain.DomainName)
	}
	return res
}

// importInstancesToEs imports EC2 instances in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importDomainsToEs(ctx context.Context, aa taws.AwsAccount, domains []DomainReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating ES domains for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixESReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, domain := range domains {
		id, err := generateId(domain)
		if err != nil {
			logger.Error("Error when marshaling domain var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, domain, TypeESReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put ES domains in ES", err.Error())
		return err
	}
	logger.Info("ES domains put in ES", nil)
	return nil
}

func generateId(domain DomainReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		domain.Account,
		domain.ReportDate,
		domain.Domain.DomainID,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}

// getDomainTag formats []*elasticsearchservice.Tag to map[string]string
func getDomainTag(tags []*elasticsearchservice.Tag) []utils.Tag {
	res := make([]utils.Tag, 0)
	for _, tag := range tags {
		res = append(res, utils.Tag{aws.StringValue(tag.Key), aws.StringValue(tag.Value)})
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
