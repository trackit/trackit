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

// Package riEc2 implements usage reports for Amazon EC2 reserved instances
package riEc2

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

// fetchDailyReservationsList sends in reservationInfoChan the reservations fetched from DescribeReservations
// and filled by DescribeReservations and getReservationStats.
func fetchDailyReservationsList(ctx context.Context, creds *credentials.Credentials, region string, reservationChan chan Reservation) error {
	defer close(reservationChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := ec2.New(sess)
	reservations, err := svc.DescribeReservedInstances(nil)
	if err != nil {
		logger.Error("Error when describing reservations", err.Error())
		return err
	}
	for _, reservation := range reservations.ReservedInstances {
		charges := getRecurringCharges(reservation)
		reservationChan <- Reservation{
			ReservationBase: ReservationBase{
				Id:                 aws.StringValue(reservation.ReservedInstancesId),
				Region:             region,
				AvailabilityZone:   aws.StringValue(reservation.AvailabilityZone),
				Type:               aws.StringValue(reservation.InstanceType),
				OfferingClass:      aws.StringValue(reservation.OfferingClass),
				OfferingType:       aws.StringValue(reservation.OfferingType),
				ProductDescription: aws.StringValue(reservation.ProductDescription),
				State:              aws.StringValue(reservation.State),
				Start:              aws.TimeValue(reservation.Start),
				End:                aws.TimeValue(reservation.End),
				InstanceCount:      aws.Int64Value(reservation.InstanceCount),
				Tenancy:            aws.StringValue(reservation.InstanceTenancy),
				UsagePrice:         aws.Float64Value(reservation.UsagePrice),
				RecurringCharges:   charges,
			},
			Tags: getReservationTag(reservation.Tags),
		}
	}
	return nil
}

// FetchDailyReservationsStats fetches the stats of the reserved instances of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchReservationsStats should be called every hour.
func FetchDailyReservationsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching ReservedInstances reservation stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorReservationStsSessionName)
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
	reservationChans := make([]<-chan Reservation, 0, len(regions))
	for _, region := range regions {
		reservationChan := make(chan Reservation)
		go fetchDailyReservationsList(ctx, creds, region, reservationChan)
		reservationChans = append(reservationChans, reservationChan)
	}
	reservations := make([]ReservationReport, 0)
	for reservation := range merge(reservationChans...) {
		reservations = append(reservations, ReservationReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Reservation: reservation,
		})
	}
	return importReservationsToEs(ctx, awsAccount, reservations)
}
