package tagging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
)

// UpdateMostUsedTagsForAccount updates most used tags in MySQL for the specified AWS account
func UpdateMostUsedTagsForAccount(ctx context.Context, account int, awsAccount string) error {
	mostUsedTags, err := getMostUsedTagsForAccount(ctx, account, []string{
		"aws:cloudformation:stack-id",
		"aws:cloudformation:logical-id",
		"aws:cloudformation:stack-name",
	})
	if err != nil {
		return err
	}

	mostUsedTagsStr, err := json.Marshal(mostUsedTags)
	if err != nil {
		return err
	}

	model := models.MostUsedTag{
		AwsAccountID: account,
		ReportDate:   time.Now(),
		Tags:         string(mostUsedTagsStr),
	}
	return model.Insert(db.Db)
}

func getMostUsedTagsForAccount(ctx context.Context, account int, ignoredTags []string) ([]string, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(account, destIndexName)

	indexExists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !indexExists {
		return []string{}, nil
	}

	filterQueries := getFilterQueriesFromIgnoredTags(ignoredTags)

	index := client.Search().Index(indexName)
	termsAgg := elastic.NewTermsAggregation().Field("tags.key").Size(5)
	filterAgg := elastic.NewFilterAggregation().Filter(elastic.NewBoolQuery().MustNot(filterQueries...)).SubAggregation("terms", termsAgg)
	nestedAgg := elastic.NewNestedAggregation().Path("tags").SubAggregation("filter", filterAgg)
	reportDateAgg := elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1).SubAggregation("nested", nestedAgg)
	res, err := index.Size(0).Query(elastic.NewMatchAllQuery()).Aggregation("reportDate", reportDateAgg).Do(ctx)
	if err != nil {
		return nil, err
	}

	return processMostUsedTagsResult(res)
}

func getFilterQueriesFromIgnoredTags(ignoredTags []string) []elastic.Query {
	queries := []elastic.Query{}

	for _, ignoredTag := range ignoredTags {
		queries = append(queries, elastic.NewTermQuery("tags.key", ignoredTag))
	}

	return queries
}

func processMostUsedTagsResult(res *elastic.SearchResult) ([]string, error) {
	reportDateRes, found := res.Aggregations.Terms("reportDate")
	if !found || len(reportDateRes.Buckets) <= 0 {
		return nil, errors.New("could not query elastic search")
	}
	nestedRes, found := reportDateRes.Buckets[0].Aggregations.Nested("nested")
	if !found {
		return nil, errors.New("could not query elastic search")
	}
	filterRes, found := nestedRes.Aggregations.Filter("filter")
	if !found {
		return nil, errors.New("could not query elastic search")
	}
	termsRes, found := filterRes.Aggregations.Terms("terms")
	if !found {
		return nil, errors.New("could not query elastic search")
	}

	mostUsedTags := []string{}

	for _, result := range termsRes.Buckets {
		mostUsedTags = append(mostUsedTags, fmt.Sprintf("%s", result.Key))
	}

	return mostUsedTags, nil
}
