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
	mostUsedTags, err := getMostUsedTagsForAccount(ctx, account)
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

func getMostUsedTagsForAccount(ctx context.Context, account int) ([]string, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(account, destIndexName)

	indexExists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !indexExists {
		return []string{}, nil
	}

	index := client.Search().Index(indexName)
	topTags := elastic.NewTermsAggregation().Field("tags.key").Size(5)
	tags := elastic.NewNestedAggregation().Path("tags").SubAggregation("topTags", topTags)
	reportDate := elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1).SubAggregation("tags", tags)
	res, err := index.Size(0).Query(elastic.NewMatchAllQuery()).Aggregation("reportDate", reportDate).Do(ctx)
	if err != nil {
		return nil, err
	}

	return processMostUsedTagsResult(res)
}

func processMostUsedTagsResult(res *elastic.SearchResult) ([]string, error) {
	reportDateRes, found := res.Aggregations.Terms("reportDate")
	if !found || len(reportDateRes.Buckets) <= 0 {
		return nil, errors.New("could not query elastic search")
	}
	tagsRes, found := reportDateRes.Buckets[0].Aggregations.Nested("tags")
	if !found {
		return nil, errors.New("could not query elastic search")
	}
	topTagsRes, found := tagsRes.Aggregations.Terms("topTags")
	if !found {
		return nil, errors.New("could not query elastic search")
	}

	mostUsedTags := []string{}

	for _, result := range topTagsRes.Buckets {
		mostUsedTags = append(mostUsedTags, fmt.Sprintf("%s", result.Key))
	}

	return mostUsedTags, nil
}
