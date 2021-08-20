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

package route53

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// getCompleteHostedZoneList calls ListHostedZones repeatedly until all the zones have been fetched (this is required if the account contains more than 100 elements, as a single ListHostedZones can only return at most 100 zones)
func getCompleteHostedZoneList(svc *route53.Route53) ([]*route53.HostedZone, error) {
	listHostedZonesOutput, err := svc.ListHostedZones(nil)
	if err != nil {
		return nil, err
	}
	result := listHostedZonesOutput.HostedZones
	for *listHostedZonesOutput.IsTruncated {
		listHostedZonesOutput, err = svc.ListHostedZones(&route53.ListHostedZonesInput{
			Marker: listHostedZonesOutput.NextMarker,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, listHostedZonesOutput.HostedZones...)
	}
	return result, nil
}

// fetchDailyRoute53List sends in hostedZoneInfoChan the Hosted Zones fetched from ListHostedZones
func fetchDailyRoute53List(ctx context.Context, creds *credentials.Credentials, region string, hostedZoneChan chan HostedZone) error {
	defer close(hostedZoneChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := route53.New(sess)
	hostedZones, err := getCompleteHostedZoneList(svc)
	if err != nil {
		logger.Error("Error when getting Route53 Hosted Zones list", err.Error())
		return err
	}
	for _, hostedZone := range hostedZones {
		ss := strings.Split(aws.StringValue(hostedZone.Id), "/")
		hostedZoneId := ss[len(ss) - 1]
		hostedZoneChan <- HostedZone{
			HostedZoneBase: HostedZoneBase{
				Name:   aws.StringValue(hostedZone.Name),
				Id:     hostedZoneId,
				Region: region,
			},
			Tags: getRoute53Tags(ctx, hostedZone, svc),
		}
	}
	return nil
}

// FetchDailyRoute53Stats fetches the stats of the Route53 Hosted Zones of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchRoute53Stats should be called every hour.
func FetchDailyRoute53Stats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching Route53 stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorRoute53StsSessionName)
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
	hostedZonesChans := make([]<-chan HostedZone, 0, len(regions))
	for _, region := range regions {
		hostedZoneChan := make(chan HostedZone)
		go fetchDailyRoute53List(ctx, creds, region, hostedZoneChan)
		hostedZonesChans = append(hostedZonesChans, hostedZoneChan)
	}
	hostedZones := make([]HostedZoneReport, 0)
	for hostedZone := range merge(hostedZonesChans...) {
		hostedZones = append(hostedZones, HostedZoneReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			HostedZone: hostedZone,
		})
	}
	return importRoute53ToEs(ctx, awsAccount, hostedZones)
}
