package taggingrds

import (
	"context"
	"errors"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/es"
)

func fetchReports(ctx context.Context, account int) ([]*elastic.SearchHit, error) {
	client := es.Client
	indexName := es.IndexNameForUserId(account, sourceIndexName)

	indexExists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !indexExists {
		return []*elastic.SearchHit{}, nil
	}

	index := client.Search().Index(indexName)
	topHitsAggregation := elastic.NewTopHitsAggregation().Size(2147483647).FetchSourceContext(elastic.NewFetchSourceContext(true).Include("instance.id", "instance.availabilityZone", "instance.tags"))
	reportDateAggregation := elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1).SubAggregation("data", topHitsAggregation)
	res, err := index.Size(0).Query(elastic.NewTermQuery("reportType", "daily")).Aggregation("reportDate", reportDateAggregation).Do(ctx)
	if err != nil {
		return nil, err
	}

	reportDateAggregationRes, found := res.Aggregations.Terms("reportDate")
	if !found || len(reportDateAggregationRes.Buckets) <= 0 {
		return nil, errors.New("could not query elastic search")
	}

	topHitsAggregationRes, found := reportDateAggregationRes.Buckets[0].Aggregations.TopHits("data")
	if !found {
		return nil, errors.New("could not query elastic search")
	}

	return topHitsAggregationRes.Hits.Hits, nil
}
