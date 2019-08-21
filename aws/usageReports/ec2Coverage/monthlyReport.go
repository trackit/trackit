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

package ec2Coverage

import (
	"context"
	"github.com/trackit/trackit/pagination"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/usageReports/ec2"
	"github.com/trackit/trackit/users"
)

// geEc2CoverageReport get EC2 coverage data from the AWS Cost Explorer
func getEc2CoverageReport(creds *credentials.Credentials, start time.Time) (*costexplorer.GetReservationCoverageOutput, error) {
	end := time.Date(start.Year(), start.Month()+1, 1, 0, 0, 0, 0, start.Location())
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(""),
	}))
	svc := costexplorer.New(sess)
	input := costexplorer.GetReservationCoverageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start.Format("2006-01-02")),
			End:   aws.String(end.Format("2006-01-02")),
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("INSTANCE_TYPE"),
				Type: aws.String("DIMENSION"),
			}, {
				Key:  aws.String("PLATFORM"),
				Type: aws.String("DIMENSION"),
			}, {
				Key:  aws.String("TENANCY"),
				Type: aws.String("DIMENSION"),
			}, {
				Key:  aws.String("REGION"),
				Type: aws.String("DIMENSION"),
			},
		},
	}
	reservations, err := svc.GetReservationCoverage(&input)
	if err != nil {
		return nil, err
	} else {
		return reservations, nil
	}
}

// getInstancesNames return an array of the names of the instances which are linked to a reservation
func getInstancesNames(report Reservation, instances []ec2.InstanceReport) []string {
	names := make([]string, 0)
	for _, instance := range instances {
		if name, ok := instance.Instance.Tags["Name"]; ok && instance.Instance.Type == report.Type &&
			strings.Contains(instance.Instance.Region, report.Region) && instance.Instance.Purchasing == "on demand" {
			names = append(names, name)
		}
	}
	return names
}

func generateEc2ReservationReport(ctx context.Context, aa taws.AwsAccount, start time.Time, instances []ec2.InstanceReport, reservations *costexplorer.GetReservationCoverageOutput) (bool, error) {
	reports := make([]ReservationReport, 0)
	for _, covByTime := range reservations.CoveragesByTime {
		for _, group := range covByTime.Groups {
			report := Reservation{
				Type:     aws.StringValue(group.Attributes["instanceType"]),
				Platform: aws.StringValue(group.Attributes["platform"]),
				Tenancy:  aws.StringValue(group.Attributes["tenancy"]),
				Region:   aws.StringValue(group.Attributes["region"]),
			}
			cov := group.Coverage.CoverageHours
			if value, err := strconv.ParseFloat(aws.StringValue(cov.CoverageHoursPercentage), 64); err == nil {
				report.AverageCoverage = value
			}
			if value, err := strconv.ParseFloat(aws.StringValue(cov.ReservedHours), 64); err == nil {
				report.CoveredHours = value
			}
			if value, err := strconv.ParseFloat(aws.StringValue(cov.OnDemandHours), 64); err == nil {
				report.OnDemandHours = value
			}
			if value, err := strconv.ParseFloat(aws.StringValue(cov.TotalRunningHours), 64); err == nil {
				report.TotalRunningHours = value
			}
			report.InstancesNames = getInstancesNames(report, instances)
			reports = append(reports, ReservationReport{
				ReportBase: utils.ReportBase{
					Account:    aa.AwsIdentity,
					ReportType: "monthly",
					ReportDate: start,
				},
				Reservation: report,
			})
		}
	}
	return importReportsToEs(ctx, aa, reports)
}

// PutEc2MonthlyCoverageReport puts a monthly report of EC2 reservation in ES
func PutEc2MonthlyCoverageReport(ctx context.Context, aa taws.AwsAccount, start, end time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting EC2 Coverage monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    start.Format("2006-01-02T15:04:05Z"),
		"endDate":      end.Format("2006-01-02T15:04:05Z"),
	})
	if already, err := utils.CheckMonthlyReportExists(ctx, start, aa, IndexPrefixEC2CoverageReport); err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an EC2 Coverage monthly report", nil)
		return false, nil
	} else if tx, err := db.Db.Begin(); err != nil {
		logger.Error("Failed to begin Tx for DB", err.Error())
		return false, err
	} else if user, err := users.GetUserWithId(tx, aa.UserId); err != nil {
		logger.Error("Failed to get User by Id", err.Error())
		return false, err
	} else if _, instances, err := ec2.GetEc2Data(ctx, ec2.Ec2QueryParams{
		AccountList: []string{aa.AwsIdentity},
		IndexList:   nil,
		Date:        start,
		Pagination:  pagination.NewPagination(nil),
	}, user, tx); err != nil {
		return false, err
	} else if creds, err := taws.GetTemporaryCredentials(aa, "monitor-coverage"); err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return false, err
	} else if reservations, err := getEc2CoverageReport(creds, start); err != nil {
		logger.Error("Failed to get EC2 Coverage report from Cost Explorer", err.Error())
		return false, err
	} else {
		return generateEc2ReservationReport(ctx, aa, start, instances, reservations)
	}
}
