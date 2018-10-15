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

package history

import (
	"time"
	"errors"
	"context"
	"net/http"
	"encoding/json"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/aws/usageReports/ec2"
	"github.com/trackit/trackit-server/aws/usageReports/rds"
)

type (
	// structures that allows to parse ES result
	EsInstancePerProductResult struct {
		Products struct {
			Buckets []struct {
				Product   string                  `json:"key"`
				Instances EsCostPerInstanceResult `json:"instances"`
			} `json:"buckets"`
		} `json:"products"`
	}

	EsCostPerInstanceResult struct {
		Buckets []struct {
			Instance string `json:"key"`
			Cost     struct {
				Value float64 `json:"value"`
			} `json:"cost"`
		} `json:"buckets"`
	}

	// struct which contain the instance list for a product
	InstancePerProduct struct {
		Product   string
		Instances []utils.CostPerInstance
	}

	// type that define the parsed response of ES
	Response []InstancePerProduct
)

// getHistoryDate return the begin and the end date of the last month
func getHistoryDate() (time.Time, time.Time) {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month() - 1, 1, 0, 0, 0, 0, now.Location()).UTC()
	end := time.Date(now.Year(), now.Month(), 0, 23, 59, 59, 999999999, now.Location()).UTC()
	return start, end
}

// makeElasticSearchRequestForCost will make the actual request to the ElasticSearch
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequestForCost(ctx context.Context, client *elastic.Client, aa aws.AwsAccount,
	startDate time.Time, endDate time.Time) (*elastic.SearchResult, int, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUserId(aa.UserId, es.IndexPrefixLineItems)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("usageAccountId", es.GetAccountIdFromRoleArn(aa.RoleArn)))
	query = query.Filter(elastic.NewTermsQuery("productCode", "AmazonEC2", "AmazonCloudWatch", "AmazonRDS"))
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(startDate).To(endDate))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("products", elastic.NewTermsAggregation().Field("productCode").Size(0x7FFFFFFF).
		SubAggregation("instances", elastic.NewTermsAggregation().Field("resourceId").Size(0x7FFFFFFF).
			SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost"))))
	result, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			logger.Warning("Query execution failed, ES index does not exists", map[string]interface{}{"index": index, "error": err.Error()})
			return nil, http.StatusOK, err
		}
		logger.Error("Query execution failed", err.Error())
		return nil, http.StatusInternalServerError, nil
	}
	return result, http.StatusOK, nil
}

// getCostPerInstance returns the parsed result of ES
// This response contains the list of the instances of products and the cost associated
func getCostPerInstance(ctx context.Context, aa aws.AwsAccount, startDate time.Time, endDate time.Time) (Response, error) {
	var parsedResult EsInstancePerProductResult
	var response     Response
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	result, returnCode, err := makeElasticSearchRequestForCost(ctx, es.Client, aa, startDate, endDate)
	if err != nil {
		if returnCode != http.StatusOK {
			return nil, err
		} else {
			return nil, nil
		}
	}
	err = json.Unmarshal(*result.Aggregations["products"], &parsedResult.Products)
	if err != nil {
		logger.Error("Error while unmarshaling", err)
		return nil, errors.New("Internal server error")
	}
	for _, product := range parsedResult.Products.Buckets {
		res := InstancePerProduct{product.Product, []utils.CostPerInstance{}}
		for _, instance := range product.Instances.Buckets {
			res.Instances = append(res.Instances, utils.CostPerInstance{instance.Instance, instance.Cost.Value})
		}
		response = append(response, res)
	}
	return response, nil
}

// checkBillingDataCompleted checks if billing data in ES are complete.
// If they are complete it returns true, otherwise it returns false.
func checkBillingDataCompleted(ctx context.Context, startDate time.Time, endDate time.Time, aa aws.AwsAccount) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("usageAccountId", es.GetAccountIdFromRoleArn(aa.RoleArn)))
	query = query.Filter(elastic.NewTermQuery("invoiceId", ""))
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(startDate).To(endDate))
	index := es.IndexNameForUserId(aa.UserId, es.IndexPrefixLineItems)
	result, err := es.Client.Search().Index(index).Size(1).Query(query).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			logger.Warning("Query execution failed, ES index does not exists", map[string]interface{}{"index": index, "error": err.Error()})
			return false, nil
		}
		logger.Error("Query execution failed", err.Error())
		return false, err
	}
	if result.Hits.TotalHits == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// getInstanceInfo sort products and call history reports (only ec2 for now)
func getInstancesInfo(ctx context.Context, aa aws.AwsAccount, response Response, startDate time.Time, endDate time.Time) (error) {
	var ec2Cost, cloudWatchCost, rdsCost []utils.CostPerInstance
	var stringError = ""
	for _, product := range response {
		switch product.Product {
		case "AmazonEC2":
			ec2Cost = product.Instances
			break
		case "AmazonCloudWatch":
			cloudWatchCost = product.Instances
			break
		case "AmazonRDS":
			rdsCost = product.Instances
			break
		}
	}
	errEc2 := ec2.GetEc2MonthlyReport(ctx, ec2Cost, cloudWatchCost, aa, startDate, endDate)
	if errEc2 != nil {
		stringError += errEc2.Error()
	}
	errRds := rds.GetRdsMonthlyReport(ctx, rdsCost, aa, startDate, endDate)
	if errRds != nil {
		stringError += " + " + errRds.Error()
	}
	if stringError != "" {
		if len(stringError) > 254 {
			stringError = stringError[0:254]
		}
		return errors.New(stringError)
	}
	return nil
}

// FetchHistoryInfos fetches billing data and stats of EC2 and RDS instances of the last month
func FetchHistoryInfos(ctx context.Context, aa aws.AwsAccount) (error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	startDate, endDate := getHistoryDate()
	logger.Info("Starting history report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	if complete, err := checkBillingDataCompleted(ctx, startDate, endDate, aa); !complete || err != nil {
		logger.Info("Billing data are not completed", err)
		return err
	}
	response, err := getCostPerInstance(ctx, aa, startDate, endDate)
	if err != nil {
		return err
	}
	return getInstancesInfo(ctx, aa, response, startDate, endDate)
}
