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

// fetchDailyRoute53List sends in hostedZoneInfoChan the Hosted Zones fetched from ListHostedZones
func fetchDailyRoute53List(ctx context.Context, creds *credentials.Credentials, region string, hostedZonesChan chan HostedZone) error {
	defer close(hostedZonesChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := route53.New(sess)
	hostedZones, err := svc.ListHostedZones(nil)
	if err != nil {
		logger.Error("Error when describing Route53 HostedZones", err.Error())
		return err
	}
	for _, hostedZone := range hostedZones.HostedZones {
		ss := strings.Split(aws.StringValue(hostedZone.Id), "/")
		hostedZoneId := ss[len(ss) - 1]
		hostedZonesChan <- HostedZone{
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

	hostedZonesChan := make(chan HostedZone)
	go fetchDailyRoute53List(ctx, creds, "", hostedZonesChan)

	hostedZones := make([]HostedZoneReport, 0)
	for hostedZone := range hostedZonesChan {
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
