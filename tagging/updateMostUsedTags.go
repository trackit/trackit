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
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
)

var ignoredTagsRegexp = []string{
	"aws:*",
	"lambda:.*",
	".*k8s\\.io.*",
	"KubernetesCluster",
}

// UpdateMostUsedTagsForUser updates most used tags in MySQL for the specified user
func UpdateMostUsedTagsForUser(ctx context.Context, userId int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Updating most used tags.", map[string]interface{}{
		"userId": userId,
	})

	// This sleep is required because ES only updates his indexes every second
	// https://stackoverflow.com/questions/18078561/elasticsearch-get-just-after-post
	time.Sleep(time.Second * 2)

	mostUsedTags, err := getMostUsedTagsForUser(ctx, userId, ignoredTagsRegexp)
	if err != nil {
		return err
	}

	mostUsedTagsStr, err := json.Marshal(mostUsedTags)
	if err != nil {
		return err
	}

	reportDate := time.Now()
	model := models.MostUsedTag{
		UserID:     userId,
		ReportDate: reportDate,
		Tags:       string(mostUsedTagsStr),
	}

	err = model.Insert(db.Db)
	if err == nil {
		logger.Info("Most used tags pushed to database.", map[string]interface{}{
			"reportDate": reportDate.String(),
		})
	}

	return err
}

func getMostUsedTagsForUser(ctx context.Context, userId int, ignoredTagsRegexp []string) ([]string, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(userId, destIndexName)

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
