package tagging

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
)

type compliance struct {
	Total           int64 `json:"total"`
	TotallyTagged   int64 `json:"totallyTagged"`
	PartiallyTagged int64 `json:"partiallyTagged"`
	NotTagged       int64 `json:"notTagged"`
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

	totallyTagged, err := getTotallyTagged(ctx, accountID, mostUsedTags)
	if err != nil {
		return err
	}

	untagged, err := getUntagged(ctx, accountID, mostUsedTags)
	if err != nil {
		return err
	}

	partiallyTagged := count - totallyTagged - untagged

	return pushComplianceToEs(ctx, accountID, compliance{
		Total:           count,
		TotallyTagged:   totallyTagged,
		PartiallyTagged: partiallyTagged,
		NotTagged:       untagged,
	})
}

func getMostUsedTagsFromDb(accountID int) ([]string, error) {
	mostUsedTags, err := models.LatestMostUsedTagsByAwsAccountID(db.Db, accountID)
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

	reportDateAgg := elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)
	res, err := index.Size(0).Query(elastic.NewMatchAllQuery()).Aggregation("reportDate", reportDateAgg).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func getTotallyTagged(ctx context.Context, accountID int, mostUsedTags []string) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(accountID, "tagging-reports")
	index := client.Search().Index(indexName)

	termQueries := mostUsedTagsToTermQueries(mostUsedTags)
	query := elastic.NewBoolQuery().Must(termQueries...)

	reportDateAgg := elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)
	res, err := index.Size(0).Query(query).Aggregation("reportDate", reportDateAgg).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func getUntagged(ctx context.Context, accountID int, mostUsedTags []string) (int64, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(accountID, "tagging-reports")
	index := client.Search().Index(indexName)

	termQueries := mostUsedTagsToTermQueries(mostUsedTags)
	query := elastic.NewBoolQuery().MustNot(termQueries...)

	reportDateAgg := elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1)
	res, err := index.Size(0).Query(query).Aggregation("reportDate", reportDateAgg).Do(ctx)
	return handleComplianceEsReponse(res, err)
}

func pushComplianceToEs(ctx context.Context, accountID int, compliance compliance) error {
	client := es.Client
	indexName := es.IndexNameForUserId(accountID, "tagging-compliance")
	_, err := client.Index().Index(indexName).Type("tagging-compliance").BodyJson(compliance).Do(ctx)
	return err
}
