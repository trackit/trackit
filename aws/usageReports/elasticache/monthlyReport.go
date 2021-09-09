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

package elasticache

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchMonthlyInstancesList sends in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances and getInstanceStats.
func fetchMonthlyInstancesList(ctx context.Context, creds *credentials.Credentials, inst utils.CostPerResource,
	account, region string, instanceChan chan Instance, startDate, endDate time.Time) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := elasticache.New(sess)
	input := elasticache.DescribeCacheClustersInput{ShowCacheNodeInfo: aws.Bool(true)}
	instances, err := svc.DescribeCacheClusters(&input)
	if err != nil {
		return err
	}
	for _, cluster := range instances.CacheClusters {
		stats := getInstanceStats(ctx, cluster, sess, startDate, endDate)
		costs := make(map[string]float64, 1)
		costs[inst.Resource] = inst.Cost
		tags := getClusterTags(ctx, cluster, svc, account, region)
		instanceChan <- Instance{
			InstanceBase: InstanceBase{
				Id:            aws.StringValue(cluster.CacheClusterId),
				Status:        aws.StringValue(cluster.CacheClusterStatus),
				Region:        aws.StringValue(cluster.PreferredAvailabilityZone),
				NodeType:      aws.StringValue(cluster.CacheNodeType),
				Nodes:         extractCacheNodes(cluster.CacheNodes),
				Engine:        aws.StringValue(cluster.Engine),
				EngineVersion: aws.StringValue(cluster.EngineVersion),
			},
			Costs: costs,
			Tags:  tags,
			Stats: stats,
		}
	}
	return nil
}

// getEc2Metrics gets credentials, accounts and region to fetch ElastiCache instances stats
func fetchMonthlyInstancesStats(ctx context.Context, instances []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) ([]InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorInstanceStsSessionName)
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
	instanceChans := make([]<-chan Instance, 0, len(regions))
	for _, instance := range instances {
		for _, region := range regions {
			if strings.Contains(instance.Region, region) {
				instanceChan := make(chan Instance)
				go fetchMonthlyInstancesList(ctx, creds, instance, account, region, instanceChan, startDate, endDate)
				instanceChans = append(instanceChans, instanceChan)
			}
		}
	}
	instancesList := make([]InstanceReport, 0)
	for instance := range merge(instanceChans...) {
		instancesList = append(instancesList, InstanceReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Instance: instance,
		})
	}
	return instancesList, nil
}

// filterElastiCacheInstances filters cost per instance to get only costs associated to an ElastiCache Instance
func filterElastiCacheInstances(elastiCacheCost []utils.CostPerResource) []utils.CostPerResource {
	costInstances := make([]utils.CostPerResource, 0)
	for _, domain := range elastiCacheCost {
		// format in billing data for an ElastiCache instance is:
		// "arn:aws:elasticache:[region]:[aws_id]:cluster:[cluster name]"
		split := strings.Split(domain.Resource, ":")
		if len(split) == 7 && split[2] == "elasticache" {
			costInstances = append(costInstances, utils.CostPerResource{
				Resource: split[6],
				Cost:     domain.Cost,
				Region:   domain.Region,
			})
		}
	}
	return costInstances
}

// PutElastiCacheMonthlyReport puts a monthly report of ElastiCache instance in ES
func PutElastiCacheMonthlyReport(ctx context.Context, costs []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting ElastiCache monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costs = filterElastiCacheInstances(costs)
	if len(costs) == 0 {
		logger.Info("No ElastiCache instances found in billing data.", nil)
		return false, nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixElastiCacheReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an ElastiCache monthly report", nil)
		return false, nil
	}
	instances, err := fetchMonthlyInstancesStats(ctx, costs, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	if err = importInstancesToEs(ctx, aa, instances); err != nil {
		return false, err
	} else {
		return true, nil
	}
}
