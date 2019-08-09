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

package s3

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/olivere/elastic"
	"github.com/satori/go.uuid"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/es"
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

// ReportUpdateConclusion represents the results of a bill ingestion job.
type ReportUpdateConclusion struct {
	BillRepository       BillRepository
	LastImportedManifest time.Time
	Error                error
}

// reportUpdateConclusionChanToSlice accepts a <-chan ReportUpdateConclusion
// and a count, and builds a []ReportUpdateConclusion from the values read on
// the channel, stopping at 'count' values.
func reportUpdateConclusionChanToSlice(rucc <-chan ReportUpdateConclusion, count int) (rucs []ReportUpdateConclusion) {
	rucs = make([]ReportUpdateConclusion, count)
	for i := range rucs {
		if r, ok := <-rucc; ok {
			rucs[i] = r
		} else {
			rucs = rucs[:i]
			return
		}
	}
	return
}

// UpdateDueReports finds all BillRepositories in need of an update and updates
// them.
func UpdateDueReports(ctx context.Context, tx *sql.Tx) ([]ReportUpdateConclusion, error) {
	var wg sync.WaitGroup
	aas := make(map[int]aws.AwsAccount)
	brs, err := GetAwsBillRepositoriesWithDueUpdate(tx)
	if err != nil {
		return nil, err
	}
	wg.Add(len(brs))
	conclusionChan := make(chan ReportUpdateConclusion, len(brs))
	defer close(conclusionChan)
	for _, br := range brs {
		var aa aws.AwsAccount
		var ok bool
		var err error
		if aa, ok = aas[br.AwsAccountId]; !ok {
			aa, err = aws.GetAwsAccountWithId(br.AwsAccountId, tx)
			if err != nil {
				return nil, err
			}
			aas[br.AwsAccountId] = aa
		}
		go func(ctx context.Context, aa aws.AwsAccount, br BillRepository) {
			lim, err := UpdateReport(ctx, aa, br)
			conclusionChan <- ReportUpdateConclusion{
				BillRepository:       br,
				LastImportedManifest: lim,
				Error:                err,
			}
			wg.Done()
		}(ctx, aa, br)
	}
	wg.Wait()
	return reportUpdateConclusionChanToSlice(conclusionChan, len(brs)), nil
}

// contextKey is a key in a context, to prevent collision with other modules.
type contextKey uint

// ingestionContextKey is used to store an 'ingestionId' in a context.
const ingestionContextKey = contextKey(iota)

// contextWithIngestionId returns a context configured so that its logger logs
// an 'ingestionId'.
func contextWithIngestionId(ctx context.Context) context.Context {
	ingestionId := uuid.NewV1().String()
	ctx = context.WithValue(ctx, ingestionContextKey, ingestionId)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger = logger.WithContextKey(ingestionContextKey, "ingestionId")
	logger = logger.WithContext(ctx)
	return jsonlog.ContextWithLogger(ctx, logger)
}

// UpdateReport updates the elasticsearch database with new data from usage and
// cost reports.
func UpdateReport(ctx context.Context, aa aws.AwsAccount, br BillRepository) (latestManifest time.Time, err error) {
	ctx = contextWithIngestionId(ctx)
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
			ingestLineItems(ctx, bp, index, br),
			manifestsModifiedAfter(br.LastImportedManifest),
		)
		logger.Info("Done ingesting data.", nil)
		return latestManifest, err
	}
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
func ingestLineItems(ctx context.Context, bp *elastic.BulkProcessor, index string, br BillRepository) OnLineItem {
	return func(li LineItem, ok bool) {
		if ok {
			if li.LineItemType == "Tax" {
				li.AvailabilityZone = "taxes"
				li.Region = "taxes"
			}
			li.BillRepositoryId = br.Id
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
func manifestsModifiedAfter(t time.Time) ManifestPredicate {
	return func(m manifest, oneMonthBefore bool) bool {
		if oneMonthBefore {
			if time.Time(m.LastModified).AddDate(0, 1, 0).After(t) {
				return true
			} else {
				return false
			}
		} else {
			if time.Time(m.LastModified).After(t) {
				return true
			} else {
				return false
			}
		}
	}
}

// extractTags extracts tags from a LineItem's Any field. It retrieves user
// tags only and stores them in the Tags map with a clean key.
func extractTags(li LineItem) LineItem {
	var tags []LineItemTags
	for k, v := range li.Any {
		if strings.HasPrefix(k, tagPrefix) {
			tags = append(tags, LineItemTags{strings.TrimPrefix(k, tagPrefix), v})
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
			})
		} else {
			logger.Info("Finished bulk ElasticSearch requests.", map[string]interface{}{
				"executionId": execId,
				"took":        resp.Took,
			})
		}

	}
}
