//   Copyright 2020 MSolution.IO
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

package routes

import (
	"github.com/olivere/elastic"
)

const maxAggregationSize = 0x7FFFFFFF

// getElasticSearchResourcesParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- params ResourcesRequestBody : contains the list of accounts and filters
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "tagging-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getElasticSeachResourcesParams(params ResourcesRequestBody, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.Accounts) > 0 {
		query = query.Filter(createQueryAccountFilterResources(params.Accounts))
	}
	if len(params.Regions) > 0 {
		query = query.Filter(createQueryRegionFilterResources(params.Regions))
	}
	if len(params.ResourceTypes) > 0 {
		query = query.Filter(createQueryTypeFilterResources(params.ResourceTypes))
	}
	if len(params.Tags) > 0 {
		query = query.Filter(createQueryTagsFilterResources(params.Tags))
	}
	if len(params.MissingTags) > 0 {
		query = query.Filter(createQueryMissingTagsFilterResources(params.MissingTags))
	}
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts", elastic.NewTermsAggregation().Field("account").
		SubAggregation("dates", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1).
			SubAggregation("resources", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(maxAggregationSize))))
	return search
}

// createQueryAccountFilterResources creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterResources(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// createQueryRegionFilterResources creates and return a new *elastic.TermsQuery on the regionList array
func createQueryRegionFilterResources(regionList []string) *elastic.TermsQuery {
	regionListFormatted := make([]interface{}, len(regionList))
	for i, v := range regionList {
		regionListFormatted[i] = v
	}
	return elastic.NewTermsQuery("region", regionListFormatted...)
}

//createQueryTypeFilterResources creates and return a new *elastic.TermsQuery on the typeList array
func createQueryTypeFilterResources(typeList []string) *elastic.TermsQuery {
	typeListFormatted := make([]interface{}, len(typeList))
	for i, v := range typeList {
		typeListFormatted[i] = v
	}
	return elastic.NewTermsQuery("resourceType", typeListFormatted...)
}

//createQueryTagsFilterResources creates and return a new *elastic.BoolQuery based on the tagList
func createQueryTagsFilterResources(tagList []Tag) *elastic.BoolQuery {
	termQueries := []elastic.Query{}
	for _, v := range tagList {
		termQueries = append(termQueries, elastic.NewNestedQuery("tags", elastic.NewBoolQuery().Must(elastic.NewTermQuery("tags.key", v.Key), elastic.NewTermQuery("tags.value", v.Value))))
	}
	return elastic.NewBoolQuery().Must(termQueries...)
}

//createQueryTypeFilterResources creates and return a new *elastic.BoolQuery based on the missingTagList
func createQueryMissingTagsFilterResources(missingTagList []Tag) *elastic.BoolQuery {
	termQueries := []elastic.Query{}
	for _, v := range missingTagList {
		termQueries = append(termQueries, elastic.NewNestedQuery("tags", elastic.NewBoolQuery().Must(elastic.NewTermQuery("tags.key", v.Key), elastic.NewTermQuery("tags.value", v.Value))))
	}
	return elastic.NewBoolQuery().MustNot(termQueries...)
}
