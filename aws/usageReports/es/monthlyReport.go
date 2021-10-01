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

package es

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/esReports"
)

func fetchMonthlyDomainsList(ctx context.Context, creds *credentials.Credentials, dom common.CostPerResource,
	region string, domainChan chan esReports.Domain, start, end time.Time) error {
	defer close(domainChan)
	var tags []common.Tag
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := elasticsearchservice.New(sess)
	domainStatus, err := svc.DescribeElasticsearchDomain(&elasticsearchservice.DescribeElasticsearchDomainInput{
		DomainName: &strings.Split(dom.Resource, "/")[1],
	})
	if err != nil {
		return err
	}
	domain := domainStatus.DomainStatus
	if esTags, err := svc.ListTags(&elasticsearchservice.ListTagsInput{ARN: domain.ARN}); err != nil {
		logger.Error("Error while listing Tags for domain", err.Error())
		tags = make([]common.Tag, 0)
	} else {
		tags = getDomainTag(esTags.TagList)
	}
	stats := getDomainStats(ctx, *domain.DomainName, sess, start, end)
	costs := make(map[string]float64, 1)
	costs["domain"] = dom.Cost
	domainChan <- esReports.Domain{
		DomainBase: esReports.DomainBase{
			Arn:               aws.StringValue(domain.ARN),
			InstanceType:      aws.StringValue(domain.ElasticsearchClusterConfig.InstanceType),
			InstanceCount:     aws.Int64Value(domain.ElasticsearchClusterConfig.InstanceCount),
			DomainID:          aws.StringValue(domain.DomainId),
			DomainName:        aws.StringValue(domain.DomainName),
			TotalStorageSpace: aws.Int64Value(domain.EBSOptions.VolumeSize),
			Region:            region,
		},
		Tags:  tags,
		Costs: costs,
		Stats: stats,
	}
	return nil
}

// getEsMetrics gets credentials, accounts and region to fetch RDS instances stats
func getEsMetrics(ctx context.Context, domainsList []common.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) ([]esReports.DomainReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, ESStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return nil, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return nil, err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return nil, err
	}
	domainChans := make([]<-chan esReports.Domain, 0, len(regions))
	for _, domain := range domainsList {
		for _, region := range regions {
			domainChan := make(chan esReports.Domain)
			go fetchMonthlyDomainsList(ctx, creds, domain, region, domainChan, startDate, endDate)
			domainChans = append(domainChans, domainChan)
		}
	}
	domains := make([]esReports.DomainReport, 0)
	for domain := range merge(domainChans...) {
		domains = append(domains, esReports.DomainReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Domain: domain,
		})
	}
	return domains, nil
}

// filterEsDomains filters cost per domain to get only costs associated to an ES domain
func filterEsDomains(esCost []common.CostPerResource) []common.CostPerResource {
	costDomains := []common.CostPerResource{}
	for _, domain := range esCost {
		split := strings.Split(domain.Resource, ":")
		if len(split) == 6 && split[2] == "es" {
			costDomains = append(costDomains, common.CostPerResource{
				Resource: split[5],
				Cost:     domain.Cost,
				Region:   domain.Region},
			)
		}
	}
	return costDomains
}

// PutEsMonthlyReport puts a monthly report of ES in ES
func PutEsMonthlyReport(ctx context.Context, esCost []common.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting ES monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance := filterEsDomains(esCost)
	if len(costInstance) == 0 {
		logger.Info("No ES domains found in billing data.", nil)
		return false, nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, esReports.Model.IndexSuffix)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an ES monthly report", nil)
		return false, nil
	}
	report, err := getEsMetrics(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importDomainsToEs(ctx, aa, report)
	if err != nil {
		return false, err
	}
	return true, nil
}
