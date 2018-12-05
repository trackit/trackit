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

package anomalies

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/s3"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
)

type (
	// AnalyzedCostProductMeta can be the additional metadata in AnalyzedCostEssentialMeta.
	// It's used to detect product anomalies and store them in ElasticSearch with more info.
	AnalyzedCostProductMeta struct {
		Product string
	}

	// AnalyzedCostEssentialMeta is the mandatory metadata ignored by the algorithm
	// and contains the Date and an additional metadata with specialized values.
	AnalyzedCostEssentialMeta struct {
		AdditionalMeta interface{}
		Date           string
	}

	// AnalyzedCost is returned by Bollinger Band algorithm and contains
	// every necessary data for it. It also contains metadata, ignored by
	// the algorithm.
	AnalyzedCost struct {
		Meta      AnalyzedCostEssentialMeta
		Cost      float64
		UpperBand float64
		Anomaly   bool
	}

	AnalyzedCosts []AnalyzedCost

	// AnomalyEsQueryParams will store the parsed query params
	AnomalyEsQueryParams struct {
		DateBegin time.Time
		DateEnd   time.Time
		Account   string
		Index     string
	}

	// ElasticSearchFunction is a function passed to makeElasticSearchRequest,
	// used to get results from ElasticSearch.
	ElasticSearchFunction func(
		account string,
		durationBegin time.Time,
		durationEnd time.Time,
		aggregationPeriod string,
		client *elastic.Client,
		index string,
	) *elastic.SearchService
)

// RunAnomaliesDetection run every anomaly detection algorithms and store results in ElasticSearch.
func RunAnomaliesDetection(account aws.AwsAccount, lastUpdate time.Time, ctx context.Context) (time.Time, error) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	begin := lastUpdate.Add(-1 * time.Hour * time.Duration(30*24))
	end := today.Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59))
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting anomalies detection", map[string]interface{}{
		"awsAccount": account.Id,
		"begin": begin,
		"end": end,
	})
	parsedParams := AnomalyEsQueryParams{
		DateBegin: begin,
		DateEnd:   end,
		Account:   account.AwsIdentity,
		Index:     es.IndexNameForUserId(account.UserId, s3.IndexPrefixLineItem),
	}
	return today, runAnomaliesDetectionForProducts(parsedParams, account, ctx)
}

// deleteOffset deletes the offset set in createQueryTimeRange.
func deleteOffset(aCosts AnalyzedCosts, dateBegin time.Time) AnalyzedCosts {
	var toDelete []int
	for i, aCost := range aCosts {
		if d, err := time.Parse("2006-01-02T15:04:05.000Z", aCost.Meta.Date); err == nil {
			if dateBegin.After(d) && !dateBegin.Equal(d) {
				toDelete = append(toDelete, i)
			}
		}
	}
	for n, i := range toDelete {
		aCosts = append(aCosts[:i-n], aCosts[i-n+1:]...)
	}
	return aCosts
}

// makeElasticSearchRequest prepares and run the request to retrieve the cost anomalies.
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, esFct ElasticSearchFunction, parsedParams AnomalyEsQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	searchService := esFct(
		parsedParams.Account,
		parsedParams.DateBegin,
		parsedParams.DateEnd,
		"day",
		es.Client,
		parsedParams.Index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": parsedParams.Index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, errors.GetErrorMessage(ctx, err)
		} else if err.(*elastic.Error).Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}
