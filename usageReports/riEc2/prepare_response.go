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

package riEc2

import (
	"context"
	"encoding/json"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/aws/usageReports/riEc2"
)

type (

	// Structure that allow to parse ES response for ReservedInstances Daily reservations
	ResponseReservedInstancesDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time         string `json:"key_as_string"`
						Reservations struct {
							Hits struct {
								Hits []struct {
									Reservation riEc2.ReservationReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"reservations"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// ReservationReport has all the information of an ReservedInstances reservation report
	ReservationReport struct {
		utils.ReportBase
		Reservation Reservation `json:"reservation"`
	}

	// Reservation contains the information of an ReservedInstances reservation
	Reservation struct {
		riEc2.ReservationBase
		Tags map[string]string `json:"tags"`
	}
)

func getReservedInstancesReportResponse(oldReservation riEc2.ReservationReport) ReservationReport {
	tags := make(map[string]string, 0)
	for _, tag := range oldReservation.Reservation.Tags {
		tags[tag.Key] = tag.Value
	}
	newReservation := ReservationReport{
		ReportBase: oldReservation.ReportBase,
		Reservation: Reservation{
			ReservationBase: oldReservation.Reservation.ReservationBase,
			Tags:            tags,
		},
	}
	return newReservation
}

// prepareResponseReservedInstancesDaily parses the results from elasticsearch and returns an array of ReservedInstances daily reservations report
func prepareResponseReservedInstancesDaily(ctx context.Context, resReservedInstances *elastic.SearchResult) ([]ReservationReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedReservedInstances ResponseReservedInstancesDaily
	reservations := make([]ReservationReport, 0)
	err := json.Unmarshal(*resReservedInstances.Aggregations["accounts"], &parsedReservedInstances.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES ReservedInstances response", err)
		return nil, err
	}
	for _, account := range parsedReservedInstances.Accounts.Buckets {
		var lastDate = ""
		for _, date := range account.Dates.Buckets {
			if date.Time > lastDate {
				lastDate = date.Time
			}
		}
		for _, date := range account.Dates.Buckets {
			if date.Time == lastDate {
				for _, reservation := range date.Reservations.Hits.Hits {
					reservations = append(reservations, getReservedInstancesReportResponse(reservation.Reservation))
				}
			}
		}
	}
	return reservations, nil
}
