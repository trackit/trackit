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

package anomalies

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sort"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/anomaliesDetection"
)

type (
	// costWithProduct is used when a cost has to be wrapped by a product name.
	costWithProduct struct {
		product string
		cost    float64
	}

	// totalCostByDay is a named type for total cost for each day.
	totalCostByDay map[string]float64

	// highestSpendersByDay contains the more costly product podium for each day.
	highestSpendersByDay map[string][]string
)

// runAnomaliesDetectionForProducts will get data from ElasticSearch,
// compute anomalies and ingest the result in ElasticSearch.
func runAnomaliesDetectionForProducts(parsedParams AnomalyEsQueryParams, account aws.AwsAccount, ctx context.Context) (err error) {
	var res AnalyzedCosts
	if res, err = productGetAnomaliesData(ctx, parsedParams); err != nil {
	} else if err = productSaveAnomaliesData(ctx, res, account); err != nil {
	} else if err = removeRecurrence(ctx, parsedParams, account); err != nil {
	}
	return
}

// productSaveAnomaliesData will save anomalies in ElasticSearch.
// If the index doesn't exist, it will be created.
// Anomalies are unique and will replace the existing ones if
// they changed (cost or upper band).
func productSaveAnomaliesData(ctx context.Context, aCosts AnalyzedCosts, account aws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating anomalies for AWS account.", map[string]interface{}{
		"awsAccount": account,
	})
	index := es.IndexNameForUserId(account.UserId, anomaliesDetection.IndexSuffix)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, aCost := range aCosts {
		doc := anomaliesDetection.EsProductAnomaly{
			Account:   account.AwsIdentity,
			Date:      aCost.Meta.Date,
			Product:   aCost.Meta.AdditionalMeta.(AnalyzedCostProductMeta).Product,
			Abnormal:  aCost.Anomaly,
			Recurrent: false,
			Cost: anomaliesDetection.EsProductAnomalyCost{
				Value:       aCost.Cost,
				MaxExpected: aCost.UpperBand,
			},
		}
		id, err := productGenerateElasticSearchDocumentId(doc)
		if err != nil {
			logger.Error("Error when marshaling anomalies var", err.Error())
			return err
		}
		bp = addDocToBulkProcessor(bp, doc, anomaliesDetection.Type, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Failed when putting anomalies in ES", err.Error())
		return err
	}
	logger.Info("Anomalies put in ES", nil)
	return nil
}

// productGenerateElasticSearchDocumentId is used to generate the document id ingested in ElasticSearch.
// The document id is not dependent on cost or upper band: if one of them change,
// it will update the document in ElasticSearch instead of recreating one.
func productGenerateElasticSearchDocumentId(doc anomaliesDetection.EsProductAnomaly) (id string, err error) {
	var ji []byte
	ji, err = json.Marshal(struct {
		Account string `json:"account"`
		Date    string `json:"date"`
		Product string `json:"product"`
	}{
		doc.Account,
		doc.Date,
		doc.Product,
	})
	if err != nil {
		return
	}
	hash := md5.Sum(ji)
	id = base64.URLEncoding.EncodeToString(hash[:])
	return
}

// productClearDisturbances clears fake alerts with thresholds in config.
func productClearDisturbances(aCosts AnalyzedCosts, totalCostByDay totalCostByDay, highestSpendersByDay highestSpendersByDay) AnalyzedCosts {
	for index, aCost := range aCosts {
		if aCost.Anomaly {
			date := aCost.Meta.Date
			increaseAmount := aCost.Cost - aCost.UpperBand
			if increaseAmount < totalCostByDay[date]*config.AnomalyDetectionDisturbanceCleaningMinPercentOfDailyBill/100 ||
				aCost.Cost < config.AnomalyDetectionDisturbanceCleaningMinAbsoluteCost {
				aCosts[index].Anomaly = false
			} else {
				spenderInPodium := false
				for _, spender := range highestSpendersByDay[date] {
					if spender == aCost.Meta.AdditionalMeta.(AnalyzedCostProductMeta).Product {
						spenderInPodium = true
						break
					}
				}
				aCosts[index].Anomaly = spenderInPodium
			}
		}
	}
	return aCosts
}

// productAddCostToCosts is a tool used by productGetHighestSpendersByDay.
func productAddCostToCosts(product string, cost float64, costs []costWithProduct) []costWithProduct {
	for idx := range costs {
		if costs[idx].product == product {
			costs[idx].cost += cost
			return costs
		}
	}
	return append(costs, costWithProduct{product, cost})
}

// productGetHighestSpendersByDay gets a podium of the highest spenders.
func productGetHighestSpendersByDay(typedDocument anomaliesDetection.EsProductTypedResult) highestSpendersByDay {
	costByDayByProduct := map[string][]costWithProduct{}
	for _, product := range typedDocument.Products.Buckets {
		for _, date := range product.Dates.Buckets {
			costByDayByProduct[date.Key] = productAddCostToCosts(product.Key, date.Cost.Value, costByDayByProduct[date.Key])
		}
	}
	highestSpendersByDay := make(highestSpendersByDay)
	for day, products := range costByDayByProduct {
		sort.Slice(products, func(i, j int) bool {
			return products[i].cost > products[j].cost
		})
		for i := 0; i < config.AnomalyDetectionDisturbanceCleaningHighestSpendingMinRank && i < len(products); i++ {
			highestSpendersByDay[day] = append(highestSpendersByDay[day], products[i].product)
		}
	}
	return highestSpendersByDay
}

// productGetTotalCostByDay gets the total cost for each day.
func productGetTotalCostByDay(typedDocument anomaliesDetection.EsProductTypedResult) totalCostByDay {
	totalCostByDay := totalCostByDay{}
	for _, product := range typedDocument.Products.Buckets {
		for _, date := range product.Dates.Buckets {
			totalCostByDay[date.Key] += date.Cost.Value
		}
	}
	return totalCostByDay
}

// productGetAnomaliesData returns product anomalies based on query params, in JSON format.
func productGetAnomaliesData(ctx context.Context, params AnomalyEsQueryParams) (AnalyzedCosts, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sr, err := makeElasticSearchRequest(ctx, getProductElasticSearchParams, params)
	if err != nil {
		logger.Error("Failed to make elasticsearch request.", err.Error())
		return nil, err
	}
	var typedDocument anomaliesDetection.EsProductTypedResult
	if err := json.Unmarshal(*sr.Aggregations["products"], &typedDocument.Products); err != nil {
		logger.Error("Failed to parse elasticsearch document.", err.Error())
		return nil, err
	}
	totalAnalyzedCosts := make(AnalyzedCosts, 0)
	totalCostsByDay := productGetTotalCostByDay(typedDocument)
	highestSpendersByDay := productGetHighestSpendersByDay(typedDocument)
	for _, product := range typedDocument.Products.Buckets {
		aCosts := make(AnalyzedCosts, 0, len(product.Dates.Buckets))
		for _, date := range product.Dates.Buckets {
			aCosts = append(aCosts, AnalyzedCost{
				Meta: AnalyzedCostEssentialMeta{
					AdditionalMeta: AnalyzedCostProductMeta{
						Product: product.Key,
					},
					Date: date.Key,
				},
				Cost:    date.Cost.Value,
				Anomaly: false,
			})
		}
		aCosts = computeAnomalies(ctx, aCosts, params.DateBegin)
		aCosts = deleteOffset(aCosts, params.DateBegin)
		totalAnalyzedCosts = append(totalAnalyzedCosts, aCosts...)
	}
	totalAnalyzedCosts = productClearDisturbances(totalAnalyzedCosts, totalCostsByDay, highestSpendersByDay)
	if err != nil {
		return nil, err
	}
	return totalAnalyzedCosts, nil
}
