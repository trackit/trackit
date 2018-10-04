//   Copyright 2017 MSolution.IO
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

// Package costs gets billing information from an ElasticSearch.
package rds

import (
	"time"

	"gopkg.in/olivere/elastic.v5"
)

// getDateForDailyReport returns the end and the begin of the date of the report based on a date
func getDateForDailyReport(date time.Time) (begin, end time.Time) {
	now := time.Now().UTC()
	if date.Year() == now.Year() && date.Month() == now.Month() {
		end = now
		begin = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location()).UTC()
		return
	} else {
		begin = date
		end = time.Date(date.Year(), date.Month() + 1, 0, 23, 59, 59, 999999999, date.Location()).UTC()
		return
	}
}

// createQueryAccountFilterRds creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterRds(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// GetElasticSearchRdsDailyParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as paramters :
// 	- params rdsQueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "rds-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchRdsDailyParams(params rdsQueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.accountList) > 0 {
		query = query.Filter(createQueryAccountFilterRds(params.accountList))
	}
	query = query.Filter(elastic.NewTermQuery("reportType", "daily"))
	dateStart, dateEnd := getDateForDailyReport(params.date)
	query = query.Filter(elastic.NewRangeQuery("reportDate").
		From(dateStart).To(dateEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("top_reports", elastic.NewTermsAggregation().Field("account").
		SubAggregation("top_reports_hits", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(1)))
	return search
}

// GetElasticSearchRdsMonthlyParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- params rdsQueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "rds-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchRdsMonthlyParams(params rdsQueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.accountList) > 0 {
		query = query.Filter(createQueryAccountFilterRds(params.accountList))
	}
	query = query.Filter(elastic.NewTermQuery("reportType", "monthly"))
	query = query.Filter(elastic.NewTermQuery("reportDate", params.date))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("top_reports", elastic.NewTermsAggregation().Field("account").
		SubAggregation("top_reports_hits", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(1)))
	return search
}

// createQueryAccountFilterBill creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterBill(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("usageAccountId", accountListFormatted...)
}

// GetElasticSearchCostParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as paramters :
// 	- params rdsQueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "rds-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchCostParams(params rdsQueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.accountList) > 0 {
		query = query.Filter(createQueryAccountFilterBill(params.accountList))
	}
	query = query.Filter(elastic.NewTermQuery("productCode", "AmazonRDS"))
	dateEnd := time.Now().UTC()
	dateStart, dateEnd := getDateForDailyReport(params.date)
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(dateStart).To(dateEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts",  elastic.NewTermsAggregation().Field("usageAccountId").Size(len(params.accountList)).
		SubAggregation("instances", elastic.NewTermsAggregation().Field("resourceId").Size(0x7FFFFFFF).
			SubAggregation("cost",  elastic.NewSumAggregation().Field("unblendedCost"))))
	return search
}
