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
package ec2

import (
	"gopkg.in/olivere/elastic.v5"
	"time"
)

// createQueryAccountFilter creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilter(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// GetElasticSearchEc2Params is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as paramters :
// 	- accountList []string : A slice of strings representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "ec2-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchEc2Params(accountList []string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(accountList) > 0 {
		query = query.Filter(createQueryAccountFilter(accountList))
	}
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("top_reports", elastic.NewTermsAggregation().Field("account").
		SubAggregation("top_reports_hits", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(1)))
	return search
}

// GetElasticSearchCostParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as paramters :
// 	- accountList []string : A slice of strings representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "ec2-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchCostParams(account string, date string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("usageAccountId", account))
	query = query.Filter(elastic.NewTermQuery("productCode", "AmazonEC2"))
	dateEnd, err := time.Parse("2006-01-02T15:04:05.000Z", date)
	if err == nil {
		dateStart := time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day()-31, 0, 0, 0, 0, dateEnd.Location())
		query = query.Filter(elastic.NewRangeQuery("usageStartDate").
			From(dateStart).To(dateEnd))
	}
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("instances", elastic.NewTermsAggregation().Field("resourceId").
		SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost")))
	return search
}