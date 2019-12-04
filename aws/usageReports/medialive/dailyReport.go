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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchDailyChannelList sends in channelInfoChan the instances fetched from ListChannels
func fetchDailyChannelsList(_ context.Context, creds *credentials.Credentials, region string, channelChan chan Channel) error {
	defer close(channelChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := medialive.New(sess)
	var nextToken *string
	var err error
	for nextToken, err = getChannelsFromAWS(channelChan, svc, region, nextToken); nextToken != nil; {
		if err != nil {
			return err
		}
	}
	return nil
}

func getChannelsFromAWS(channelChan chan Channel, svc *medialive.MediaLive, region string, token *string) (*string, error) {
	listChannel, err := svc.ListChannels(&medialive.ListChannelsInput{NextToken: token})
	if err != nil {
		return nil, err
	}
	for _, channel := range listChannel.Channels {
		channelChan <- Channel{
			ChannelBase:           ChannelBase{
				Arn: aws.StringValue(channel.Arn),
				Id: aws.StringValue(channel.Id),
				Name: aws.StringValue(channel.Name),
				Region: region,
			},
			ChannelClass:          aws.StringValue(channel.ChannelClass),
			LogLevel:              aws.StringValue(channel.LogLevel),
			PipelinesRunningCount: aws.Int64Value(channel.PipelinesRunningCount),
			State:                 aws.StringValue(channel.State),
			Tags:                  getChannelTags(channel.Tags),
			Cost:                  nil,
		}
	}
	return listChannel.NextToken, nil
}

// fetchDailyInputList sends in inputInfoChan the inputs fetched from ListInputs
func fetchDailyInputsList(_ context.Context, creds *credentials.Credentials, region string, channelChan chan Input) error {
	defer close(channelChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := medialive.New(sess)
	var nextToken *string
	var err error
	for nextToken, err = getInputsFromAWS(channelChan, svc, region, nextToken); nextToken != nil; {
		if err != nil {
			return err
		}
	}
	return nil
}

func getInputsFromAWS(inputChan chan Input, svc *medialive.MediaLive, region string, token *string) (*string, error) {
	listInput, err := svc.ListInputs(&medialive.ListInputsInput{NextToken: token})
	if err != nil {
		return nil, err
	}
	for _, input := range listInput.Inputs {
		inputChan <- Input{
			InputBase:           InputBase{
				Arn: aws.StringValue(input.Arn),
				Id: aws.StringValue(input.Id),
				Name: aws.StringValue(input.Name),
				Region: region,
			},
			AttachedChannels:          aws.StringValueSlice(input.AttachedChannels),
			InputClass:              aws.StringValue(input.InputClass),
			RoleArn: aws.StringValue(input.RoleArn),
			SecurityGroups: aws.StringValueSlice(input.SecurityGroups),
			State:                 aws.StringValue(input.State),
			Type: aws.StringValue(input.Type),
			Tags:                  getChannelTags(input.Tags),
			Cost:                  nil,
		}
	}
	return listInput.NextToken, nil
}

// FetchDailyChannelStats fetches the stats of the Medialive Channels of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchChannelStats should be called every hour.
func fetchDailyChannelStats(ctx context.Context, awsAccount taws.AwsAccount, now time.Time, account string, regions []string, creds *credentials.Credentials) error {
	channelsChan := make([]<-chan Channel, 0, len(regions))
	for _, region := range regions {
		chanChan := make(chan Channel)
		go fetchDailyChannelsList(ctx, creds, region, chanChan)
		channelsChan = append(channelsChan, chanChan)
	}
	instances := make([]ChannelReport, 0)
	for instance := range mergeChannels(channelsChan...) {
		instances = append(instances, ChannelReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Channel: instance,
		})
	}
	return importChannelsToEs(ctx, awsAccount, instances)
}

// FetchDailyInputStats fetches the stats of the Medialive Inputs of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchInputStats should be called every hour.
func fetchDailyInputStats(ctx context.Context, awsAccount taws.AwsAccount, now time.Time, account string, regions []string, creds *credentials.Credentials) error {
	inputsChan := make([]<-chan Input, 0, len(regions))
	for _, region := range regions {
		inputChan := make(chan Input)
		go fetchDailyInputsList(ctx, creds, region, inputChan)
		inputsChan = append(inputsChan, inputChan)
	}
	inputs := make([]InputReport, 0)
	for input := range mergeInput(inputsChan...) {
		inputs = append(inputs, InputReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Input: input,
		})
	}
	return importInputsToEs(ctx, awsAccount, inputs)
}

func FetchDailyChannelsInputsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching EC2 instance stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorChannelStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	now := time.Now().UTC()
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	if err = fetchDailyChannelStats(ctx, awsAccount, now, account, regions, creds); err != nil {
		return err
	} else if err = fetchDailyInputStats(ctx, awsAccount, now, account, regions, creds); err != nil {
		return err
	}
	return nil
}

