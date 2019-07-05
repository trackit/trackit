//   Copyright 2018 MSolution.IO
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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
)

// getElasticSearchEbsSnapshot prepares and run the request to retrieve the a report of an snapshot
// It will return the data and an error.
func getElasticSearchEbsSnapshot(ctx context.Context, account, snapshot string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("snapshot.id", snapshot))
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

// getSnapshotInfoFromEs gets information about an snapshot from previous report to put it in the new report
func getSnapshotInfoFromES(ctx context.Context, snapshot utils.CostPerResource, account string, userId int) Snapshot {
	var docType SnapshotReport
	var snap = Snapshot{
		SnapshotBase: SnapshotBase{
			Id:          snapshot.Resource,
			Description: "N/A",
			State:       "N/A",
			Encrypted:   false,
			StartTime:   time.Time{},
		},
		Tags:  make([]utils.Tag, 0),
		Costs: make(map[string]float64, 0),
		Volume: Volume{
			Id:    "N/A",
			Size:  -1,
		},
	}
	res, err := getElasticSearchEbsSnapshot(ctx, account, snapshot.Resource,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixEBSReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			snap.Id = docType.Snapshot.Id
			snap.Description = docType.Snapshot.Description
			snap.State = docType.Snapshot.State
			snap.Encrypted = docType.Snapshot.Encrypted
			snap.StartTime = docType.Snapshot.StartTime
			snap.Tags = docType.Snapshot.Tags
		}
	}
	return snap
}

// fetchMonthlySnapshotsList sends in snapshotInfoChan the snapshots fetched from DescribeSnapshots
// and filled by DescribeSnapshots and getSnapshotStats.
func fetchMonthlySnapshotsList(ctx context.Context, creds *credentials.Credentials, snap utils.CostPerResource,
	account, region string, snapshotChan chan Snapshot, userId int) error {
	defer close(snapshotChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	desc := ec2.DescribeSnapshotsInput{SnapshotIds: aws.StringSlice([]string{snap.Resource})}
	snapshots, err := svc.DescribeSnapshots(&desc)
	if err != nil {
		snapshotChan <- getSnapshotInfoFromES(ctx, snap, account, userId)
		return err
	}
	for _, snapshot := range snapshots.Snapshots {
		costs := make(map[string]float64, 0)
		costs["snapshot"] = snap.Cost
		snapshotChan <- Snapshot{
			SnapshotBase: SnapshotBase{
				Id:          aws.StringValue(snapshot.SnapshotId),
				Description: aws.StringValue(snapshot.Description),
				State:       aws.StringValue(snapshot.State),
				Encrypted:   aws.BoolValue(snapshot.Encrypted),
				StartTime:   aws.TimeValue(snapshot.StartTime),
			},
			Tags:  getSnapshotTag(snapshot.Tags),
			Costs: costs,
			Volume: Volume{
				Id:   aws.StringValue(snapshot.VolumeId),
				Size: aws.Int64Value(snapshot.VolumeSize),
			},
		}
	}
	return nil
}

// getEbsMetrics gets credentials, accounts and region to fetch EBS snapshots stats
func fetchMonthlySnapshotsStats(ctx context.Context, snapshots []utils.CostPerResource, aa taws.AwsAccount, startDate time.Time) ([]SnapshotReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorSnapshotStsSessionName)
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
	snapshotChans := make([]<-chan Snapshot, 0, len(regions))
	for _, snapshot := range snapshots {
		for _, region := range regions {
			if strings.Contains(snapshot.Region, region) {
				snapshotChan := make(chan Snapshot)
				go fetchMonthlySnapshotsList(ctx, creds, snapshot, account, region, snapshotChan, aa.UserId)
				snapshotChans = append(snapshotChans, snapshotChan)
			}
		}
	}
	snapshotsList := make([]SnapshotReport, 0)
	for snapshot := range merge(snapshotChans...) {
		snapshotsList = append(snapshotsList, SnapshotReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Snapshot: snapshot,
		})
	}
	return snapshotsList, nil
}

// filterSnapshotsCosts filters snapshots, cloudwatch and volumes of EBS snapshots costs
func filterSnapshotsCosts(ebsCost, cloudwatchCost []utils.CostPerResource) ([]utils.CostPerResource, []utils.CostPerResource, []utils.CostPerResource) {
	newSnapshot := make([]utils.CostPerResource, 0)
	newVolume := make([]utils.CostPerResource, 0)
	newCloudWatch := make([]utils.CostPerResource, 0)
	for _, snapshot := range ebsCost {
		if len(snapshot.Resource) == 19 && strings.HasPrefix(snapshot.Resource, "i-") {
			newSnapshot = append(newSnapshot, snapshot)
		}
		if len(snapshot.Resource) == 21 && strings.HasPrefix(snapshot.Resource, "vol-") {
			newVolume = append(newVolume, snapshot)
		}
	}
	for _, snapshot := range cloudwatchCost {
		for _, cost := range newSnapshot {
			if strings.Contains(snapshot.Resource, cost.Resource) {
				newCloudWatch = append(newCloudWatch, snapshot)
			}
		}
	}
	return newSnapshot, newVolume, newCloudWatch
}

func addCostToSnapshots(snapshots []SnapshotReport, costVolume, costCloudWatch []utils.CostPerResource) []SnapshotReport {
	for i, snapshot := range snapshots {
		for _, volume := range snapshot.Snapshot.Stats.Volumes {
			for _, costPerVolume := range costVolume {
				if volume.Id == costPerVolume.Resource {
					snapshots[i].Snapshot.Costs[volume.Id] += costPerVolume.Cost
				}
			}
		}
		for _, cloudWatch := range costCloudWatch {
			if strings.Contains(cloudWatch.Resource, snapshot.Snapshot.Id) {
				snapshots[i].Snapshot.Costs["cloudwatch"] += cloudWatch.Cost
			}
		}
	}
	return snapshots
}

// PutEbsMonthlyReport puts a monthly report of EBS snapshot in ES
func PutEbsMonthlyReport(ctx context.Context, ebsCost, cloudWatchCost []utils.CostPerResource, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting EBS monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costSnapshot, _, _ := filterSnapshotsCosts(ebsCost, cloudWatchCost)
	//costSnapshot, costVolume, costCloudWatch := filterSnapshotsCosts(ebsCost, cloudWatchCost)
	if len(costSnapshot) == 0 {
		logger.Info("No EBS snapshots found in billing data.", nil)
		return false, nil
	}
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixEBSReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an EBS monthly report", nil)
		return false, nil
	}
	snapshots, err := fetchMonthlySnapshotsStats(ctx, costSnapshot, aa, startDate)
	if err != nil {
		return false, err
	}
	//snapshots = addCostToSnapshots(snapshots, costVolume, costCloudWatch)
	err = importSnapshotsToEs(ctx, aa, snapshots)
	if err != nil {
		return false, err
	}
	return true, nil
}
