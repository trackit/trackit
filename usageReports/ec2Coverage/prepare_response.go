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
	"encoding/json"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/aws/usageReports/ec2Coverage"
	"github.com/trackit/trackit-server/errors"
)

type (
	// Structure that allow to parse ES response for EC2 Coverage Monthly report
	ResponseEc2CoverageMonthly struct {
		Accounts struct {
			Buckets []struct {
				Reservations struct {
					Hits struct {
						Hits []struct {
							Reservation ec2Coverage.ReservationReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"reservations"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// ReservationReport has all the information of an EC2 Coverage report
	ReservationReport struct {
		utils.ReportBase
		Reservation ec2Coverage.Reservation `json:"reservation"`
	}
)

// prepareResponseEc2CoverageMonthly parses the results from elasticsearch and returns an array of EC2 Coverage monthly reservations report
func prepareResponseEc2CoverageMonthly(ctx context.Context, resEc2Coverage *elastic.SearchResult) ([]ReservationReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseEc2CoverageMonthly
	reservations := make([]ReservationReport, 0)
	err := json.Unmarshal(*resEc2Coverage.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES EC2 Coverage response", err)
		return nil, errors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, reservation := range account.Reservations.Hits.Hits {
			reservations = append(reservations, ReservationReport{
				ReportBase:  reservation.Reservation.ReportBase,
				Reservation: reservation.Reservation.Reservation,
			})
		}
	}
	return reservations, nil
}
