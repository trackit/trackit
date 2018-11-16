package anomalies

import (
	"context"
	"database/sql"
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
	AnalyzedCostProductMeta struct {
		Product string
	}

	AnalyzedCostEssentialMeta struct {
		AdditionalMeta interface{}
		Date           string
	}

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

	ElasticSearchFunction func(
		account string,
		durationBegin time.Time,
		durationEnd time.Time,
		aggregationPeriod string,
		client *elastic.Client,
		index string,
	) *elastic.SearchService
)

func RunAnomaliesDetection(account aws.AwsAccount, ctx context.Context, tx *sql.Tx) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	parsedParams := AnomalyEsQueryParams{
		DateBegin: today.Add(-1 * time.Hour * time.Duration(7*24)),
		DateEnd:   today.Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
		Account:   account.AwsIdentity,
		Index:     es.IndexNameForUserId(account.UserId, s3.IndexPrefixLineItem),
	}
	return runAnomaliesDetectionForProducts(parsedParams, account, ctx)
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
