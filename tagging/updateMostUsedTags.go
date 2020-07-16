package tagging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/olivere/elastic"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
)

var ignoredTagsRegexp = []string{
	"aws:cloudformation:.*",
	"aws:autoscaling:.*",
	"lambda:.*",
	".*k8s\\.io.*",
	"KubernetesCluster",
}

// UpdateMostUsedTagsForAccount updates most used tags in MySQL for the specified account
func UpdateMostUsedTagsForAccount(ctx context.Context, account int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	mostUsedTags, err := getMostUsedTagsForAccount(ctx, account, ignoredTagsRegexp)
	if err != nil {
		return err
	}

	mostUsedTagsStr, err := json.Marshal(mostUsedTags)
	if err != nil {
		return err
	}

	reportDate := time.Now()
	model := models.MostUsedTag{
		AwsAccountID: account,
		ReportDate:   reportDate,
		Tags:         string(mostUsedTagsStr),
	}

	err = model.Insert(db.Db)
	if err == nil {
		logger.Info("Most used tags pushed to database.", map[string]interface{}{
			"reportDate": reportDate.String(),
		})
	}

	return err
}

func getMostUsedTagsForAccount(ctx context.Context, account int, ignoredTagsRegexp []string) ([]string, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(account, destIndexName)

	indexExists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !indexExists {
		return []string{}, nil
	}

	filterQueries := getFilterQueriesFromIgnoredTags(ignoredTagsRegexp)

	index := client.Search().Index(indexName)
	res, err := index.Size(0).Query(elastic.NewMatchAllQuery()).
		Aggregation("reportDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1).
			SubAggregation("nested", elastic.NewNestedAggregation().Path("tags").
				SubAggregation("filter", elastic.NewFilterAggregation().Filter(elastic.NewBoolQuery().MustNot(filterQueries...)).
					SubAggregation("terms", elastic.NewTermsAggregation().Field("tags.key").Size(5))))).Do(ctx)
	if err != nil {
		return nil, err
	}

	return processMostUsedTagsResult(res)
}

func getFilterQueriesFromIgnoredTags(ignoredTagsRegexp []string) []elastic.Query {
	queries := []elastic.Query{}

	for _, ignoredTag := range ignoredTagsRegexp {
		queries = append(queries, elastic.NewRegexpQuery("tags.key", ignoredTag))
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
