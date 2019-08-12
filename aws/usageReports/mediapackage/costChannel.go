package mediapackage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	elastic2 "github.com/olivere/elastic"
	"regexp"
	"time"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/jsonlog"
	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

const aggregationMaxSize = 0x7FFFFFFF

type (
	ResponseChannelIdCostMonthly struct {
		Id struct {
			Buckets []struct {
				Key string `json:"key"`
				Date struct {
					Buckets []struct {
						Key string `json:"key_as_string"`
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
	ChannelInformations struct {
		Id     string
		Region string
		Cost   map[time.Time]float64
		Arn    string
	}
)

func getElasticSearchCost(ctx context.Context, startDate, endDate time.Time, userId int) (*elastic2.SearchResult, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewRangeQuery("usageEndDate").From(startDate).To(endDate))
	query = query.Filter(elastic.NewTermQuery("productCode", "AWSElementalMediaPackage"))
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

func getMediaPackageChannelCosts(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) []ChannelInformations {
	var response ResponseChannelIdCostMonthly
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
	channelInformations := make([]ChannelInformations, 0)
	for _, id := range response.Id.Buckets {
		channelId := getChannelId(id.Key)
		channelRegion := getChannelRegion(id.Key)
		datesCosts := make(map[time.Time]float64)
		for _, date := range id.Date.Buckets {
			totalCosts := 0.0
			for _, cost := range date.Cost.Buckets {
				totalCosts += cost.Key
			}
			dateTime, _ := time.Parse("2006-01-02T15:04:05.000Z", date.Key)
			datesCosts[dateTime] = totalCosts
		}
		channelInformations = append(channelInformations, ChannelInformations{
			Id:     channelId,
			Region: channelRegion,
			Cost:   datesCosts,
			Arn:    id.Key,
		})
	}
	return channelInformations
}

func getChannelId(resourceId string) string {
	var rgxArray []string
	reg, err := regexp.Compile("^arn:aws:mediapackage:[\\w\\d\\-]+:\\d+:channel/([\\w\\d\\-]+)")
	if err != nil {
		return ""
	} else if rgxArray = reg.FindStringSubmatch(resourceId); len(rgxArray) < 2 {
		return ""
	}
	return rgxArray[1]
}

func getChannelRegion(resourceId string) string {
	var rgxArray []string
	reg, err := regexp.Compile("^arn:aws:mediapackage:([\\w\\d\\-]+):\\d+:channel")
	if err != nil {
		return ""
	} else if rgxArray = reg.FindStringSubmatch(resourceId); len(rgxArray) < 2 {
		return ""
	}
	return rgxArray[1]
}

func getChannelTags(tags map[string]*string) map[string]string {
	formatTags := make(map[string]string)
	for key, value := range tags {
		formatTags[key] = aws.StringValue(value)
	}
	return formatTags
}