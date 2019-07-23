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
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/es"
)

type (
	// ReservationReport is saved in ES to have all the information of an EC2 reservation
	ReservationReport struct {
		utils.ReportBase
		Reservation Reservation `json:"reservation"`
	}

	// Reservation contains basics information of an EC2 reservation
	Reservation struct {
		Type              string   `json:"type"`
		Platform          string   `json:"platform"`
		Tenancy           string   `json:"tenancy"`
		Region            string   `json:"region"`
		AverageCoverage   float64  `json:"averageCoverage"`
		CoveredHours      float64  `json:"coveredHours"`
		OnDemandHours     float64  `json:"onDemandHours"`
		TotalRunningHours float64  `json:"totalRunningHours"`
		InstancesNames    []string `json:"instancesNames"`
	}
)

// importReportsToEs imports EC2 Coverage report in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importReportsToEs(ctx context.Context, aa aws.AwsAccount, reservations []ReservationReport) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating EC2 Coverage report for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixEC2CoverageReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return false, err
	}
	for _, reservation := range reservations {
		id, err := generateId(reservation)
		if err != nil {
			logger.Error("Error when marshaling reservation var", err.Error())
			return false, err
		}
		bp = utils.AddDocToBulkProcessor(bp, reservation, TypeEC2CoverageReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Failed to put EC2 Coverage report in ES", err.Error())
		return false, err
	}
	logger.Info("EC2 Coverage report put in ES", nil)
	return true, nil
}

func generateId(reservation ReservationReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"reservationId"`
	}{
		reservation.Account,
		reservation.ReportDate,
		uuid.NewV1().String(),
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}
