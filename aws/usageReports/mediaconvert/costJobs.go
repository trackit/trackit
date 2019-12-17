package mediaconvert

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

const aggregationMaxSize = 0x7FFFFFFF

type (
	ResponseJobIdCostMonthly struct {
		Id struct {
			Buckets []struct {
				Key  string `json:"key"`
				Date struct {
					Buckets []struct {
						Key  string `json:"key_as_string"`
						Cost struct {
							Buckets []struct {
								Key float64 `json:"key"`
							} `json:"buckets"`
						} `json:"unblendedCost"`
					} `json:"buckets"`
				} `json:"usageStartDate"`
			} `json:"buckets"`
		} `json:"resourceId"`
	}
	JobInformations struct {
		Id     string
		Region string
		Cost   float64
		Arn    string
	}
)

//getElasticSearchCost return a result from a request to the ES about Elemental MediaConvert billing in lineitems
func getElasticSearchCost(ctx context.Context, startDate, endDate time.Time, userId int) (*elastic.SearchResult, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewRangeQuery("usageEndDate").From(startDate).To(endDate))
	query = query.Filter(elastic.NewTermQuery("productCode", "AWSElementalMediaConvert"))
	search := es.Client.Search().Index(es.IndexNameForUserId(userId, es.IndexPrefixLineItems)).Size(0).Query(query)
	search.Aggregation("resourceId", elastic.NewTermsAggregation().Field("resourceId").Size(aggregationMaxSize).
		SubAggregation("usageStartDate", elastic.NewDateHistogramAggregation().Field("usageStartDate").MinDocCount(0).Interval("hour").
			SubAggregation("unblendedCost", elastic.NewTermsAggregation().Field("unblendedCost").Size(aggregationMaxSize))))
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			logger.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": es.IndexNameForUserId(userId, es.IndexPrefixLineItems),
				"error": err.Error(),
			})
			return nil, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			logger.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			logger.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, errors.GetErrorMessage(ctx, err)
	}
	return res, nil
}

func getMediaConvertJobCosts(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) []JobInformations {
	var response ResponseJobIdCostMonthly
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res, err := getElasticSearchCost(ctx, startDate, endDate, aa.UserId)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(*res.Aggregations["resourceId"], &response.Id)
	if err != nil {
		logger.Error("Unmarshal execution failedd", err)
		return nil
	}
	jobInformations := make([]JobInformations, 0)
	for _, id := range response.Id.Buckets {
		jobId := getJobId(id.Key)
		jobRegion := getJobRegion(id.Key)
		datesCosts := make(map[time.Time]float64)
		for _, date := range id.Date.Buckets {
			totalCosts := 0.0
			for _, cost := range date.Cost.Buckets {
				totalCosts += cost.Key
			}
			dateTime, _ := time.Parse("2006-01-02T15:04:05.000Z", date.Key)
			datesCosts[dateTime] = totalCosts
		}
		jobInformations = append(jobInformations, JobInformations{
			Id:     jobId,
			Region: jobRegion,
			Cost:   getTotalCost(datesCosts),
			Arn:    id.Key,
		})
	}
	return jobInformations
}

func getTotalCost(costs map[time.Time]float64) float64 {
	var totalCost float64

	totalCost = 0
	for _, cost := range costs {
		totalCost += cost
	}
	return totalCost
}

func getJobId(resourceId string) string {
	var rgxArray []string
	reg, err := regexp.Compile("^arn:aws:mediaconvert:[\\w\\d\\-]+:\\d+:job/([\\w\\d\\-]+)")
	if err != nil {
		return ""
	} else if rgxArray = reg.FindStringSubmatch(resourceId); len(rgxArray) < 2 {
		return ""
	}
	return rgxArray[1]
}

func getJobRegion(resourceId string) string {
	var rgxArray []string
	reg, err := regexp.Compile("^arn:aws:mediaconvert:([\\w\\d\\-]+):\\d+:job")
	if err != nil {
		return ""
	} else if rgxArray = reg.FindStringSubmatch(resourceId); len(rgxArray) < 2 {
		return ""
	}
	return rgxArray[1]
}
