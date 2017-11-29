package s3

import (
	"context"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit2/aws"
	"github.com/trackit/trackit2/es"
)

const (
	kibibyte = 1 << 10
	mebibyte = 1 << 20
	gibibyte = 1 << 30

	esBulkInsertSize    = 8 * mebibyte
	esBulkInsertWorkers = 4

	opTypeIndex  = "index"
	opTypeCreate = "create"

	tagPrefix = `resourceTags/user:`
)

/*
func UpdateDueReports(ctx context.Context, tx *sql.Tx) error {
	aas := make(map[int]aws.AwsAccount)
	brs := aws.AwsBillRepositoriesWithDueUpdate(tx)
	for _, br := range repositoriesWithDueUpdate {
		var aa AwsAccount
		if aa, ok := aas[br.AwsAccountId]; !ok {
			aa, err := aws.GetAwsAccountWithId(br.AwsAccountId, tx)
			if err != nil {
				return err
			}
			aas[br.AwsAccountId] = aa
		}
		go UpdateReport(ctx, aa, br)
	}
}
*/

// UpdateReport updates the elasticsearch database with new data from usage and
// cost reports.
func UpdateReport(ctx context.Context, aa aws.AwsAccount, br BillRepository) (latestManifest time.Time, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating reports for AWS account.", map[string]interface{}{
		"awsAccount":     aa,
		"billRepository": br,
	})
	if bp, err := getBulkProcessor(ctx); err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return latestManifest, err
	} else {
		index := es.IndexNameForUserId(aa.UserId, IndexPrefixLineItem)
		latestManifest, err = ReadBills(
			ctx,
			aa,
			br,
			ingestLineItems(ctx, bp, index),
			manifestsStartingAfter(br.LastImportedPeriod),
		)
	}
	return
}

// getBulkProcessor builds a bulk processor for ElasticSearch.
func getBulkProcessor(ctx context.Context) (*elastic.BulkProcessor, error) {
	bps := elastic.NewBulkProcessorService(es.Client)
	bps = bps.BulkActions(-1)
	bps = bps.BulkSize(esBulkInsertSize)
	bps = bps.Workers(esBulkInsertWorkers)
	bps = bps.Before(beforeBulk(ctx))
	bps = bps.After(afterBulk(ctx))
	return bps.Do(context.Background()) // use of background context is not an error
}

// ingestLineItems returns an OnLineItem handler which ingests LineItems in an
// ElasticSearch index.
func ingestLineItems(ctx context.Context, bp *elastic.BulkProcessor, index string) OnLineItem {
	return func(li LineItem, ok bool) {
		if ok {
			li = extractTags(li)
			rq := elastic.NewBulkIndexRequest()
			rq = rq.Index(index)
			rq = rq.OpType(opTypeCreate)
			rq = rq.Type(TypeLineItem)
			rq = rq.Id(li.EsId())
			rq = rq.Doc(li)
			bp.Add(rq)
		} else {
			bp.Flush()
			bp.Close()
		}
	}
}

// manifestsStartingAfter returns a manifest predicate which is true for all
// manifests starting after a given date.
func manifestsStartingAfter(t time.Time) ManifestPredicate {
	return func(m manifest) bool {
		if time.Time(m.BillingPeriod.Start).After(t) {
			return true
		} else {
			return false
		}
	}
}

// extractTags extracts tags from a LineItem's Any field. It retrieves user
// tags only and stores them in the Tags map with a clean key.
func extractTags(li LineItem) LineItem {
	tags := make(map[string]string)
	for k, v := range li.Any {
		if strings.HasPrefix(k, tagPrefix) {
			tags[strings.TrimPrefix(k, tagPrefix)] = v
		}
	}
	li.Tags = tags
	li.Any = nil
	return li
}

func beforeBulk(ctx context.Context) func(int64, []elastic.BulkableRequest) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	return func(execId int64, reqs []elastic.BulkableRequest) {
		logger.Info("Performing bulk ElasticSearch requests.", map[string]interface{}{
			"executionId":   execId,
			"requestsCount": len(reqs),
		})
	}
}

func afterBulk(ctx context.Context) func(int64, []elastic.BulkableRequest, *elastic.BulkResponse, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	return func(execId int64, reqs []elastic.BulkableRequest, resp *elastic.BulkResponse, err error) {
		if err != nil {
			logger.Error("Failed bulk ElasticSearch requests.", map[string]interface{}{
				"executionId": execId,
				"error":       err.Error(),
				"took":        resp.Took,
			})
		} else {
			logger.Info("Finished bulk ElasticSearch requests.", map[string]interface{}{
				"executionId": execId,
				"took":        resp.Took,
			})
		}

	}
}
