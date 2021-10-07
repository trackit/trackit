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

package plugins_account_core

import (
	"github.com/olivere/elastic"
)

const maxAggregationSize = 0x7FFFFFFF

// createQueryAccountFilterPlugins creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterPlugins(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// GetElasticSearchPluginsParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- accountList []string : A slice of strings representing aws account ids
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query.
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchPluginsParams(accountList []string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(accountList) > 0 {
		query = query.Filter(createQueryAccountFilterPlugins(accountList))
	}
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("top_plugins_account", elastic.NewTermsAggregation().Field("accountPluginIdx").Size(maxAggregationSize).
		SubAggregation("top_reports_hits", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(1)))
	return search
}
