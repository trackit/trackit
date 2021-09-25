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

package onDemandToRiEc2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	terrors "github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

type (
	// RiEc2QueryParams will store the parsed query params
	RiEc2QueryParams struct {
		AccountList []string
		IndexList   []string
		DateBegin   time.Time
		DateEnd     time.Time
	}

	// ResponseRiEc2Reports allows us to parse ES response for RI EC2 reports
	ResponseRiEc2Reports struct {
		Accounts struct {
			Buckets []struct {
				Key string `json:"key"`
				Reports struct {
					Hits struct {
						Hits []struct {
							Report OdToRiEc2Report `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"reports"`
			} `json:"buckets"`
		} `json:"accounts"`
	}
)

// makeElasticSearchRequest prepares and run an ES request
// based on the RiEc2QueryParams and search params
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 Internal Server Error status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams RiEc2QueryParams,
	esSearchParams func(RiEc2QueryParams, *elastic.Client, string) *elastic.SearchService) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.IndexList, ",")
	searchService := esSearchParams(
		parsedParams,
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, terrors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, terrors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}

// createQueryAccountFilterRiEc2 creates and returns a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterRiEc2(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// createQueryTimeRange creates and returns a new *elastic.RangeQuery based on the duration
// defined by durationBegin and durationEnd
func createQueryTimeRange(durationBegin time.Time, durationEnd time.Time) *elastic.RangeQuery {
	return elastic.NewRangeQuery("reportDate").
		From(durationBegin).To(durationEnd)
}

// getElasticSearchRiEc2Params is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as parameters :
// 	- params RiEc2QueryParams : contains the list of accounts and the date
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on which to execute the query. In this context the default value
//	should be "od-to-ri-ec2-report"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func getElasticSearchRiEc2Params(params RiEc2QueryParams, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilterRiEc2(params.AccountList))
	}
	query = query.Filter(createQueryTimeRange(params.DateBegin, params.DateEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts", elastic.NewTermsAggregation().Field("account").
		SubAggregation("reports", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(1)))
	return search
}

// prepareResponseRiEc2 parses the results from elasticsearch and returns an array of RI EC2 report
func prepareResponseRiEc2(ctx context.Context, resEc2 *elastic.SearchResult) ([]OdToRiEc2Report, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	var response ResponseRiEc2Reports
	reports := make([]OdToRiEc2Report, 0)
	err := json.Unmarshal(*resEc2.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES EC2 response", err)
		return nil, terrors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, report := range account.Reports.Hits.Hits {
			reports = append(reports, report.Report)
		}
	}
	return reports, nil
}
