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

	uuid "github.com/satori/go.uuid"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	utils "github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/ec2CoverageReports"
)

// importReportsToEs imports EC2 Coverage report in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importReportsToEs(ctx context.Context, aa aws.AwsAccount, reservations []ec2CoverageReports.ReservationReport) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating EC2 Coverage report for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, ec2CoverageReports.Model.IndexSuffix)
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
		bp = utils.AddDocToBulkProcessor(bp, reservation, ec2CoverageReports.Model.Type, index, id)
	}
	err = bp.Flush()
	if closeErr := bp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		logger.Error("Failed to put EC2 Coverage report in ES", err.Error())
		return false, err
	}
	logger.Info("EC2 Coverage report put in ES", nil)
	return true, nil
}

func generateId(reservation ec2CoverageReports.ReservationReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"reservationId"`
		Type       string    `json:"reportType"`
	}{
		reservation.Account,
		reservation.ReportDate,
		uuid.NewV1().String(),
		reservation.ReportType,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}
