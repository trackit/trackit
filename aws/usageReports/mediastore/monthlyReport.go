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

package mediastore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediastore"
	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

// getElasticSearchMediaStoreContainer prepares and run the request to retrieve the a report of an instance
// It will return the data and an error.
func getElasticSearchMediaStoreContainer(ctx context.Context, account, instance string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
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

// getContainerInfoFromEs gets information about an instance from previous report to put it in the new report
func getContainerInfoFromES(ctx context.Context, cost ContainerInformations, account string, userId int) Container {
	var docType ContainerReport
	var container = Container{
		ContainerBase: ContainerBase{
			Name:         "N/A",
			Region:     "N/A",
			Arn:      "N/A",
		},
		Costs: make(map[time.Time]float64, 0),
	}
	res, err := getElasticSearchMediaStoreContainer(ctx, account, cost.Arn,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixMediaStoreReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			container.Region = docType.Container.Region
			container.Name = docType.Container.Name
			container.Arn = docType.Container.Arn
			container.Costs = docType.Container.Costs
		}
	}
	return container
}

// fetchMonthlyContainersList sends in instanceInfoChan the instances fetched from DescribeContainers
// and filled by DescribeContainers and getContainerStats.
func fetchMonthlyContainersList(ctx context.Context, creds *credentials.Credentials,
	account, region, containerId string, cost ContainerInformations, instanceChan chan Container, startDate, endDate time.Time, userId int) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := mediastore.New(sess)
	container, err := svc.DescribeContainer(&mediastore.DescribeContainerInput{ContainerName: &containerId})
	if err != nil {
		instanceChan <- getContainerInfoFromES(ctx, cost, account, userId)
		return err
	}
	instanceChan <- Container{
		ContainerBase: ContainerBase{
			Name: aws.StringValue(container.Container.Name),
			Arn: aws.StringValue(container.Container.ARN),
			Region: cost.Region,
		},
		Costs:   cost.Cost,
	}
	return nil
}

// getMediaStoreMetrics gets credentials, accounts and region to fetch MediaStore instances stats
func fetchMonthlyContainersStats(ctx context.Context, aa taws.AwsAccount, costs []ContainerInformations, startDate, endDate time.Time) ([]ContainerReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorContainerStsSessionName)
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
	containerChans := make([]<-chan Container, 0, len(regions))
	for _, cost := range costs {
		containerRegion := getContainerRegion(cost.Arn)
		containerId := getContainerId(cost.Arn)
		for _, region := range regions {
			if strings.Contains(region, containerRegion) {
				containerChan := make(chan Container)
				go fetchMonthlyContainersList(ctx, creds, account, region, containerId, cost, containerChan, startDate, endDate, aa.UserId)
				containerChans = append(containerChans, containerChan)
			}
		}
	}
	containersList := make([]ContainerReport, 0)
	for instance := range merge(containerChans...) {
		containersList = append(containersList, ContainerReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Container: instance,
		})
	}
	return containersList, nil
}

// PutMediaStoreMonthlyReport puts a monthly report of MediaStore instance in ES
func PutMediaStoreMonthlyReport(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting MediaStore monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	costs := getMediaStoreContainerCosts(ctx, aa, startDate, endDate)
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixMediaStoreReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an MediaStore monthly report", nil)
		return false, nil
	}
	containers, err := fetchMonthlyContainersStats(ctx, aa, costs, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importContainersToEs(ctx, aa, containers)
	if err != nil {
		return false, err
	}
	return true, nil
}
