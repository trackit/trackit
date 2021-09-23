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

package history

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/aws/usageReports/ebs"
	"github.com/trackit/trackit/aws/usageReports/ec2"
	"github.com/trackit/trackit/aws/usageReports/ec2Coverage"
	"github.com/trackit/trackit/aws/usageReports/elasticache"
	tes "github.com/trackit/trackit/aws/usageReports/es"
	"github.com/trackit/trackit/aws/usageReports/instanceCount"
	"github.com/trackit/trackit/aws/usageReports/rds"
	"github.com/trackit/trackit/es"
)

const numPartition = 5

var ErrBillingDataIncomplete = errors.New("Billing data are not completed")

type (
	// EsRegionPerResourceResult allows us to parse ES results
	EsRegionPerResourceResult struct {
		Resources struct {
			Buckets []struct {
				Resource string                `json:"key"`
				Regions  EsCostPerRegionResult `json:"regions"`
			} `json:"buckets"`
		} `json:"products"`
	}

	EsCostPerRegionResult struct {
		Buckets []struct {
			Region string `json:"key"`
			Cost   struct {
				Value float64 `json:"value"`
			} `json:"cost"`
		} `json:"buckets"`
	}
)

// GetHistoryDate return the begin and the end date of the last month
func GetHistoryDate() (time.Time, time.Time) {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location()).UTC()
	end := time.Date(now.Year(), now.Month(), 0, 23, 59, 59, 999999999, now.Location()).UTC()
	return start, end
}

// makeElasticSearchRequestForCost will make the actual request to the ElasticSearch
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 Internal Server Error status code, it will return the provided status code
// with empty data
func makeElasticSearchRequestForCost(ctx context.Context, client *elastic.Client, aa aws.AwsAccount,
	startDate, endDate time.Time, product string, partition int) (*elastic.SearchResult, int, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUserId(aa.UserId, es.IndexPrefixLineItems)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("usageAccountId", aa.AwsIdentity))
	query = query.Filter(elastic.NewTermQuery("productCode", product))
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(startDate).To(endDate))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("resources", elastic.NewTermsAggregation().Field("resourceId").Size(utils.MaxAggregationSize).Partition(partition).NumPartitions(numPartition).
		SubAggregation("regions", elastic.NewTermsAggregation().Field("availabilityZone").Size(utils.MaxAggregationSize).
			SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost"))))
	result, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			logger.Warning("Query execution failed, ES index does not exists", map[string]interface{}{"index": index, "error": err.Error()})
			return nil, http.StatusOK, err
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
			logger.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			logger.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, err
	}
	return result, http.StatusOK, nil
}

// getCostPerResource returns the parsed result of ES
// This response contains the list of the resources of the specified product with the cost and region associated
func getCostPerResource(ctx context.Context, aa aws.AwsAccount, startDate time.Time, endDate time.Time,
	product string) ([]utils.CostPerResource, error) {
	var parsedResult EsRegionPerResourceResult
	response := make([]utils.CostPerResource, 0)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	for i := 0; i < numPartition; i++ {
		result, returnCode, err := makeElasticSearchRequestForCost(ctx, es.Client, aa, startDate, endDate, product, i)
		if err != nil {
			if returnCode != http.StatusOK {
				return response, err
			} else {
				return response, nil
			}
		}
		err = json.Unmarshal(*result.Aggregations["resources"], &parsedResult.Resources)
		if err != nil {
			logger.Error("Error while unmarshaling", err)
			return response, errors.New("Internal server error")
		}
		for _, resource := range parsedResult.Resources.Buckets {
			element := utils.CostPerResource{resource.Resource, 0, ""}
			for _, region := range resource.Regions.Buckets {
				if region.Region != "" {
					element.Region = region.Region
				}
				element.Cost += region.Cost.Value
			}
			response = append(response, element)
		}
	}
	return response, nil
}

func concatErrors(tabError []error) error {
	var stringError = ""
	for _, err := range tabError {
		if err != nil {
			if stringError != "" {
				stringError += " + "
			}
			stringError += err.Error()
		}
	}
	if stringError != "" {
		if len(stringError) > 254 {
			stringError = stringError[0:254]
		}
		return errors.New(stringError)
	}
	return nil
}

// getInstanceInfo sort products and call history reports
func getInstancesInfo(ctx context.Context, aa aws.AwsAccount, startDate time.Time, endDate time.Time) (bool, error) {
	var ebsCreated, ec2Created, rdsCreated, esCreated, elastiCacheCreated bool
	var ebsErr error
	ec2Cost, ec2Err := getCostPerResource(ctx, aa, startDate, endDate, "AmazonEC2")
	cloudWatchCost, cloudWatchErr := getCostPerResource(ctx, aa, startDate, endDate, "AmazonCloudWatch")
	if ec2Err == nil && cloudWatchErr == nil {
		ec2Created, ec2Err = ec2.PutEc2MonthlyReport(ctx, ec2Cost, cloudWatchCost, aa, startDate, endDate)
		ebsCreated, ebsErr = ebs.PutEbsMonthlyReport(ctx, ec2Cost, aa, startDate, endDate)
	}
	rdsCost, rdsErr := getCostPerResource(ctx, aa, startDate, endDate, "AmazonRDS")
	if rdsErr == nil {
		rdsCreated, rdsErr = rds.PutRdsMonthlyReport(ctx, rdsCost, aa, startDate, endDate)
	}
	esCost, esErr := getCostPerResource(ctx, aa, startDate, endDate, "AmazonES")
	if esErr == nil {
		esCreated, esErr = tes.PutEsMonthlyReport(ctx, esCost, aa, startDate, endDate)
	}
	elastiCacheCost, elastiCacheErr := getCostPerResource(ctx, aa, startDate, endDate, "AmazonElastiCache")
	if elastiCacheErr == nil {
		elastiCacheCreated, elastiCacheErr = elasticache.PutElastiCacheMonthlyReport(ctx, elastiCacheCost, aa, startDate, endDate)
	}
	ec2CoverageCreated, ec2CoverageErr := ec2Coverage.PutEc2MonthlyCoverageReport(ctx, aa, startDate, endDate)
	instanceCountCreated, instanceCountErr := instanceCount.PutInstanceCountMonthlyReport(ctx, aa, startDate, endDate)
	reportsCreated := ebsCreated || ec2Created || rdsCreated || esCreated || elastiCacheCreated || ec2CoverageCreated || instanceCountCreated
	return reportsCreated, concatErrors([]error{ec2Err, ebsErr, cloudWatchErr, rdsErr, esErr, elastiCacheErr, ec2CoverageErr, instanceCountErr})
}

// CheckBillingDataCompleted checks if billing data in ES are complete.
// If they are complete it returns true, otherwise it returns false.
func CheckBillingDataCompleted(ctx context.Context, startDate time.Time, endDate time.Time, aa aws.AwsAccount) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("usageAccountId", aa.AwsIdentity))
	query = query.Filter(elastic.NewTermQuery("invoiceId", ""))
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(startDate).To(endDate))
	index := es.IndexNameForUserId(aa.UserId, es.IndexPrefixLineItems)
	result, err := es.Client.Search().Index(index).Size(1).Query(query).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			logger.Warning("Query execution failed, ES index does not exists", map[string]interface{}{"index": index, "error": err.Error()})
			return false, nil
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
			logger.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			logger.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return false, err
	}
	if result.Hits.TotalHits == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// FetchHistoryInfos fetches billing data and stats of EC2 and RDS instances of the last month
func FetchHistoryInfos(ctx context.Context, aa aws.AwsAccount, date time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var startDate, endDate time.Time
	if date.IsZero() {
		startDate, endDate = GetHistoryDate()
	} else {
		startDate = date
		endDate = time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location())
	}
	logger.Info("Starting history report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	complete, err := CheckBillingDataCompleted(ctx, startDate, endDate, aa)
	if err != nil {
		return false, err
	} else if !complete {
		logger.Info("Billing data are not completed", nil)
		return false, ErrBillingDataIncomplete
	}
	return getInstancesInfo(ctx, aa, startDate, endDate)
}
