package medialive

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
)

func fetchChannels(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time, regions []string, account string,
	creds *credentials.Credentials) []ChannelReport {
	channelsCosts := getMediaLiveChannelCosts(ctx, aa, startDate, endDate)
	channelChans := make([]<-chan Channel, 0, len(regions))
	for _, cost := range channelsCosts {
		for _, region := range regions {
			if strings.Contains(cost.Region, region) && cost.Id != "" {
				channelChan := make(chan Channel)
				go fetchMonthlyChannelsList(ctx, creds, cost, account, region, channelChan, aa.UserId)
				channelChans = append(channelChans, channelChan)
			}
		}
	}
	channelsList := make([]ChannelReport, 0)
	for channel := range mergeChannels(channelChans...) {
		channelsList = append(channelsList, ChannelReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Channel: channel,
		})
	}
	return channelsList
}

func fetchInputs(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time, regions []string, account string,
	creds *credentials.Credentials) []InputReport {
	inputsCosts := getMediaLiveInputCosts(ctx, aa, startDate, endDate)
	inputChans := make([]<-chan Input, 0, len(regions))
	for _, cost := range inputsCosts {
		for _, region := range regions {
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
	}
	return inputsList
}
