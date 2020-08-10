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

package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/tagging"
	"github.com/trackit/trackit/users"
)

const maxSuggestionsCount = 10

func routeGetTaggingSuggestions(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)
	tagKey := a[suggestionsQueryArgs[0]].(string)

	suggestions, err := getSuggestions(r.Context(), u.Id, tagKey)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, map[string]interface{}{
		"tagKey":      tagKey,
		"suggestions": suggestions,
	}
}

func getSuggestions(ctx context.Context, userId int, tagKey string) ([]string, error) {
	client := es.Client

	res, err := client.Search().Index(es.IndexNameForUserId(userId, tagging.IndexPrefixTaggingReport)).Size(0).
		Aggregation("byDate", elastic.NewTermsAggregation().Field("reportDate").Order("_term", false).Size(1).
			SubAggregation("nested", elastic.NewNestedAggregation().Path("tags").
				SubAggregation("byTagKey", elastic.NewFilterAggregation().Filter(elastic.NewTermQuery("tags.key", tagKey)).
					SubAggregation("results", elastic.NewTermsAggregation().Field("tags.value").Size(maxSuggestionsCount))))).Do(ctx)
	if err != nil {
		return nil, err
	}
	return processTaggingSuggestionsResult(res)
}

func processTaggingSuggestionsResult(res *elastic.SearchResult) ([]string, error) {
	byDateRes, found := res.Aggregations.Terms("byDate")
	if !found || len(byDateRes.Buckets) <= 0 {
		return []string{}, nil
	}

	nestedRes, found := byDateRes.Buckets[0].Nested("nested")
	if !found {
		return []string{}, nil
	}

	byTagKeyRes, found := nestedRes.Aggregations.Terms("byTagKey")
	if !found || len(byDateRes.Buckets) <= 0 {
		return []string{}, nil
	}

	resultsRes, found := byTagKeyRes.Aggregations.Terms("results")
	if !found {
		return []string{}, nil
	}

	results := []string{}

	for _, buck := range resultsRes.Buckets {
		results = append(results, fmt.Sprintf("%s", buck.Key))
	}

	return results, nil
}
