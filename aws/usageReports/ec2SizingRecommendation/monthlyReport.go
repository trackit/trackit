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

package ec2SizingRecommendation

import (
	"context"
	//"encoding/json"
	"fmt"
	"github.com/trackit/trackit-server/aws/s3"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
)

type (
	EsQueryParams struct {
		DateBegin         time.Time
		DateEnd           time.Time
		AccountList       []string
		IndexList         []string
		AggregationParams []string
	}

	ReponseEC2SizingRecommendationMonthly struct {

	}
)

// getElasticSearchEc2Instance prepares and run the request to retrieve the a report of an instance
// It will return the data and an error.
func getElasticSearchEc2Instance(ctx context.Context, account, instance string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("instance.id", instance))
	search := client.Search().Index(index).Size(1).Query(query)
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, errors.GetErrorMessage(ctx, err)
	}
	return res, nil
}

// getInstanceInfoFromEs gets information about an instance from previous report to put it in the new report
func getInstanceInfoFromES(ctx context.Context, instance utils.CostPerResource, account string, userId int) Instance {
	/*var docType InstanceReport
	var inst = Instance{
		InstanceBase: InstanceBase{
			Id:         instance.Resource,
			Region:     "N/A",
			State:      "N/A",
			Purchasing: "N/A",
			KeyPair:    "",
			Type:       "N/A",
			Platform:   "Linux/UNIX",
		},
		Tags:  make([]utils.Tag, 0),
		Costs: make(map[string]float64, 0),
		Stats: Stats{
			Cpu: Cpu{
				Average: -1,
				Peak:    -1,
			},
			Network: Network{
				In:  -1,
				Out: -1,
			},
			Volumes: make([]Volume, 0),
		},
	}
	inst.Costs["instance"] = instance.Cost
	res, err := getElasticSearchEc2Instance(ctx, account, instance.Resource,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixEC2Report))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			inst.Region = docType.Instance.Region
			inst.Purchasing = docType.Instance.Purchasing
			inst.KeyPair = docType.Instance.KeyPair
			inst.Type = docType.Instance.Type
			inst.Platform = docType.Instance.Platform
			inst.Tags = docType.Instance.Tags
		}
	}
	return inst*/
	return Instance{}
}

func FormatInstanceSizeResult(res *elastic.SearchResult) []InstanceReport {
	return []InstanceReport{}
}

// getEc2Metrics gets credentials, accounts and region to fetch EC2 instances stats
func fetchMonthlyInstancesStats(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) ([]InstanceReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUserId(aa.UserId, s3.IndexPrefixLineItem)
	parsedParams := EsQueryParams{
		AccountList: []string{aa.AwsIdentity},
		IndexList:   []string{index},
		DateBegin:   startDate,
		DateEnd:     endDate,
	}
	search := GetElasticSearchParams(parsedParams.AccountList, startDate, endDate, es.Client, index)
	res, err := search.Do(ctx)
	if err != nil {
		logger.Error("Error when doing the search", err)
	}
	logger.Debug("RESULT ==========", map[string]interface{}{
		"result": res,
	})
	reports := FormatInstanceSizeResult(res)
	return reports, nil
}

// filterInstancesCosts filters instances, cloudwatch and volumes of EC2 instances costs
func filterInstancesCosts(ec2Cost, cloudwatchCost []utils.CostPerResource) ([]utils.CostPerResource, []utils.CostPerResource, []utils.CostPerResource) {
	newInstance := make([]utils.CostPerResource, 0)
	newVolume := make([]utils.CostPerResource, 0)
	newCloudWatch := make([]utils.CostPerResource, 0)
	for _, instance := range ec2Cost {
		if len(instance.Resource) == 19 && strings.HasPrefix(instance.Resource, "i-") {
			newInstance = append(newInstance, instance)
		}
		if len(instance.Resource) == 21 && strings.HasPrefix(instance.Resource, "vol-") {
			newVolume = append(newVolume, instance)
		}
	}
	for _, instance := range cloudwatchCost {
		for _, cost := range newInstance {
			if strings.Contains(instance.Resource, cost.Resource) {
				newCloudWatch = append(newCloudWatch, instance)
			}
		}
	}
	return newInstance, newVolume, newCloudWatch
}

// PutEc2MonthlyReport puts a monthly report of EC2 instance in ES
func PutEc2SizingRecommendationMonthlyReport(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting EC2 monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	/*already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixEC2SizeRecommendationReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an EC2 monthly report", nil)
		return false, nil
	}*/
	_, err := fetchMonthlyInstancesStats(ctx, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	//instances = addCostToInstances(instances, costVolume, costCloudWatch)
	//err = importInstancesToEs(ctx, aa, instances)
	/*if err != nil {
		return false, err
	}*/
	return true, nil
}
