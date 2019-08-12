//   Copyright 2019 MSolution.IO
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

package mediapackage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediapackage"
	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

// getElasticSearchMediaPackageChannel prepares and run the request to retrieve the a report of an instance
// It will return the data and an error.
func getElasticSearchMediaPackageChannel(ctx context.Context, account, instance string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("instance.id", instance))
	search := client.Search().Index(index).Size(1).Query(query)
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, errors.GetErrorMessage(ctx, err)
	}
	return res, nil
}

// getChannelInfoFromEs gets information about an instance from previous report to put it in the new report
func getChannelInfoFromES(ctx context.Context, cost ChannelInformations, account string, userId int) Channel {
	var docType ChannelReport
	var channel = Channel{
		ChannelBase: ChannelBase{
			Id:         "N/A",
			Region:     "N/A",
			Arn:      "N/A",
		},
		Costs: make(map[time.Time]float64, 0),
	}
	res, err := getElasticSearchMediaPackageChannel(ctx, account, cost.Arn,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixMediaPackageReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			channel.Region = docType.Channel.Region
			channel.Id = docType.Channel.Id
			channel.Arn = docType.Channel.Arn
			channel.Costs = docType.Channel.Costs
		}
	}
	return channel
}

// fetchMonthlyChannelsList sends in instanceInfoChan the instances fetched from DescribeChannels
// and filled by DescribeChannels and getChannelStats.
func fetchMonthlyChannelsList(ctx context.Context, creds *credentials.Credentials,
	account, region, channelId string, cost ChannelInformations, instanceChan chan Channel, startDate, endDate time.Time, userId int) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := mediapackage.New(sess)
	channel, err := svc.DescribeChannel(&mediapackage.DescribeChannelInput{Id: &channelId})
	if err != nil {
		instanceChan <- getChannelInfoFromES(ctx, cost, account, userId)
		return err
	}
	instanceChan <- Channel{
		ChannelBase: ChannelBase{
			Id: aws.StringValue(channel.Id),
			Arn: aws.StringValue(channel.Arn),
			Region: cost.Region,
		},
		Tags: getChannelTags(channel.Tags),
		Costs:   cost.Cost,
	}
	return nil
}

// getMediaPackageMetrics gets credentials, accounts and region to fetch MediaPackage instances stats
func fetchMonthlyChannelsStats(ctx context.Context, aa taws.AwsAccount, costs []ChannelInformations, startDate, endDate time.Time) ([]ChannelReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorChannelStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return nil, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return nil, err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return nil, err
	}
	channelChans := make([]<-chan Channel, 0, len(regions))
	for _, cost := range costs {
		channelRegion := getChannelRegion(cost.Arn)
		channelId := getChannelId(cost.Arn)
		for _, region := range regions {
			if strings.Contains(region, channelRegion) {
				channelChan := make(chan Channel)
				go fetchMonthlyChannelsList(ctx, creds, account, region, channelId, cost, channelChan, startDate, endDate, aa.UserId)
				channelChans = append(channelChans, channelChan)
			}
		}
	}
	channelsList := make([]ChannelReport, 0)
	for instance := range merge(channelChans...) {
		channelsList = append(channelsList, ChannelReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Channel: instance,
		})
	}
	return channelsList, nil
}

// PutMediaPackageMonthlyReport puts a monthly report of MediaPackage instance in ES
func PutMediaPackageMonthlyReport(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting MediaPackage monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costs := getMediaPackageChannelCosts(ctx, aa, startDate, endDate)
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixMediaPackageReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an MediaPackage monthly report", nil)
		return false, nil
	}
	channels, err := fetchMonthlyChannelsStats(ctx, aa, costs, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importChannelsToEs(ctx, aa, channels)
	if err != nil {
		return false, err
	}
	return true, nil
}
