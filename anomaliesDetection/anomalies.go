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

// Package anomalies handles detecting anomalies in trackit, providing a function to run the anomalies detection algorithms and all the types necessary to inspect detected anomalies
package anomalies

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/es"
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

	// elasticSearchDateElem is used to get usageStartDate from awsdetailedlineitems.
	elasticSearchDateElem struct {
		UsageStartDate string `json:"usageStartDate"`
	}
)

// RunAnomaliesDetection run every anomaly detection algorithms and store results in ElasticSearch.
func RunAnomaliesDetection(account aws.AwsAccount, lastUpdate time.Time, ctx context.Context) (time.Time, error) {
	esIndex := es.IndexNameForUserId(account.UserId, s3.IndexPrefixLineItem)
	begin, end, err := getDateRange(account, lastUpdate, ctx)
	if err != nil {
		return begin, err
	}
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting anomalies detection", map[string]interface{}{
		"awsAccount": account.Id,
		"begin":      begin,
		"end":        end,
	})
	parsedParams := AnomalyEsQueryParams{
		DateBegin: begin,
		DateEnd:   end,
		Account:   account.AwsIdentity,
		Index:     esIndex,
	}
	return end, runAnomaliesDetectionForProducts(parsedParams, account, ctx)
}

// makeElasticSearchDateRangeRequest makes the ElasticSearch request to get begin or end date
func makeElasticSearchDateRangeRequest(ctx context.Context, begin bool, account string, index string) (time.Time, error) {
	searchService := getDateRangeElasticSearchParams(account, begin, es.Client, index)
	res, err := searchService.Do(ctx)
	if err != nil {
		return time.Time{}, err
	}
	if len(res.Hits.Hits) == 0 {
		return time.Time{}, errors.New("empty index")
	}
	raw, err := res.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return time.Time{}, err
	}
	var elem elasticSearchDateElem
	err = json.Unmarshal(raw, &elem)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("2006-01-02T15:04:05Z", elem.UsageStartDate)
}

// getEsDateRange gets the begin and the end date from Es.
func getEsDateRange(ctx context.Context, account string, index string) (time.Time, time.Time, error) {
	begin, err := makeElasticSearchDateRangeRequest(ctx, true, account, index)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := makeElasticSearchDateRangeRequest(ctx, false, account, index)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return begin, end, err
}

// getDateRange gets the begin and the end date to launch anomaly detection.
func getDateRange(account aws.AwsAccount, lastUpdate time.Time, ctx context.Context) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	esIndex := es.IndexNameForUserId(account.UserId, s3.IndexPrefixLineItem)
	begin, end, err := getEsDateRange(ctx, account.AwsIdentity, esIndex)
	if err != nil {
		return begin, end, err
	}
	if lastUpdate.Before(begin) {
		regularEnd := begin.AddDate(0, 6, 0)
		if end.After(regularEnd) {
			end = regularEnd
		}
	} else {
		begin = time.Date(lastUpdate.Year(), lastUpdate.Month()-1, lastUpdate.Day(), 23, 59, 59, 0, lastUpdate.Location())
		end = time.Date(lastUpdate.Year(), lastUpdate.Month()+6, lastUpdate.Day(), 23, 59, 59, 0, lastUpdate.Location())
	}
	if end.After(now) {
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	}
	return begin, end, nil
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
// It will return the data and an error.
func makeElasticSearchRequest(ctx context.Context, esFct ElasticSearchFunction, parsedParams AnomalyEsQueryParams) (*elastic.SearchResult, error) {
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
		return nil, err
	}
	return res, nil
}
