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
