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

package medialive

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

// getElasticSearchMedialiveChannel prepares and run the request to retrieve the report of a channel
// It will return the data and an error.
// Index is the index of the name for the user ID
func getElasticSearchMedialiveChannel(ctx context.Context, account, channel string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("channel.id", channel))
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

// getChannelInfoFromEs gets information about an channel from previous report to put it in the new report
func getChannelInfoFromES(ctx context.Context, channelCosts ChannelInformations, account string, userId int) Channel {
	var docType ChannelReport
	var chann = Channel{
		ChannelBase: ChannelBase{
			Arn:       channelCosts.Arn,
			Id:        channelCosts.Id,
			Name:      "N/A",
			Region:    channelCosts.Region,
		},
		Tags: make(map[string]string, 0),
		Cost: channelCosts.Cost,
	}
	res, err := getElasticSearchMedialiveChannel(ctx, account, channelCosts.Arn,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixMediaLiveReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			chann.Id = docType.Channel.Id
			chann.Arn = docType.Channel.Arn
			chann.Name = docType.Channel.Name
			chann.Region = docType.Channel.Region
			chann.Cost = docType.Channel.Cost
			chann.Tags = docType.Channel.Tags
		}
	}
	return chann
}

// getChannelInfoFromEs gets information about an channel from previous report to put it in the new report
func getInputInfoFromES(ctx context.Context, inputCosts InputInformations, account string, userId int) Input {
	var docType InputReport
	var input = Input{
		InputBase: InputBase{
			Arn:       inputCosts.Arn,
			Id:        inputCosts.Id,
			Name:      "N/A",
			Region:    inputCosts.Region,
		},
		Tags: make(map[string]string, 0),
		Cost: inputCosts.Cost,
	}
	res, err := getElasticSearchMedialiveChannel(ctx, account, inputCosts.Arn,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixMediaLiveReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			input.Id = docType.Input.Id
			input.Arn = docType.Input.Arn
			input.Name = docType.Input.Name
			input.Region = docType.Input.Region
			input.Cost = docType.Input.Cost
			input.Tags = docType.Input.Tags
		}
	}
	return input
}

// fetchMonthlyChannelsList sends in channelInfoChan the channels fetched from DescribeChannels
// and filled by DescribeChannels and getChannelStats.
func fetchMonthlyChannelsList(ctx context.Context, creds *credentials.Credentials, cost ChannelInformations,
	account, region string, channelChan chan Channel, userId int) error {
	defer close(channelChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := medialive.New(sess)
	describeChannel, err := svc.DescribeChannel(&medialive.DescribeChannelInput{ChannelId: &cost.Id})
	if err != nil {
		channelChan <- getChannelInfoFromES(ctx, cost, account, userId)
		return err
	}
	channelChan <- Channel{
		ChannelBase: ChannelBase{
			Arn:    aws.StringValue(describeChannel.Arn),
			Id:     aws.StringValue(describeChannel.Id),
			Name:   aws.StringValue(describeChannel.Name),
			Region: region,
		},
		Tags: getChannelTags(describeChannel.Tags),
		Cost: cost.Cost,
	}
	return nil
}

// fetchMonthlyInput sends in channelInfoChan the channels fetched from DescribeChannels
// and filled by DescribeChannels and getChannelStats.
func fetchMonthlyInput(ctx context.Context, creds *credentials.Credentials, cost InputInformations,
	account, region string, inputChan chan Input, userId int) error {
	defer close(inputChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := medialive.New(sess)
	describeInput, err := svc.DescribeInput(&medialive.DescribeInputInput{InputId: &cost.Id})
	if err != nil {
		inputChan <- getInputInfoFromES(ctx, cost, account, userId)
		return err
	}
	inputChan <- Input{
		InputBase: InputBase{
			Arn:    aws.StringValue(describeInput.Arn),
			Id:     aws.StringValue(describeInput.Id),
			Name:   aws.StringValue(describeInput.Name),
			Region: region,
		},
		Tags: getInputTags(describeInput.Tags),
		Cost: cost.Cost,
	}
	return nil
}

// fetchMonthlyChannelsInputsStats fetch MediaLive channels stats
func fetchMonthlyChannelsInputsStats(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) ([]ChannelReport, []InputReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorChannelStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return nil, nil, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return nil, nil, err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return nil, nil, err
	}
	/*channelsCosts := getMediaLiveChannelCosts(ctx, aa, startDate, endDate)
	channelChans := make([]<-chan Channel, 0, len(regions))
	for _, cost := range channelsCosts {
		for _, region := range regions {
			log.Printf("REGION = %v\n", region)
			if strings.Contains(cost.Region, region) && cost.Id != "" {
				channelChan := make(chan Channel)
				go fetchMonthlyChannelsList(ctx, creds, cost, account, region, channelChan, aa.UserId)
				channelChans = append(channelChans, channelChan)
			}
		}
	}
	channelsList := make([]ChannelReport, 0)
	for channel := range merge(channelChans...) {
		channelsList = append(channelsList, ChannelReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Channel: channel,
		})
	}
	inputsCosts := getMediaLiveInputCosts(ctx, aa, startDate, endDate)
	inputChans := make([]<-chan Input, 0, len(regions))
	for _, cost := range inputsCosts {
		for _, region := range regions {
			log.Printf("REGION = %v\n", region)
			if strings.Contains(cost.Region, region) && cost.Id != "" {
				inputChan := make(chan Input)
				go fetchMonthlyInput(ctx, creds, cost, account, region, inputChan, aa.UserId)
				inputChans = append(inputChans, inputChan)
			}
		}
	}
	inputsList := make([]InputReport, 0)
	for input := range mergeInput(inputChans...) {
		inputsList = append(inputsList, InputReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Input: input,
		})
	}*/
	channelsList := fetchChannels(ctx, aa, startDate, endDate, regions, account, creds)
	inputsList := fetchInputs(ctx, aa, startDate, endDate, regions, account, creds)
	return channelsList, inputsList, nil
}

// PutMedialiveMonthlyReport puts a monthly report of MediaLive channel in ES
func PutMedialiveMonthlyReport(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting MediaLive monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})/*
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixMediaLiveReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an MediaLive monthly report", nil)
		return false, nil
	}*/
	channels, inputs, err := fetchMonthlyChannelsInputsStats(ctx, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importChannelsToEs(ctx, aa, channels)
	if err != nil {
		return false, err
	}
	err = importInputsToEs(ctx, aa, inputs)
	if err != nil {
		return false, err
	}
	return true, nil
}