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
	"errors"
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

// UpdateTaggingComplianceForUser updates tagging compliance based on latest tagging reports and latest most used tags reports
func UpdateTaggingComplianceForUser(ctx context.Context, userId int) error {
	mostUsedTags, err := getMostUsedTagsFromDb(userId)
	if err != nil {
		return err
	}
	if mostUsedTags == nil {
		return errors.New("no most used tags data available")
	}

	count, err := getReportsCount(ctx, userId)
	if err != nil {
		return err
	}

	if len(mostUsedTags) == 0 {
		return pushComplianceToEs(ctx, userId, compliance{
			Total:           count,
			TotallyTagged:   0,
			PartiallyTagged: 0,
			NotTagged:       count,
			ReportDate:      time.Now().UTC(),
		})
	}

	totallyTagged, err := getTotallyTaggedReportsCount(ctx, userId, mostUsedTags)
	if err != nil {
		return err
	}

	untagged, err := getNotTaggedReportsCount(ctx, userId, mostUsedTags)
	if err != nil {
		return err
	}

	partiallyTagged := count - totallyTagged - untagged

	return pushComplianceToEs(ctx, userId, compliance{
		Total:           count,
		TotallyTagged:   totallyTagged,
		PartiallyTagged: partiallyTagged,
		NotTagged:       untagged,
		ReportDate:      time.Now().UTC(),
	})
}

func getMostUsedTagsFromDb(userId int) ([]string, error) {
	mostUsedTags, err := models.MostUsedTagsInUseByUser(db.Db, userId)
	if err != nil {
		return nil, err
	}
	if mostUsedTags == nil {
		return nil, nil
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

func getReportsCount(ctx context.Context, userId int) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(userId, "tagging-reports")
	index := client.Search().Index(indexName)

	res, err := index.Size(0).Query(elastic.NewMatchAllQuery()).
		Aggregation("reportDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func getTotallyTaggedReportsCount(ctx context.Context, userId int, mostUsedTags []string) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(userId, "tagging-reports")
	index := client.Search().Index(indexName)

	termQueries := mostUsedTagsToTermQueries(mostUsedTags)
	query := elastic.NewBoolQuery().Must(termQueries...)

	res, err := index.Size(0).Query(query).
		Aggregation("reportDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func getNotTaggedReportsCount(ctx context.Context, userId int, mostUsedTags []string) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(userId, "tagging-reports")
	index := client.Search().Index(indexName)

	termQueries := mostUsedTagsToTermQueries(mostUsedTags)
	query := elastic.NewBoolQuery().MustNot(termQueries...)

	res, err := index.Size(0).Query(query).
		Aggregation("reportDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func pushComplianceToEs(ctx context.Context, userId int, compliance compliance) error {
	client := es.Client
	indexName := es.IndexNameForUserId(userId, "tagging-compliance")
	_, err := client.Index().Index(indexName).Type("tagging-compliance").BodyJson(compliance).Do(ctx)
	return err
}
