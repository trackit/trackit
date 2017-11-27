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
package costs

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

// aggregationBuilder is an alias for the function type that is used in the
// aggregationBuilder functions.
type aggregationBuilder func([]string) []paramAggrAndName

// paramNameToFuncPtr maps parameter names to functions building the aggregations.
// map of string keys and functions pointer as values. For each possible param
// after parsing (removing the ':<TAG_KEY>' in the case of the tag), there is a function associated to it
// that create the Aggregations needed by this param.
// If a new param, that is only creating aggregations, needs to be added,
// a functions with an aggregationBuilder prototype need to be added to the list below.
var paramNameToFuncPtr = map[string]aggregationBuilder{
	"product": createAggregationPerProduct,
	"region":  createAggregationPerRegion,
	"account": createAggregationPerAccount,
	"tag":     createAggregationPerTag,
	"cost":    createCostSumAggregation,
	"day":     createAggregationPerDay,
	"week":    createAggregationPerWeek,
	"month":   createAggregationPerMonth,
	"year":    createAggregationPerYear,
}

// paramAggrAndName is a structure containing the name of the parameter and
// corresponding aggregation. A parameter is a string that is passed to the
// GetElasticSearchParams and represents an aggregation. A list of those
// parameters can be found in the paramNameToFuncPtr map.
type paramAggrAndName struct {
	name string
	aggr elastic.Aggregation
}

// aggregationMaxSize is the maximum size of an Elastic Search Aggregation
const aggregationMaxSize = 0x7FFFFFFF

// createQueryAccountFilter creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilter(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("linked_account_id", accountListFormatted...)
}

// createQueryTimeRange creates and return a new *elastic.RangeQuery based on the duration
// defined by durationBegin and durationEnd
func createQueryTimeRange(durationBegin time.Time, durationEnd time.Time) *elastic.RangeQuery {
	return elastic.NewRangeQuery("usage_start_date").
		From(durationBegin).To(durationEnd)
}

// createAggregationPerProduct creates and returns a new []paramAggrAndName of size 1 which creates a
// bucket aggregation on the field 'product_name'
func createAggregationPerProduct(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "product",
			aggr: elastic.NewTermsAggregation().
				Field("product_name").Size(aggregationMaxSize),
		},
	}
}

// createAggregationPerRegion creates and returns a new []paramAggrAndName of size 1 which creates a
// bucket aggregation on the field 'availability_zone'
func createAggregationPerRegion(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "region",
			aggr: elastic.NewTermsAggregation().
				Field("availability_zone").Size(aggregationMaxSize),
		},
	}
}

// createAggregationPerAccount creates and returns a new []paramAggrAndName of size 1 which creates a
// bucket aggregation on the field 'linked_account_id'
func createAggregationPerAccount(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "account",
			aggr: elastic.NewTermsAggregation().
				Field("linked_account_id").Size(aggregationMaxSize),
		},
	}
}

// createAggregationPerDay creates and returns a new []paramAggrAndName of size 1 which creates a
// date histogram aggregation on the field 'usage_start_date' with a time range of a day
func createAggregationPerDay(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "day",
			aggr: elastic.NewDateHistogramAggregation().
				Field("usage_start_date").Interval("day"),
		},
	}
}

// createAggregationPerWeek creates and returns a new []paramAggrAndName of size 1 which creates a
// date histogram aggregation on the field 'usage_start_date' with a time range of a week
func createAggregationPerWeek(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "week",
			aggr: elastic.NewDateHistogramAggregation().
				Field("usage_start_date").Interval("week"),
		},
	}
}

// createAggregationPerMonth creates and returns a new []paramAggrAndName of size 1 which creates a
// date histogram aggregation on the field 'usage_start_date' with a time range of a month
func createAggregationPerMonth(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "month",
			aggr: elastic.NewDateHistogramAggregation().
				Field("usage_start_date").Interval("month"),
		},
	}
}

// createAggregationPerYear creates and returns a new []paramAggrAndName of size 1 which creates a
// date histogram aggregation on the field 'usage_start_date' with a time range of a year
func createAggregationPerYear(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "year",
			aggr: elastic.NewDateHistogramAggregation().
				Field("usage_start_date").Interval("year"),
		},
	}
}

// createAggregationPerTag creates and returns a new []paramAggrAndName of size 2 which consits
// of two aggregations that are required for the tag param.
// The first aggregation is a FilterAggregation on the field 'tag.key', and with a value of
// the tag key passed in the parameter 'paramSplit' in the form "user:<TAG_KEY_VALUE>".
// The second aggregation is a TermsAggregation that creates bucket aggregation on the field
// 'tag.value'.
// No SubAggregation is created in this function, as it needs to be created in the nestAggregation function
func createAggregationPerTag(paramSplit []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "tag_key",
			aggr: elastic.NewFilterAggregation().
				Filter(elastic.NewTermQuery("tag.key", fmt.Sprintf("user:%v", paramSplit[1])))},
		paramAggrAndName{
			name: "tag_value",
			aggr: elastic.NewTermsAggregation().
				Field("tag.value").Size(aggregationMaxSize)},
	}
}

// createCostSumAggregation : Creates and return a new []paramAggrAndName of size 1, which creates a
// SumAggregation on the field 'cost'
func createCostSumAggregation(_ []string) []paramAggrAndName {
	return []paramAggrAndName{
		paramAggrAndName{
			name: "cost",
			aggr: elastic.NewSumAggregation().Field("cost"),
		},
	}
}

// reverseAggregationArray : reverse the paramAggrAndName slice that is passed to it
func reverseAggregationArray(aggregationArray []paramAggrAndName) []paramAggrAndName {
	for i := len(aggregationArray)/2 - 1; i >= 0; i-- {
		opp := len(aggregationArray) - 1 - i
		aggregationArray[i], aggregationArray[opp] = aggregationArray[opp], aggregationArray[i]
	}
	return aggregationArray
}

// nestAggregation takes a slice of paramAggrAndName type, and will nest the different aggregations.
// Aggregations are nested by creating a chain of SubAggregation
// A type switch is required to simulate downcasting from the interface elastic.Aggregation.
// Current types on the type switch are TermsAggregation, FilterAggregation, SumAggregation and
// DateHistogramAggregation.
// If a new function creating a type that is not listed here is added to the paramNameToFuncPtr map
// it should be added to the type switch, or the function will create bugged SubAggregations
func nestAggregation(allAggrSlice []paramAggrAndName) elastic.Aggregation {
	allAggrSlice = reverseAggregationArray(allAggrSlice)
	aggrToNest := allAggrSlice[0]
	for _, baseAggr := range allAggrSlice[1:] {
		switch assertedBaseAggr := baseAggr.aggr.(type) {
		case *elastic.TermsAggregation:
			aggrBuff := assertedBaseAggr.SubAggregation(aggrToNest.name, aggrToNest.aggr)
			aggrToNest = paramAggrAndName{name: baseAggr.name, aggr: aggrBuff}
		case *elastic.FilterAggregation:
			aggrBuff := assertedBaseAggr.SubAggregation(aggrToNest.name, aggrToNest.aggr)
			aggrToNest = paramAggrAndName{name: baseAggr.name, aggr: aggrBuff}
		case *elastic.DateHistogramAggregation:
			aggrBuff := assertedBaseAggr.SubAggregation(aggrToNest.name, aggrToNest.aggr)
			aggrToNest = paramAggrAndName{name: baseAggr.name, aggr: aggrBuff}
		}
	}
	return aggrToNest.aggr
}

// GetElasticSearchParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES.
// It takes as paramters :
// 	- accountList []string : A slice of strings representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//	- durationBeing time.Time : A time.Time struct representing the begining of the time range in the query
//	- durationEnd time.Time : A time.Time struct representing the end of the time range in the query
//	- param []string : A slice of strings representing the different parameters, in the nesting order,
//	that will create aggregations.
//	Those can be :
//		- "product" : It will create a TermsAggregation on the field 'product_name'
//		- "region" : It will create a TermsAggregation on the field 'availability_zone'
//		- "account" : It will create a TermsAggregation on the field 'linked_account_id'
//		- "tag:<TAG_KEY>" : It will create a FilterAggregation on the field 'tag.key',
//		filtering on the value 'user:<TAG_KEY>'.
//		It will then create a TermsAggregation on the field 'tag.value'
//		- "[day|week|month|year]": It will create a DateHistogramAggregation on the specified duration on
//		the field 'usage_start_date'
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "awsdetailedlineitems"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- For the 'tag:<TAG_KEY>' param, if the separator is not present, or if there is no key that is passed to it,
//	the program will crash
//	- If a param in the slice is not present in the detailedLineItemsFieldsName, the program will crash.
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchParams(accountList []string, durationBegin time.Time,
	durationEnd time.Time, params []string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	query = query.Filter(createQueryAccountFilter(accountList))
	query = query.Filter(createQueryTimeRange(durationBegin, durationEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	params = append(params, "cost")
	var allAggregationSlice []paramAggrAndName
	for _, paramName := range params {
		paramNameSplit := strings.Split(paramName, ":")
		paramAggr := paramNameToFuncPtr[paramNameSplit[0]](paramNameSplit)
		allAggregationSlice = append(allAggregationSlice, paramAggr...)
	}
	aggregationParamName := allAggregationSlice[0].name
	nestedAggregation := nestAggregation(allAggregationSlice)
	search.Aggregation(aggregationParamName, nestedAggregation)
	return search
}
