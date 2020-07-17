//   Copyright 2020 MSolution.IO
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

package tagging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
)

type compliance struct {
	ReportDate      time.Time `json:"reportDate"`
	Total           int64     `json:"total"`
	TotallyTagged   int64     `json:"totallyTagged"`
	PartiallyTagged int64     `json:"partiallyTagged"`
	NotTagged       int64     `json:"notTagged"`
}

// UpdateTaggingComplianceForAccount updates tagging compliance based on latest tagging reports and latest most used tags reports
func UpdateTaggingComplianceForAccount(ctx context.Context, accountID int) error {
	mostUsedTags, err := getMostUsedTagsFromDb(accountID)
	if err != nil {
		return err
	}

	count, err := getReportsCount(ctx, accountID)
	if err != nil {
		return err
	}

	totallyTagged, err := getTotallyTaggedReportsCount(ctx, accountID, mostUsedTags)
	if err != nil {
		return err
	}

	untagged, err := getNotTaggedReportsCount(ctx, accountID, mostUsedTags)
	if err != nil {
		return err
	}

	partiallyTagged := count - totallyTagged - untagged

	return pushComplianceToEs(ctx, accountID, compliance{
		Total:           count,
		TotallyTagged:   totallyTagged,
		PartiallyTagged: partiallyTagged,
		NotTagged:       untagged,
		ReportDate:      time.Now().UTC(),
	})
}

func getMostUsedTagsFromDb(accountID int) ([]string, error) {
	mostUsedTags, err := models.MostUsedTagsInUseByAwsAccountID(db.Db, accountID)
	if err != nil {
		return nil, err
	}

	res := []string{}
	err = json.Unmarshal([]byte(mostUsedTags.Tags), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func mostUsedTagsToTermQueries(mostUsedTags []string) []elastic.Query {
	termQueries := []elastic.Query{}
	for _, mostUsedTag := range mostUsedTags {
		termQueries = append(termQueries, elastic.NewNestedQuery("tags", elastic.NewBoolQuery().Must(elastic.NewTermQuery("tags.key", mostUsedTag))))
	}
	return termQueries
}

func handleComplianceEsReponse(res *elastic.SearchResult, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	reportDateRes, found := res.Aggregations.Terms("reportDate")
	if !found || len(reportDateRes.Buckets) <= 0 {
		return 0, nil
	}

	return reportDateRes.Buckets[0].DocCount, nil
}

func getReportsCount(ctx context.Context, accountID int) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(accountID, "tagging-reports")
	index := client.Search().Index(indexName)

	res, err := index.Size(0).Query(elastic.NewMatchAllQuery()).
		Aggregation("reportDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func getTotallyTaggedReportsCount(ctx context.Context, accountID int, mostUsedTags []string) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(accountID, "tagging-reports")
	index := client.Search().Index(indexName)

	termQueries := mostUsedTagsToTermQueries(mostUsedTags)
	query := elastic.NewBoolQuery().Must(termQueries...)

	res, err := index.Size(0).Query(query).
		Aggregation("reportDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func getNotTaggedReportsCount(ctx context.Context, accountID int, mostUsedTags []string) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(accountID, "tagging-reports")
	index := client.Search().Index(indexName)

	termQueries := mostUsedTagsToTermQueries(mostUsedTags)
	query := elastic.NewBoolQuery().MustNot(termQueries...)

	res, err := index.Size(0).Query(query).
		Aggregation("reportDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func pushComplianceToEs(ctx context.Context, accountID int, compliance compliance) error {
	client := es.Client
	indexName := es.IndexNameForUserId(accountID, "tagging-compliance")
	_, err := client.Index().Index(indexName).Type("tagging-compliance").BodyJson(compliance).Do(ctx)
	return err
}
