package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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
	if err != nil {
		logger.Error("Could not get tagging compliance data.", map[string]interface{}{"error": err.Error()})
		return 500, nil
	}

	return 200, res
}

func getTaggingComplianceInRange(ctx context.Context, accountID int, begin time.Time, end time.Time) (map[string]interface{}, error) {
	client := es.Client

	topHitsAgg := elastic.NewTopHitsAggregation()
	rangeAgg := elastic.NewDateRangeAggregation().Field("reportDate").AddRange(begin, end).SubAggregation("topHits", topHitsAgg)
	res, err := client.Search().Index(es.IndexNameForUserId(accountID, "tagging-compliance")).Query(elastic.NewMatchAllQuery()).Aggregation("range", rangeAgg).Do(ctx)
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
