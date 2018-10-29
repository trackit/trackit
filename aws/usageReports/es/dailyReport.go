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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/config"
)

// fetchDailyDomainsList fetches the list of domains for a specific region
func fetchDailyDomainsList(ctx context.Context, creds *credentials.Credentials, region string, domainChan chan Domain) error {
	defer close(domainChan)
	start, end := utils.GetCurrentCheckedDay()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := elasticsearchservice.New(sess)
	domainNames, err := svc.ListDomainNames(nil)
	if err != nil {
		logger.Error("Error when listing domain names", err.Error())
		return err
	}
	domains, err := svc.DescribeElasticsearchDomains(&elasticsearchservice.DescribeElasticsearchDomainsInput{
		DomainNames: transformDomainsListToString(domainNames.DomainNames),
	})
	if err != nil {
		logger.Error("Error when describing domains", err.Error())
	}
	for _, domain := range domains.DomainStatusList {
		tags, err := svc.ListTags(&elasticsearchservice.ListTagsInput{
			ARN: domain.ARN,
		})
		if err != nil {
			logger.Error("Error while listing Tags for domain", err.Error())
			return err
		}
		stats, err := getDomainStats(ctx, *domain.DomainName, sess, start, end)
		if err != nil {
			logger.Error("Error while getting Domain stats", err.Error())
			return err
		}
		domainChan <- Domain{
			Arn:               aws.StringValue(domain.ARN),
			InstanceType:      aws.StringValue(domain.ElasticsearchClusterConfig.InstanceType),
			InstanceCount:     aws.Int64Value(domain.ElasticsearchClusterConfig.InstanceCount),
			DomainID:          aws.StringValue(domain.DomainId),
			DomainName:        aws.StringValue(domain.DomainName),
			TotalStorageSpace: aws.Int64Value(domain.EBSOptions.VolumeSize),
			Region:            region,
			Tags:              getDomainTag(tags.TagList),
			CPUUtilizationAverage:    stats.CPUUtilizationAverage,
			CPUUtiliztionPeak:        stats.CPUUtiliztionPeak,
			FreeStorageSpace:         stats.FreeStorageSpace,
			JVMMemoryPressureAverage: stats.JVMMemoryPressureAverage,
			JVMMemoryPressurePeak:    stats.JVMMemoryPressurePeak,
			Cost:       0,
			CostDetail: make(map[string]float64, 0),
		}
	}
	return nil
}

// FetchDomainsStats retrieces ES information from the AWS API and generate a report
func FetchDomainsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching EC2 instance stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorDomainStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return err
	}
	report := Report{
		account,
		time.Now().UTC(),
		"daily",
		make([]Domain, 0),
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	domainChans := make([]<-chan Domain, 0, len(regions))
	for _, region := range regions {
		domainChan := make(chan Domain)
		go fetchDailyDomainsList(ctx, creds, region, domainChan)
		domainChans = append(domainChans, domainChan)
	}
	for domain := range merge(domainChans...) {
		report.Domains = append(report.Domains, domain)
	}
	return importReportToEs(ctx, awsAccount, report)
}
