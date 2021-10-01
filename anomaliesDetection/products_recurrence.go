package anomalies

import (
	"context"
	"encoding/json"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/anomaliesDetection"
)

type (
	// anomaliesByDate is used to get an anomaly with its date more easily.
	anomaliesByDate map[time.Time]anomaliesDetection.EsProductAnomalyWithId

	// anomaliesByProduct is used to get an anomaly with its product more easily.
	anomaliesByProduct map[string]anomaliesByDate
)

// removeRecurrence gets all anomalies from ElasticSearch and removes recurrent anomalies.
func removeRecurrence(ctx context.Context, params AnomalyEsQueryParams, account aws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Removing recurrent anomalies", nil)
	params.Index = es.IndexNameForUserId(account.UserId, anomaliesDetection.Model.IndexSuffix)
	if raw, err := getAnomaliesFromEs(ctx, params); err != nil {
		return err
	} else {
		res := transformAnomaliesToMap(raw)
		var recurrentAnomalies anomaliesDetection.EsProductAnomaliesWithId
		for product := range res {
			recurrentAnomalies = append(recurrentAnomalies, detectRecurrence(res[product])...)
		}
		err := applyRecurrentAnomaliesToEs(ctx, account, recurrentAnomalies)
		return err
	}
}

// applyRecurrentAnomaliesToEs will save in ElasticSearch all recurrent anomalies
// by setting recurrent field to true.
func applyRecurrentAnomaliesToEs(ctx context.Context, account aws.AwsAccount, recurrentAnomalies anomaliesDetection.EsProductAnomaliesWithId) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating recurrent anomalies.", map[string]interface{}{
		"awsAccount": account,
	})
	index := es.IndexNameForUserId(account.UserId, anomaliesDetection.Model.IndexSuffix)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, recurrentAnomaly := range recurrentAnomalies {
		recurrentAnomaly.Source.Recurrent = true
		if err != nil {
			logger.Error("Error when marshaling recurrent anomalies var", err.Error())
			return err
		}
		bp = addDocToBulkProcessor(bp, recurrentAnomaly.Source, anomaliesDetection.Model.Type, index, recurrentAnomaly.Id)
	}
	err = bp.Flush()
	if closeErr := bp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		logger.Error("Failed when putting recurrent anomalies in ES", err.Error())
		return err
	}
	logger.Info("Recurrent anomalies put in ES", nil)
	return nil
}

// approximateCostComparison compares two float64 with
// an approximation of config.AnomalyDetectionRecurrenceCleaningThreshold.
// For example +/- 10% if it is set to 0.1.
func approximateCostComparison(a, b float64) bool {
	t := config.AnomalyDetectionRecurrenceCleaningThreshold
	return a+a*t > b && a-a*t < b
}

// detectRecurrence detects recurrent anomalies.
func detectRecurrence(an anomaliesByDate) (res anomaliesDetection.EsProductAnomaliesWithId) {
	for date := range an {
		prev := date.AddDate(0, -1, 0)
		if an[prev].Source.Abnormal && approximateCostComparison(an[date].Source.Cost.Value, an[prev].Source.Cost.Value) {
			res = append(res, an[date])
		}
	}
	return res
}

// transformAnomaliesToMap transform a raw slice of anomalies in a parsed map.
func transformAnomaliesToMap(raw anomaliesDetection.EsProductAnomaliesWithId) anomaliesByProduct {
	res := make(anomaliesByProduct)
	for _, r := range raw {
		if res[r.Source.Product] == nil {
			res[r.Source.Product] = make(anomaliesByDate)
		}
		if date, err := time.Parse("2006-01-02T15:04:05Z", r.Source.Date); err == nil {
			res[r.Source.Product][date] = r
		}
	}
	return res
}

// getAnomaliesFromEs returns product anomalies in ElasticSearch
func getAnomaliesFromEs(ctx context.Context, params AnomalyEsQueryParams) (anomaliesDetection.EsProductAnomaliesWithId, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sr, err := getAnomalyElasticSearchParams(params.Account, params.DateBegin, params.DateEnd, es.Client, params.Index, anomaliesDetection.Model.Type).Do(ctx)
	if err != nil {
		return nil, err
	}
	logger.Info("Found anomalies", map[string]interface{}{
		"begin":   params.DateBegin.String(),
		"end":     params.DateEnd.String(),
		"account": params.Account,
		"amount":  sr.Hits.TotalHits,
	})
	typedDocuments := make(anomaliesDetection.EsProductAnomaliesWithId, sr.Hits.TotalHits)
	for i, h := range sr.Hits.Hits {
		typedDocuments[i].Id = h.Id
		if b, err := h.Source.MarshalJSON(); err != nil {
			logger.Error("Failed to marshal one of the documents", map[string]interface{}{"document": h.Source})
		} else if err := json.Unmarshal(b, &typedDocuments[i].Source); err != nil {
			logger.Error("Failed to unmarshal one of the documents", map[string]interface{}{"document": string(b)})
		}
	}
	return typedDocuments, nil
}
