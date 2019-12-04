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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchDailyInstancesList sends in instanceInfoChan the instances fetched from DescribeInstances
// and filled by DescribeInstances and getInstanceStats.
func fetchDailyChannelsList(ctx context.Context, creds *credentials.Credentials, region string, channelChan chan Channel) error {
	defer close(instanceChan)
	start, end := utils.GetCurrentCheckedDay()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	var nextToken *string
	var err error
	for nextToken, err = getChannelsFromAWS(channelChan, svc, region, nextToken); nextToken != nil; {
		if err != nil {
			return err
		}
	}
	return nil
}

func getChannelsFromAWS(channelChan chan Channel, svc *mediaconvert.)

// FetchDailyInstancesStats fetches the stats of the EC2 instances of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchInstancesStats should be called every hour.
func FetchDailyChannelInputStats(ctx context.Context, awsAccount taws.AwsAccount) error {
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
	channelChan := make([]<-chan Channel, 0, len(regions))
	for _, region := range regions {
		instanceChan := make(chan Channel)
		go fetchDailyChannelsList(ctx, creds, region, instanceChan)
		channelChan = append(channelChan, instanceChan)
	}
	instances := make([]ChannelReport, 0)
	for instance := range mergeChannels(channelChan...) {
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
