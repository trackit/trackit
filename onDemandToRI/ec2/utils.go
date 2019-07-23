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

package onDemandToRiEc2

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/usageReports/ec2"
	"github.com/trackit/trackit/usageReports/riEc2"
)

var (
	HoursPerMonth = 730.0
)

// getRegionName takes an availability zone or a region name and returns a region name
func getRegionName(az string) string {
	if _, err := strconv.Atoi(string(az[len(az)-1])); err == nil {
		// The "az" finishes by a number, so it's a region name
		return az
	}
	return az[:len(az)-1]
}

// instanceMatchReservation takes an ec2.InstanceReport and an riEc2.ReservationReport
// It returns true if the InstanceReport matches the ReservationReport
func instanceMatchReservation(instanceReport ec2.InstanceReport, reservationReport riEc2.ReservationReport) bool {
	if (getRegionName(instanceReport.Instance.Region) == reservationReport.Reservation.Region ||
		instanceReport.Instance.Region == reservationReport.Reservation.AvailabilityZone) &&
		instanceReport.Instance.Type == reservationReport.Reservation.Type &&
		instanceReport.Instance.Platform == reservationReport.Reservation.ProductDescription {
		return true
	}
	return false
}

// instanceMatchSpecs takes an ec2.InstanceReport and an InstancesSpecs
// it returns true if the InstanceReport matches the InstancesSpecs
func instanceMatchSpecs(instanceReport ec2.InstanceReport, specs InstancesSpecs) bool {
	if getRegionName(instanceReport.Instance.Region) == specs.Region &&
		instanceReport.Instance.Type == specs.Type && instanceReport.Instance.Platform == specs.Platform {
		return true
	}
	return false
}

// getMonthlyCostPerUnit returns the monthly cost based on the hourlyCost
// it returns 0.0 if the hourlyCost is -1.0 (which means the pricing term does not exist)
func getMonthlyCostPerUnit(hourlyCost float64) float64 {
	if hourlyCost != -1.0 {
		return hourlyCost * HoursPerMonth
	}
	return 0.0
}

// IngestOdToRiEc2Result saves a OdToRiEc2Report into elasticsearch
func IngestOdToRiEc2Result(ctx context.Context, aa aws.AwsAccount, report OdToRiEc2Report) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Saving od to ri ec2 result for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	client := es.Client
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
	}{
		report.Account,
		report.ReportDate,
	})
	if err != nil {
		logger.Error("Error when marshaling instance var", err.Error())
		return err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixOdToRiEC2Report)
	if res, err := client.
		Index().
		Index(index).
		Type(TypeOdToRiEC2Report).
		BodyJson(report).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting od to ri ec2 result in ES", err.Error())
		return err
	} else {
		logger.Info("od to ri ec2 result put in ES", *res)
	}
	return nil
}
