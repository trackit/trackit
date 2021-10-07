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
	"testing"
	"time"

	"github.com/olivere/elastic"
)

// Fonction not used
/* func createAndConfigureTestClient(t *testing.T) *elastic.Client {
	client, err := elastic.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	return client
} */

func TestQueryAccountFiltersMultipleAccounts(t *testing.T) {
	linkedAccountID := []string{
		"123456",
		"98765432",
	}
	expectedResult := `{"terms":{"usageAccountId":["123456","98765432"]}}`
	res := createQueryAccountFilter(linkedAccountID)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestQueryAccountFiltersSingleAccount(t *testing.T) {
	linkedAccountID := []string{
		"123456",
	}
	expectedResult := `{"terms":{"usageAccountId":["123456"]}}`
	res := createQueryAccountFilter(linkedAccountID)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestQueryTimeRange(t *testing.T) {
	durationBegin, _ := time.Parse("2006-1-2 15:04", "2017-01-12 11:23")
	durationEnd, _ := time.Parse("2006-1-2 15:04", "2017-05-23 11:23")
	expectedResult := `{"range":{"usageStartDate":{"from":"2017-01-12T11:23:00Z","include_lower":true,"include_upper":true,"to":"2017-05-23T11:23:00Z"}}}`

	res := createQueryTimeRange(durationBegin, durationEnd)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerProduct(t *testing.T) {
	res := createAggregationPerProduct([]string{""})
	expectedResult := `{"terms":{"field":"productCode","size":2147483647}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerRegion(t *testing.T) {
	res := createAggregationPerRegion([]string{""})
	expectedResult := `{"terms":{"field":"region","size":2147483647}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerAccount(t *testing.T) {
	res := createAggregationPerAccount([]string{""})
	expectedResult := `{"terms":{"field":"usageAccountId","size":2147483647}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerDay(t *testing.T) {
	res := createAggregationPerDay([]string{""})
	expectedResult := `{"date_histogram":{"field":"usageStartDate","interval":"day","min_doc_count":0}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerMonth(t *testing.T) {
	res := createAggregationPerMonth([]string{""})
	expectedResult := `{"date_histogram":{"field":"usageStartDate","interval":"month","min_doc_count":0}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestCostSumAggregation(t *testing.T) {
	res := createCostSumAggregation([]string{""})
	expectedResult := `{"sum":{"field":"unblendedCost"}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerWeek(t *testing.T) {
	res := createAggregationPerWeek([]string{""})
	expectedResult := `{"date_histogram":{"field":"usageStartDate","interval":"week","min_doc_count":0}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerYear(t *testing.T) {
	res := createAggregationPerYear([]string{""})
	expectedResult := `{"date_histogram":{"field":"usageStartDate","interval":"year","min_doc_count":0}}`
	src, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestAggregationPerTag(t *testing.T) {
	res := createAggregationPerTag([]string{"tag", "test"})
	expectedFirstResult := `{"filter":{"term":{"tag.key":"user:test"}}}`
	expectedSecondResult := `{"terms":{"field":"tag.value","size":2147483647}}`
	srcFirst, err := res[0].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	srcSecond, err := res[1].aggr.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonResFirst, err := json.Marshal(srcFirst)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonResFirst) != expectedFirstResult {
		t.Fatalf("Expected %v but got %v", expectedFirstResult, string(jsonResFirst))
	}
	jsonResSecond, err := json.Marshal(srcSecond)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonResSecond) != expectedSecondResult {
		t.Fatalf("Expected %v but got %v", expectedSecondResult, string(jsonResSecond))
	}
}

func TestReverseAggregationArray(t *testing.T) {
	sliceTobeReversed := []paramAggrAndName{
		{
			name: "first",
			aggr: nil},
		{
			name: "second",
			aggr: nil},
		{
			name: "thrice",
			aggr: nil}}
	expectedReversedSlice := []paramAggrAndName{
		{
			name: "thrice",
			aggr: nil},
		{
			name: "second",
			aggr: nil},
		{
			name: "first",
			aggr: nil}}
	reversedSlice := reverseAggregationArray(sliceTobeReversed)
	for i, val := range reversedSlice {
		if val.name != expectedReversedSlice[i].name {
			t.Fatalf("Expected %v on index %v but got %v", expectedReversedSlice[i].name, i, val.name)
		}
	}
}

func TestAggregationNestingWithSingleElementSlice(t *testing.T) {
	singleAggregationSlice := createAggregationPerAccount([]string{""})
	expectedResult := `{"terms":{"field":"usageAccountId","size":2147483647}}`
	res := nestAggregation(singleAggregationSlice)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonResult, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonResult) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonResult))
	}
}

func TestAggregationNestingWithCoupleElementsSlice(t *testing.T) {
	coupleAggregationSlice := createAggregationPerTag([]string{"", "test"})
	expectedResult := `{
	"aggregations": {
		"tag_value": {
			"terms": {
				"field": "tag.value",
				"size": 2147483647
			}
		}
	},
	"filter": {
		"term": {
			"tag.key": "user:test"
		}
	}
}`
	res := nestAggregation(coupleAggregationSlice)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonResult, err := json.MarshalIndent(src, "", "	")
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonResult) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonResult))
	}
}

func TestAggregationNestingWithFewElementsSlice(t *testing.T) {
	fewAggregationSlice := createAggregationPerTag([]string{"", "test"})
	buffAggregation := createCostSumAggregation([]string{""})
	fewAggregationSlice = append(fewAggregationSlice, buffAggregation...)
	expectedResult := `{
	"aggregations": {
		"tag_value": {
			"aggregations": {
				"value": {
					"sum": {
						"field": "unblendedCost"
					}
				}
			},
			"terms": {
				"field": "tag.value",
				"size": 2147483647
			}
		}
	},
	"filter": {
		"term": {
			"tag.key": "user:test"
		}
	}
}`
	res := nestAggregation(fewAggregationSlice)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonResult, err := json.MarshalIndent(src, "", "	")
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonResult) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonResult))
	}
}

// TODO
/* func TestAggregationNestingWithAllHandledElasticAggregationTypes(t *testing.T) {
	allTypesSlice := createAggregationPerYear([]string{""})
	allTypesSlice = append(allTypesSlice, createAggregationPerTag([]string{"", "test"})...)
	allTypesSlice = append(allTypesSlice, createAggregationPerProduct([]string{""})...)
	expectedResult := `{
	"aggregations": {
		"by-tag_key": {
			"aggregations": {
				"tag_value": {
					"aggregations": {
						"by-product": {
							"terms": {
								"field": "productCode",
								"size": 2147483647
							}
						}
					},
					"terms": {
						"field": "tag.value",
						"size": 2147483647
					}
				}
			},
			"filter": {
				"term": {
					"tag.key": "user:test"
				}
			}
		}
	},
	"date_histogram": {
		"field": "usageStartDate",
		"interval": "year",
		"min_doc_count": 0
	}
}`
	res := nestAggregation(allTypesSlice)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonResult, err := json.MarshalIndent(src, "", "	")
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonResult) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonResult))
	}
}

func TestElasticSearchParamWithNoResults(t *testing.T) {
	client, err := elastic.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	accountList := []string{"123456"}
	durationBegin, _ := time.Parse("2006-1-2 15:04", "2017-01-12 11:23")
	durationEnd, _ := time.Parse("2006-1-2 15:04", "2017-05-23 11:23")
	params := []string{"product"}
	index := "awsdetailedlineitem"
	expectedResult := `{
	"product": {
		"doc_count_error_upper_bound": 0,
		"sum_other_doc_count": 0,
		"buckets": []
	}
}`
	searchService := GetElasticSearchParams(accountList, durationBegin, durationEnd, params, client, index)
	res, err := searchService.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	aggregationMarshalled, err := json.MarshalIndent(res.Aggregations, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if string(aggregationMarshalled) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(aggregationMarshalled))
	}
}

func TestElasticSearchParamWithFewResultsAndNoAggregationNesting(t *testing.T) {
	client, err := elastic.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	accountList := []string{"394125495069"}
	durationBegin, _ := time.Parse("2006-1-2 15:04", "2017-01-12 11:23")
	durationEnd, _ := time.Parse("2006-1-2 15:04", "2017-05-23 11:23")
	params := []string{"product"}
	index := "awsdetailedlineitem"
	expectedResult := `{
	"product": {
		"doc_count_error_upper_bound": 0,
		"sum_other_doc_count": 0,
		"buckets": [
			{
				"key": "Amazon Elastic Compute Cloud",
				"doc_count": 49,
				"cost": {
					"value": 0.14448613000000002
				}
			},
			{
				"key": "Amazon DynamoDB",
				"doc_count": 18,
				"cost": {
					"value": 2.0E-8
				}
			},
			{
				"key": "Amazon Simple Storage Service",
				"doc_count": 13,
				"cost": {
					"value": 2.4008000000000002E-4
				}
			},
			{
				"key": "AWS Lambda",
				"doc_count": 7,
				"cost": {
					"value": 2.34E-6
				}
			},
			{
				"key": "Amazon Elasticsearch Service",
				"doc_count": 3,
				"cost": {
					"value": 0.1993281
				}
			},
			{
				"key": "Amazon RDS Service",
				"doc_count": 3,
				"cost": {
					"value": 0.0016027799999999998
				}
			},
			{
				"key": "Amazon Simple Notification Service",
				"doc_count": 3,
				"cost": {
					"value": 0.0
				}
			},
			{
				"key": "Amazon EC2 Container Registry (ECR)",
				"doc_count": 1,
				"cost": {
					"value": 2.765E-5
				}
			},
			{
				"key": "Amazon Route 53",
				"doc_count": 1,
				"cost": {
					"value": 3.6E-6
				}
			},
			{
				"key": "Amazon Simple Queue Service",
				"doc_count": 1,
				"cost": {
					"value": 0.0
				}
			},
			{
				"key": "AmazonCloudWatch",
				"doc_count": 1,
				"cost": {
					"value": 0.0
				}
			}
		]
	}
}`
	searchService := GetElasticSearchParams(accountList, durationBegin, durationEnd, params, client, index)
	res, err := searchService.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	aggregationMarshalled, err := json.MarshalIndent(res.Aggregations, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if string(aggregationMarshalled) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(aggregationMarshalled))
	}
}

func TestElasticSearchParamWithFewResultsAndNesting(t *testing.T) {
	client, err := elastic.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	accountList := []string{"394125495069"}
	durationBegin, _ := time.Parse("2006-1-2 15:04", "2017-01-12 11:23")
	durationEnd, _ := time.Parse("2006-1-2 15:04", "2017-05-23 11:23")
	params := []string{"product", "region", "day"}
	index := "000000-lineitems"
	expectedResult := `{
	"by-product": {
		"doc_count_error_upper_bound": 0,
		"sum_other_doc_count": 0,
		"buckets": []
	}
}`
	searchService := GetElasticSearchParams(accountList, durationBegin, durationEnd, params, client, index)
	res, err := searchService.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	aggregationMarshalled, err := json.MarshalIndent(res.Aggregations, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if string(aggregationMarshalled) != expectedResult {
		t.Fatalf("Expected %v but got: %v", expectedResult, string(aggregationMarshalled))
	}
} */