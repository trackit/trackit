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
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/olivere/elastic"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type taggingComplianceSource struct {
	ReportDate      time.Time `json:"reportDate"`
	Total           int64     `json:"total"`
	TotallyTagged   int64     `json:"totallyTagged"`
	PartiallyTagged int64     `json:"partiallyTagged"`
	NotTagged       int64     `json:"notTagged"`
}

func routeGetTaggingCompliance(r *http.Request, a routes.Arguments) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(r.Context())
	u := a[users.AuthenticatedUser].(users.User)
	dateBegin := a[taggingComplianceQueryArgs[0]].(time.Time)
	dateEnd := a[taggingComplianceQueryArgs[1]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59))

	res, err := getTaggingComplianceInRange(r.Context(), u.Id, dateBegin, dateEnd)
	if err != nil && strings.Contains(err.Error(), "404") {
		return 200, map[string]interface{}{}
	}
	if err != nil {
		logger.Error("Could not get tagging compliance data.", map[string]interface{}{"error": err.Error()})
		return 500, nil
	}

	return 200, res
}

func getTaggingComplianceInRange(ctx context.Context, accountID int, begin time.Time, end time.Time) (map[string]interface{}, error) {
	client := es.Client

	res, err := client.Search().Index(es.IndexNameForUserId(accountID, "tagging-compliance")).Query(elastic.NewMatchAllQuery()).
		Aggregation("range", elastic.NewDateRangeAggregation().Field("reportDate").AddRange(begin, end).
			SubAggregation("topHits", elastic.NewTopHitsAggregation())).Do(ctx)
	if err != nil {
		return map[string]interface{}{}, err
	}

	rangeRes, found := res.Aggregations.DateRange("range")
	if !found {
		return map[string]interface{}{}, errors.New("could not query elastic search 1")
	}

	if len(rangeRes.Buckets) <= 0 {
		return map[string]interface{}{}, nil
	}

	topHitsRes, found := rangeRes.Buckets[0].Aggregations.TopHits("topHits")
	if !found {
		return map[string]interface{}{}, errors.New("could not query elastic search")
	}

	return processTaggingComplianceInRangeResults(topHitsRes)
}

func processTaggingComplianceInRangeResults(res *elastic.AggregationTopHitsMetric) (map[string]interface{}, error) {
	output := map[string]interface{}{}

	for _, hit := range res.Hits.Hits {
		source := taggingComplianceSource{}
		err := json.Unmarshal(*hit.Source, &source)
		if err != nil {
			return map[string]interface{}{}, err
		}

		output[source.ReportDate.String()] = source
	}

	return output, nil
}
