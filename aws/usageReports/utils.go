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

package utils

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/es"
)

type (
	// CostPerResource associates a cost to an aws resourceId with a region
	CostPerResource struct {
		Resource string
		Cost     float64
		Region   string
	}

	// BaseReport contains basic information for any kin of usage report
	ReportBase struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		ReportType string    `json:"reportType"`
	}

	// Tag contains the key of a tag and his value
	Tag struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

const (
	mebibyte = 1 << 20

	esBulkInsertSize    = 8 * mebibyte
	esBulkInsertWorkers = 4

	opTypeCreate = "create"

	MaxAggregationSize = 0x7FFFFFFF
)

// AddDocToBulkProcessor adds a document in a bulk processor to ingest them in ES
func AddDocToBulkProcessor(bp *elastic.BulkProcessor, doc interface{}, docType, index, id string) *elastic.BulkProcessor {
	rq := elastic.NewBulkIndexRequest()
	rq = rq.Index(index)
	rq = rq.OpType(opTypeCreate)
	rq = rq.Type(docType)
	rq = rq.Id(id)
	rq = rq.Doc(doc)
	bp.Add(rq)
	return bp
}

// GetBulkProcessor builds a bulk processor for ElasticSearch.
func GetBulkProcessor(ctx context.Context) (*elastic.BulkProcessor, error) {
	bps := elastic.NewBulkProcessorService(es.Client)
	bps = bps.BulkActions(-1)
	bps = bps.BulkSize(esBulkInsertSize)
	bps = bps.Workers(esBulkInsertWorkers)
	bps = bps.Before(beforeBulk(ctx))
	bps = bps.After(afterBulk(ctx))
	return bps.Do(context.Background()) // use of background context is not an error
}

// beforeBulk returns a function that will be called before a bulk to log it
func beforeBulk(ctx context.Context) func(int64, []elastic.BulkableRequest) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	return func(execId int64, reqs []elastic.BulkableRequest) {
		logger.Info("Performing bulk ElasticSearch requests.", map[string]interface{}{
			"executionId":   execId,
			"requestsCount": len(reqs),
		})
	}
}

// beforeBulk returns a function that will be called after a bulk to log it
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

// GetAccountId gets the AWS Account ID for the given credentials
func GetAccountId(ctx context.Context, sess *session.Session) (string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := sts.New(sess)
	res, err := svc.GetCallerIdentity(nil)
	if err != nil {
		logger.Error("Error when getting caller identity", err.Error())
		return "", err
	}
	return aws.StringValue(res.Account), nil
}

// GetCurrentCheckDay returns the actual date at midnight and this date the month before
func GetCurrentCheckedDay() (start time.Time, end time.Time) {
	now := time.Now().UTC()
	end = time.Date(now.Year(), now.Month(), now.Day()-1, 24, 0, 0, 0, now.Location())
	start = time.Date(now.Year(), now.Month(), now.Day()-31, 0, 0, 0, 0, now.Location())
	return start, end
}

// FetchRegionsList fetches the regions list from AWS and returns an array of their name.
func FetchRegionsList(ctx context.Context, sess *session.Session) ([]string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := ec2.New(sess)
	regions, err := svc.DescribeRegions(nil)
	if err != nil {
		logger.Error("Error when describing regions", err.Error())
		return []string{}, err
	}
	res := make([]string, 0)
	for _, region := range regions.Regions {
		res = append(res, aws.StringValue(region.RegionName))
	}
	return res, nil
}

// CheckMonthlyReportExists checks if there is already a monthly report in ES based on the prefix.
// If there is already one it returns true, otherwise it returns false.
func CheckMonthlyReportExists(ctx context.Context, date time.Time, aa taws.AwsAccount, prefix string) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", aa.AwsIdentity))
	query = query.Filter(elastic.NewTermQuery("reportDate", date))
	index := es.IndexNameForUserId(aa.UserId, prefix)
	result, err := es.Client.Search().Index(index).Size(1).Query(query).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			logger.Warning("Query execution failed, ES index does not exists", map[string]interface{}{"index": index, "error": err.Error()})
			return false, nil
		}
		logger.Error("Query execution failed", err.Error())
		return false, err
	}
	if result.Hits.TotalHits == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
