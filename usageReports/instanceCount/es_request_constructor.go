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

package instanceCount

import (
	"github.com/olivere/elastic"
)

const maxAggregationSize = 0x7FFFFFFF

// createQueryAccountFilterInstanceCount creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterInstanceCount(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// getElasticSearchInstanceCountMonthlyParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- params InstanceCountQueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an Instance Count of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "instanceCount-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getElasticSearchInstanceCountMonthlyParams(params InstanceCountQueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilterInstanceCount(params.AccountList))
	}
	query = query.Filter(elastic.NewTermQuery("reportType", "monthly"))
	query = query.Filter(elastic.NewTermQuery("reportDate", params.Date))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts", elastic.NewTermsAggregation().Field("account").
		SubAggregation("reports", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(maxAggregationSize)))
	return search
}
