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

package ebs

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/ebsReports"
)

// fetchDailySnapshotsList sends in snapshotInfoChan the snapshots fetched from DescribeSnapshots
// and filled by DescribeSnapshots and getSnapshotStats.
func fetchDailySnapshotsList(ctx context.Context, creds *credentials.Credentials, awsAccount taws.AwsAccount, region string, snapshotChan chan ebsReports.Snapshot) error {
	defer close(snapshotChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	snapshots, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{
			aws.String(awsAccount.AwsIdentity),
		},
	})
	if err != nil {
		logger.Error("Error when describing snapshots", err.Error())
		return err
	}
	for _, snapshot := range snapshots.Snapshots {
		snapshotChan <- ebsReports.Snapshot{
			SnapshotBase: ebsReports.SnapshotBase{
				Id:          aws.StringValue(snapshot.SnapshotId),
				Description: aws.StringValue(snapshot.Description),
				State:       aws.StringValue(snapshot.State),
				Encrypted:   aws.BoolValue(snapshot.Encrypted),
				StartTime:   aws.TimeValue(snapshot.StartTime),
				Region:      region,
			},
			Tags: getSnapshotTag(snapshot.Tags),
			Volume: ebsReports.Volume{
				Id:   aws.StringValue(snapshot.VolumeId),
				Size: aws.Int64Value(snapshot.VolumeSize),
			},
			Cost: 0,
		}
	}
	return nil
}

// FetchDailySnapshotsStats fetches the stats of the EBS snapshots of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchSnapshotsStats should be called every hour.
func FetchDailySnapshotsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching EBS snapshot stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorSnapshotStsSessionName)
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
	snapshotChans := make([]<-chan ebsReports.Snapshot, 0, len(regions))
	for _, region := range regions {
		snapshotChan := make(chan ebsReports.Snapshot)
		go fetchDailySnapshotsList(ctx, creds, awsAccount, region, snapshotChan)
		snapshotChans = append(snapshotChans, snapshotChan)
	}
	snapshots := make([]ebsReports.SnapshotReport, 0)
	for snapshot := range merge(snapshotChans...) {
		snapshots = append(snapshots, ebsReports.SnapshotReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Snapshot: snapshot,
		})
	}
	return importSnapshotsToEs(ctx, awsAccount, snapshots)
}
