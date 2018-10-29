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
	"strings"
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

func fetchMonthlyDomainsList(ctx context.Context, creds *credentials.Credentials, domainList []utils.CostPerResource, region string, domainChan chan Domain, start, end time.Time) error {
	defer close(domainChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := elasticsearchservice.New(sess)
	for _, domainCost := range domainList {
		domainStatus, err := svc.DescribeElasticsearchDomain(&elasticsearchservice.DescribeElasticsearchDomainInput{
			DomainName: &strings.Split(domainCost.Resource, "/")[1],
		})
		if err != nil {
			continue // Error from missing region in domainList, needs removal when region is added to it
		}
		domain := domainStatus.DomainStatus
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
		detail := make(map[string]float64, 0)
		detail["domain"] = domainCost.Cost
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
			Cost:       domainCost.Cost,
			CostDetail: detail,
		}
	}
	return nil
}

// getEsMetrics gets credentials, accounts and region to fetch RDS instances stats
func getEsMetrics(ctx context.Context, domains []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, ESStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return Report{}, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return Report{}, err
	}
	report := Report{
		Account:    account,
		ReportDate: startDate,
		ReportType: "monthly",
		Domains:    make([]Domain, 0),
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return Report{}, err
	}
	domainChans := make([]<-chan Domain, 0, len(regions))
	for _, region := range regions {
		domainChan := make(chan Domain)
		go fetchMonthlyDomainsList(ctx, creds, domains, region, domainChan, startDate, endDate)
		domainChans = append(domainChans, domainChan)
	}
	for domain := range merge(domainChans...) {
		report.Domains = append(report.Domains, domain)
	}
	return report, nil
}

func filterEsDomains(esCost []utils.CostPerResource) []utils.CostPerResource {
	costDomains := []utils.CostPerResource{}
	for _, domain := range esCost {
		split := strings.Split(domain.Resource, ":")
		if len(split) == 6 && split[2] == "es" {
			costDomains = append(costDomains, utils.CostPerResource{
				Resource: split[5],
				Cost:     domain.Cost,
			})
		}
	}
	return costDomains
}

// PutEsMonthlyReport puts a monthly report of ES in ES
func PutEsMonthlyReport(ctx context.Context, esCost []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting ES monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costInstance := filterEsDomains(esCost)
	if len(costInstance) == 0 {
		logger.Info("No ES domains found in billing data.", nil)
		return nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixESReport)
	if err != nil {
		return err
	} else if already {
		logger.Info("There is already an ES monthly report", nil)
		return nil
	}
	report, err := getEsMetrics(ctx, costInstance, aa, startDate, endDate)
	if err != nil {
		return err
	}
	return importReportToEs(ctx, aa, report)
}
