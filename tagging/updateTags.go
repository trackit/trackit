package tagging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/es"
)

type tagObj struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type instance struct {
	ID     string   `json:"id"`
	Region string   `json:"region"`
	Tags   []tagObj `json:"tags"`
}

type source struct {
	Instance instance `json:"instance"`
}

type indexTagReport struct {
	Account      string    `json:"account"`
	ReportDate   time.Time `json:"reportDate"`
	ResourceID   string    `json:"resourceId"`
	ResourceType string    `json:"resourceType"`
	Region       string    `json:"region"`
	URL          string    `json:"url"`
	Tags         []tagObj  `json:"tags"`
}

// UpdateTagsForAccount updates tags in ES for the specified AWS account
func UpdateTagsForAccount(ctx context.Context, account int, awsAccount string) error {
	client := es.Client
	indexName := es.IndexNameForUserId(account, "ec2-reports")

	indexExists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if !indexExists {
		return nil
	}

	index := client.Search().Index(indexName)
	topHitsAggregation := elastic.NewTopHitsAggregation().Size(2147483647).FetchSourceContext(elastic.NewFetchSourceContext(true).Include("instance.id", "instance.region", "instance.tags"))
	reportDateAggregation := elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1).SubAggregation("data", topHitsAggregation)
	res, err := index.Size(0).Query(elastic.NewTermQuery("reportType", "daily")).Aggregation("reportDate", reportDateAggregation).Do(ctx)
	if err != nil {
		return err
	}

	reportDateAggregationRes, found := res.Aggregations.Terms("reportDate")
	if !found || len(reportDateAggregationRes.Buckets) <= 0 {
		return errors.New("could not query elastic search")
	}

	topHitsAggregationRes, found := reportDateAggregationRes.Buckets[0].Aggregations.TopHits("data")
	if !found {
		return errors.New("could not query elastic search")
	}

	destIndexName := es.IndexNameForUserId(account, "tagging")

	for _, hit := range topHitsAggregationRes.Hits.Hits {
		var source source

		err = json.Unmarshal(*hit.Source, &source)
		if err != nil {
			// Log error
			continue
		}

		regionForURL := source.Instance.Region
		lastChar := regionForURL[len(regionForURL)-1:]
		_, err := strconv.ParseInt(lastChar, 10, 32)

		if err != nil {
			regionForURL = regionForURL[:len(regionForURL)-1]
		}

		document := indexTagReport{
			Account:      awsAccount,
			ReportDate:   time.Now().UTC(),
			ResourceID:   source.Instance.ID,
			ResourceType: "ec2",
			Region:       source.Instance.Region,
			URL:          fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#Instances:instanceId=%s", regionForURL, regionForURL, source.Instance.ID),
			Tags:         source.Instance.Tags,
		}

		_, err = client.Index().Index(destIndexName).Type("account").BodyJson(document).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
