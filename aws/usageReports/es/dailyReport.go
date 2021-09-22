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

// Package es implements usage reports for Amazon ElasticSearch
package es

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchDailyDomainsList fetches the list of domains for a specific region
func fetchDailyDomainsList(ctx context.Context, creds *credentials.Credentials, region string, domainChan chan Domain) error {
	var tags []utils.Tag
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
	var domainsStatus []*elasticsearchservice.ElasticsearchDomainStatus
	for i := range domainNames.DomainNames {
		if i%5 == 0 {
			m := i + 5
			if m >= len(domainNames.DomainNames) {
				m = len(domainNames.DomainNames)
			}
			domains, err := svc.DescribeElasticsearchDomains(&elasticsearchservice.DescribeElasticsearchDomainsInput{
				DomainNames: transformDomainsListToString(domainNames.DomainNames[i:m]),
			})
			if err != nil {
				logger.Error("Error when describing domains", err.Error())
				return err
			}
			domainsStatus = append(domainsStatus, domains.DomainStatusList...)
		}
	}
	for _, domain := range domainsStatus {
		if esTags, err := svc.ListTags(&elasticsearchservice.ListTagsInput{ARN: domain.ARN}); err != nil {
			logger.Error("Error while listing Tags for domain", err.Error())
			tags = make([]utils.Tag, 0)
		} else {
			tags = getDomainTag(esTags.TagList)
		}
		stats := getDomainStats(ctx, *domain.DomainName, sess, start, end)
		domainChan <- Domain{
			DomainBase: DomainBase{
				Arn:               aws.StringValue(domain.ARN),
				InstanceType:      aws.StringValue(domain.ElasticsearchClusterConfig.InstanceType),
				InstanceCount:     aws.Int64Value(domain.ElasticsearchClusterConfig.InstanceCount),
				DomainID:          aws.StringValue(domain.DomainId),
				DomainName:        aws.StringValue(domain.DomainName),
				TotalStorageSpace: aws.Int64Value(domain.EBSOptions.VolumeSize),
				Region:            region,
			},
			Tags:  tags,
			Costs: make(map[string]float64),
			Stats: stats,
		}
	}
	return nil
}

// FetchDomainsStats retrieves ES information from the AWS API and generates a report
func FetchDomainsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching ES instance stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorDomainStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	now := time.Now().UTC()
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return err
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
	domains := make([]DomainReport, 0)
	for domain := range merge(domainChans...) {
		domains = append(domains, DomainReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Domain: domain,
		})
	}
	return importDomainsToEs(ctx, awsAccount, domains)
}
