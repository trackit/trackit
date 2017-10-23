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

package costs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/sha1sum/aws_signing_client"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"
)

var detailedLineItemFieldsName = map[string]func([]string) []paramAggrAndName{
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

type paramAggrAndName struct {
	paramName string
	paramAggr elastic.Aggregation
}

const aggregationMaxSize = 0x7FFFFFFF

func createQueryAccountFilter(accountList []string) *elastic.TermsQuery {
	return elastic.NewTermsQuery("linked_account_id", accountList)
}

func createQueryTimeRange(durationBegin time.Time, durationEnd time.Time) *elastic.RangeQuery {
	return elastic.NewRangeQuery("usage_start_date").From(durationBegin).To(durationEnd)
}

func createAggregationPerProduct(paramSplit []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "product",
		paramAggr: elastic.NewTermsAggregation().Field("product_name").Size(aggregationMaxSize)}
	return res
}

func createAggregationPerRegion(paramSplit []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "region",
		paramAggr: elastic.NewTermsAggregation().Field("availability_zone").Size(aggregationMaxSize)}
	return res
}

func createAggregationPerAccount(paramSplit []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "account",
		paramAggr: elastic.NewTermsAggregation().Field("linked_account_id").Size(aggregationMaxSize)}
	return res
}

func createAggregationPerDay(paramList []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "day",
		paramAggr: elastic.NewDateHistogramAggregation().Field("usage_start_date").Interval("day")}
	return res
}

func createAggregationPerWeek(paramList []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "week",
		paramAggr: elastic.NewDateHistogramAggregation().Field("usage_start_date").Interval("week")}
	return res
}

func createAggregationPerMonth(paramList []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "month",
		paramAggr: elastic.NewDateHistogramAggregation().Field("usage_start_date").Interval("month")}
	return res
}

func createAggregationPerYear(paramList []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "year",
		paramAggr: elastic.NewDateHistogramAggregation().Field("usage_start_date").Interval("year")}
	return res
}

func createAggregationPerTag(paramSplit []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 2)
	res[0] = paramAggrAndName{
		paramName: "tag_key",
		paramAggr: elastic.NewFilterAggregation().Filter(elastic.NewTermQuery("tag.key", fmt.Sprintf("user:%v", paramSplit[1])))}
	res[1] = paramAggrAndName{
		paramName: "tag_value",
		paramAggr: elastic.NewTermsAggregation().Field("tag.value").Size(aggregationMaxSize)}
	return res
}

func createCostSumAggregation(paramSplit []string) []paramAggrAndName {
	res := make([]paramAggrAndName, 1)
	res[0] = paramAggrAndName{
		paramName: "cost",
		paramAggr: elastic.NewSumAggregation().Field("cost")}
	return res
}

func reverseAggregationArray(aggregationArray []paramAggrAndName) []paramAggrAndName {
	for i := len(aggregationArray)/2 - 1; i >= 0; i-- {
		opp := len(aggregationArray) - 1 - i
		aggregationArray[i], aggregationArray[opp] = aggregationArray[opp], aggregationArray[i]
	}
	return aggregationArray
}

func nestAggregation(allAggrSlice []paramAggrAndName) elastic.Aggregation {
	allAggrSlice = reverseAggregationArray(allAggrSlice)
	aggrToNest := allAggrSlice[0]
	for _, baseAggr := range allAggrSlice[1:] {
		fmt.Printf("aggrToNest.paramName = %v; baseAggr.paramName = %v\n", aggrToNest.paramName, baseAggr.paramName)
		switch assertedBaseAggr := baseAggr.paramAggr.(type) {
		case *elastic.TermsAggregation:
			aggrBuff := assertedBaseAggr.SubAggregation(aggrToNest.paramName, aggrToNest.paramAggr)
			aggrToNest = paramAggrAndName{paramName: baseAggr.paramName, paramAggr: aggrBuff}
		case *elastic.FilterAggregation:
			aggrBuff := assertedBaseAggr.SubAggregation(aggrToNest.paramName, aggrToNest.paramAggr)
			aggrToNest = paramAggrAndName{paramName: baseAggr.paramName, paramAggr: aggrBuff}
		case *elastic.SumAggregation:
			aggrBuff := assertedBaseAggr.SubAggregation(aggrToNest.paramName, aggrToNest.paramAggr)
			aggrToNest = paramAggrAndName{paramName: baseAggr.paramName, paramAggr: aggrBuff}
		case *elastic.DateHistogramAggregation:
			aggrBuff := assertedBaseAggr.SubAggregation(aggrToNest.paramName, aggrToNest.paramAggr)
			aggrToNest = paramAggrAndName{paramName: baseAggr.paramName, paramAggr: aggrBuff}
		}
	}
	return aggrToNest.paramAggr
}

// GetElasticSearchParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as paramters :
// 	- accountList []string : A slice of strings representing aws account number, in the format of the field 'awsdetailedlineitem.linked_account_id'
//	- durationBeing time.TIme : A time.Time struct representing the begining of the time range in the query
//	- durationEnd time.Time : A time.Time struct representing the end of the time range in the query
//	- param []string : A slice of strings representing the different parameters, in the nesting order, that will create aggregations.
//	  Those can be :
//		- "product" : It will create a TermsAggregation on the field 'product_name'
//		- "region" : It will create a TermsAggregation on the field 'availability_zone'
//		- "account" : It will create a TermsAggregation on the field 'linked_account_id'
//		- "tag:<TAG_KEY>" : It will create a FilterAggregation on the field 'tag.key', filtering on the value 'user:<TAG_KEY>'. It will then create a TermsAggregation on the field 'tag.value'
//		- "[day|week|month|year]": It will create a DateHistogramAggregation on the specified duration on the field 'usage_start_date'
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client. It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value should be "awsdetailedlineitems"
func GetElasticSearchParams(accountList []string, durationBegin time.Time, durationEnd time.Time, params []string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(accountList) > 0 {
		query = query.Filter(createQueryAccountFilter(accountList))
	}
	query = query.Filter(createQueryTimeRange(durationBegin, durationEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	params = append(params, "cost")
	var allAggregationSlice []paramAggrAndName
	for _, paramName := range params {
		paramNameSplit := strings.Split(paramName, ":")
		// fmt.Printf("param = %v, paramNameSplit = %v\n", paramName, paramNameSplit)
		paramAggr := detailedLineItemFieldsName[paramNameSplit[0]](paramNameSplit)
		allAggregationSlice = append(allAggregationSlice, paramAggr...)
	}
	aggregationParamName := allAggregationSlice[0].paramName
	nestedAggregation := nestAggregation(allAggregationSlice)
	search.Aggregation(aggregationParamName, nestedAggregation)
	return search
}

func deserialize(searchResult *elastic.SearchResult) {
	if searchResult == nil || searchResult.TotalHits() == 0 {
		fmt.Println("Well that's not suppose to happen")
		return
	}
	fmt.Printf("Total hits : %v\n", searchResult.TotalHits())
	for i, hit := range searchResult.Hits.Hits {
		fmt.Printf("i = %v, hit = %v\n", i, hit)
	}
}

func PlaceholderFunc() {
	credentials := credentials.NewSharedCredentials("", "default")
	signer := v4.NewSigner(credentials)
	awsClient, err := aws_signing_client.New(signer, nil, "es", "us-west-2")
	if err != nil {
		fmt.Printf("error received : %v\n", err)
	}
	client, err := elastic.NewClient(
		elastic.SetURL("https://search-job-msol-prod-trackit-j6ofkgxgmxezkamcywpmqwfsn4.us-west-2.es.amazonaws.com"),
		elastic.SetScheme("https"),
		elastic.SetHttpClient(awsClient),
		elastic.SetSniff(false),
	)
	var accountList = []string{"298868543803", "135736586752"}
	// var accountList = []string{}
	var paramList = []string{"account", "product", "day"}
	// search := GetElasticSearchParams([]string{"394125495069"}, time.Date(2017, 10, 1, 0, 0, 0, 0, time.Local), time.Date(2017, 11, 1, 0, 0, 0, 0, time.Local), []string{"product"}, client, "awsdetailedlineitem")
	search := GetElasticSearchParams(accountList, time.Unix(1504572038, 0), time.Now(), paramList, client, "awsdetailedlineitem")
	a, err := search.Do(context.Background())
	if err != nil {
		panic(err)
	}
	res, err := json.MarshalIndent(*a, "", "  ")
	fmt.Printf("%v %v\n", string(res), err)
	deserialize(a)
	fmt.Println(err)

}

// HandleRequest is a dummy request handler function. It does nothing except
// some logging and returns static data.
func HandleRequest(response http.ResponseWriter, request *http.Request, logger jsonlog.Logger) {
	logger.Debug("Request headers.", request.Header)
	response.WriteHeader(200)
	response.Write([]byte("Costs."))
}
