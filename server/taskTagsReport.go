//   Copyright 2017 MSolution.IO
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

package main

import (
	"context"
	"flag"
	"strconv"
	"errors"
	"database/sql"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/users"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"time"
	"github.com/trackit/trackit-server/usageReports/ec2"
	"github.com/trackit/trackit-server/usageReports/rds"
	"github.com/trackit/trackit-server/usageReports/es"
)

type (
	Resource struct {
		ResourceType string `json:"resourceType"`
		ResourceName string `json:"resourceName"`
		Tags map[string]string `json:"-"`
	}

	TagsReport map[string]map[string][]Resource
)

func taskTagsReport(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'Tags Report'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("taskTagsReport requires an integer argument")
	} else if aaId, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else {
		return prepareTagsReport(ctx, aaId)
	}
}

func prepareTagsReport(ctx context.Context, aaId int) (err error) {
	var report TagsReport
	var tx *sql.Tx
	date, _ := history.GetHistoryDate()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err := aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if user, err := users.GetUserWithId(tx, aa.UserId); err != nil {
	} else if report, err = generateTagsReport(ctx, tx, aa, user, date); err != nil {
	}
	if err == nil {
		logger.Info("Tags Report", report)
	} else {
		logger.Error("Failed to generate Tags Report", err.Error())
	}
	return
}

func generateTagsReport(ctx context.Context, tx *sql.Tx, aa aws.AwsAccount, user users.User, date time.Time) (TagsReport, error) {
	ec2Resources := getEc2Resources(ctx, tx, aa, user, date)
	rdsResources := getRdsResources(ctx, tx, aa, user, date)
	esResources := getEsResources(ctx, tx, aa, user, date)
	resources := [][]Resource{ec2Resources, rdsResources, esResources}
	report := make(TagsReport, 0)
	for _, resourceType := range resources {
		for _, resource := range resourceType {
			for key := range resource.Tags {
				report[key] = make(map[string][]Resource, 0)
			}
		}
	}
	for reportKey := range report {
		for _, resourceType := range resources {
			for _, resource := range resourceType {
				var gotKey bool = false
				for key, tag := range resource.Tags {
					if key == reportKey && tag != "" {
						report[reportKey][tag] = append(report[reportKey][tag], resource)
						gotKey = true
						break
					}
				}
				if gotKey == false {
					report[reportKey]["noTag"] = append(report[reportKey]["noTag"], resource)
				}
			}
		}
	}
	return report, nil
}

func getEc2Resources(ctx context.Context, tx *sql.Tx, aa aws.AwsAccount, user users.User, date time.Time) []Resource {
	resources := make([]Resource, 0)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	_, instances, err := ec2.GetEc2Data(ctx, ec2.Ec2QueryParams{[]string{aa.AwsIdentity}, nil, date}, user, tx)
	if err != nil {
		logger.Error("Failed to get EC2 instances", err.Error())
		return resources
	}
	for _, instance := range instances {
		resources = append(resources, Resource{
			ResourceType: "EC2",
			ResourceName: instance.Instance.Id,
			Tags: instance.Instance.Tags,
		})
	}
	return resources
}

func getRdsResources(ctx context.Context, tx *sql.Tx, aa aws.AwsAccount, user users.User, date time.Time) []Resource {
	resources := make([]Resource, 0)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	_, instances, err := rds.GetRdsData(ctx, rds.RdsQueryParams{[]string{aa.AwsIdentity}, nil, date}, user, tx)
	if err != nil {
		logger.Error("Failed to get EC2 instances", err.Error())
		return resources
	}
	for _, instance := range instances {
		resources = append(resources, Resource{
			ResourceType: "RDS",
			ResourceName: instance.Instance.DBInstanceIdentifier,
			Tags: instance.Instance.Tags,
		})
	}
	return resources
}

func getEsResources(ctx context.Context, tx *sql.Tx, aa aws.AwsAccount, user users.User, date time.Time) []Resource {
	resources := make([]Resource, 0)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	_, instances, err := es.GetEsData(ctx, es.EsQueryParams{[]string{aa.AwsIdentity}, nil, date}, user, tx)
	if err != nil {
		logger.Error("Failed to get EC2 instances", err.Error())
		return resources
	}
	for _, instance := range instances {
		resources = append(resources, Resource{
			ResourceType: "ES",
			ResourceName: instance.Domain.DomainID,
			Tags: instance.Domain.Tags,
		})
	}
	return resources
}