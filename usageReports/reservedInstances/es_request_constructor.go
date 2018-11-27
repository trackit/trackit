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

package reservedInstances

import (
	"time"

	"gopkg.in/olivere/elastic.v5"
)

const maxAggregationSize = 0x7FFFFFFF

// getDateForDailyReport returns the end and the begin of the date of the report based on a date
// if the date given as parameter is in the actual month, it returns the the the begin of the month et now at midnight
// if the date is before the actual month, it returns the begin and the end of the month given as parameter
func getDateForDailyReport(date time.Time) (begin, end time.Time) {
	now := time.Now().UTC()
	if date.Year() == now.Year() && date.Month() == now.Month() {
		end = now
		begin = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location()).UTC()
		return
	} else {
		begin = date
		end = time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location()).UTC()
		return
	}
}

// createQueryAccountFilterReservedInstances creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterReservedInstances(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// getElasticSearchReservedInstancesDailyParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- params ReservedInstancesQueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "reservedInstances-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getElasticSearchReservedInstancesDailyParams(params ReservedInstancesQueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilterReservedInstances(params.AccountList))
	}
	query = query.Filter(elastic.NewTermQuery("reportType", "daily"))
	dateStart, dateEnd := getDateForDailyReport(params.Date)
	query = query.Filter(elastic.NewRangeQuery("reportDate").
		From(dateStart).To(dateEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts", elastic.NewTermsAggregation().Field("account").
		SubAggregation("dates", elastic.NewTermsAggregation().Field("reportDate").
			SubAggregation("instances", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(maxAggregationSize))))
	return search
}

// getElasticSearchReservedInstancesMonthlyParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- params ReservedInstancesQueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "reservedInstances-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getElasticSearchReservedInstancesMonthlyParams(params ReservedInstancesQueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilterReservedInstances(params.AccountList))
	}
	query = query.Filter(elastic.NewTermQuery("reportType", "monthly"))
	query = query.Filter(elastic.NewTermQuery("reportDate", params.Date))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts", elastic.NewTermsAggregation().Field("account").
		SubAggregation("instances", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(maxAggregationSize)))
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

// getElasticSearchCostParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- params ReservedInstancesQueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getElasticSearchCostParams(params ReservedInstancesQueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilterBill(params.AccountList))
	}
	query = query.Filter(elastic.NewTermsQuery("productCode", "AmazonEC2", "AmazonCloudWatch"))
	dateStart, dateEnd := getDateForDailyReport(params.Date)
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(dateStart).To(dateEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts", elastic.NewTermsAggregation().Field("usageAccountId").Size(maxAggregationSize).
		SubAggregation("instances", elastic.NewTermsAggregation().Field("resourceId").Size(maxAggregationSize).
			SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost"))))
	return search
}
