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
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
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
								Key string `json:"key_as_string"`
								Amount struct {
									Buckets []struct {
										Key float64 `json:"key"`
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

func getInstanceCountHours(ctx context.Context, res ResponseInstanceCountMonthly, idxRegion, idxType int) []InstanceCountHours {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	hours := make([]InstanceCountHours, 0)
	for _, date := range res.Region.Buckets[idxRegion].Type.Buckets[idxType].Date.Buckets {
		hour, err := time.Parse("2006-01-02T15:04:05.000Z", date.Key)
		if err != nil {
			logger.Error("Failed to Parse date", err.Error())
		}
		totalAmount := 0.0
		for _, amount := range date.Amount.Buckets {
			totalAmount += amount.Key
		}
		hours = append(hours, InstanceCountHours{
			Hour:  hour,
			Count: totalAmount,
		})
	}
	return hours
}

func FormatResultInstanceCount(ctx context.Context, res *elastic.SearchResult, aa taws.AwsAccount, startDate time.Time) []InstanceCountReport {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseInstanceCountMonthly
	err := json.Unmarshal(*res.Aggregations["region"], &response.Region)
	if err != nil {
		logger.Error("Failed to parse JSON Instance Count document.", err.Error())
	}
	reports := make([]InstanceCountReport, 0)
	for regionIdx, region := range response.Region.Buckets {
		for typeIdx, usageType := range region.Type.Buckets {
			hours := getInstanceCountHours(ctx, response, regionIdx, typeIdx)
			reports = append(reports, InstanceCountReport{
				ReportBase: utils.ReportBase{
					Account: aa.AwsIdentity,
					ReportDate: startDate,
					ReportType: "monthly",
				},
				InstanceCount: InstanceCount{
					Type:   usageType.Key,
					Region: region.Key,
					Hours:  hours,
				},
			})
		}
	}
	return reports
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
	search := GetElasticSearchParams(parsedParams.AccountList, startDate, endDate, es.Client, index)
	res, err := search.Do(ctx)
	if err != nil {
		logger.Error("Error when doing the search", err)
	}
	reports := FormatResultInstanceCount(ctx, res, aa, startDate)
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
		logger.Info("There is already an Instance Count monthly report", nil)
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
