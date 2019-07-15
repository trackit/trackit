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

package instanceCount

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/s3"
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

	ResponseInstanceCountMonthly struct {
		Region struct {
			Buckets []struct {
				Key  string `json:"key"`
				Type struct {
					Buckets []struct {
						Key  string `json:"key"`
						Date struct {
							Buckets []struct {
								Key    string `json:"key_as_string"`
								Amount struct {
									Buckets []struct {
										Key string `json:"key_as_string"`
									} `json:"buckets"`
								} `json:"amount"`
							} `json:"buckets"`
						} `json:"date"`
					} `json:"buckets"`
				} `json:"type"`
			} `json:"buckets"`
		} `json:"region"`
	}
)

// getElasticSearchInstanceCount prepares and run the request to retrieve the a report of an instance count
// It will return the data and an error.
func getElasticSearchInstanceCount(ctx context.Context, account, report string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("report.id", report))
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

// getInstanceCountInfoFromEs gets information about an instance count from previous report to put it in the new report
func getInstanceCountInfoFromES(ctx context.Context, report utils.CostPerResource, account string, userId int) InstanceCount {
	var docType InstanceCountReport
	var inst = InstanceCount{
		// dont let it empty
		Type: "",
		Hours: []InstanceCountHours{},
	}
	res, err := getElasticSearchInstanceCount(ctx, account, report.Resource,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixInstanceCountReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			inst.Type = docType.InstanceCount.Type
			inst.Hours = docType.InstanceCount.Hours
		}
	}
	return inst
}
/*
// fetchMonthlyInstanceCountList sends in reportInfoChan the instance count fetched from ES
func fetchMonthlyInstanceCountList(ctx context.Context, creds *credentials.Credentials, inst utils.CostPerResource,
	account, region string, reportChan chan InstanceCount, userId int) error {
	defer close(reportChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	reports, err := svc.DescribeSnapshots(&desc)
	if err != nil {
		reportChan <- getInstanceCountInfoFromES(ctx, inst, account, userId)
		return err
	}
	for _, report := range reports.Snapshots {
		reportChan <- InstanceCount{
			Tags: getSnapshotTag(report.Tags),
			Cost: snap.Cost,
			Volume: Volume{
				Id:   aws.StringValue(report.VolumeId),
				Size: aws.Int64Value(report.VolumeSize),
			},
		}
	}
	return nil
}
*/

func FormatResultInstanceCount(ctx context.Context, res *elastic.SearchResult, parsedParams EsQueryParams) []InstanceCountReport {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("AGGREGATIONS ======", map[string]interface{}{
		"INDEX": "region",
		"DATE": res.Aggregations["region"],
	})
	var response ResponseInstanceCountMonthly
	err := json.Unmarshal(*res.Aggregations["region"], &response.Region)
	if err != nil {
		logger.Error("Failed to parse JSON Instance Count document.", err.Error())
	}
	logger.Debug("VALUE =======", map[string]interface{}{
		"DATA": response.Region,
	})
	/*for _, region := range response.Region.Buckets {
		for _, usageType := range region.Type.Buckets {
			for _, date := range usageType.Date.Buckets {
				logger.Debug("REGION AND USAGE TYPE ======", map[string]interface{}{
					"REGION": region.Key,
					"TYPE":   usageType.Key,
					"DATE":   date.Key,
				})
			}
		}
	}*/
	return []InstanceCountReport{}
}

// getInstanceCountMetrics gets credentials, accounts and region to fetch InstanceCount report stats
func fetchMonthlyInstanceCountReports(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) ([]InstanceCountReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUserId(aa.UserId, s3.IndexPrefixLineItem)
	parsedParams := EsQueryParams{
		AccountList: []string{aa.AwsIdentity},
		IndexList:   []string{index},
		DateBegin:   startDate,
		DateEnd:     endDate,
	}

	// Préparer la requête pour l'ES avec les queryParams nécessaires
	search := GetElasticSearchParams(parsedParams.AccountList, startDate, endDate, es.Client, index)
	res, err := search.Do(ctx)
	if err != nil {
		logger.Error("Error when doing the search", err)
	}
	logger.Debug("HERE IT IS RESULT FROM SEARCH ON ES", map[string]interface{}{
		"account": parsedParams.AccountList,
		"res": res,
	})
	reports := FormatResultInstanceCount(ctx, res, parsedParams)
	// Trier les données et les formater correctement
	// Générer les InstanceCountReport avec les données formatées dedans

	return reports, nil
}

// PutInstanceCountMonthlyReport puts a monthly report of InstanceCount in ES
func PutInstanceCountMonthlyReport(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting InstanceCount monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixInstanceCountReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an InstanceCount monthly report", nil)
		return false, nil
	}
	reports, err := fetchMonthlyInstanceCountReports(ctx, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importInstanceCountToEs(ctx, aa, reports)
	if err != nil {
		return false, err
	}
	return true, nil
}
