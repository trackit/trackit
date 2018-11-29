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

package reservedInstances

import (
	"context"
	"encoding/json"
	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/aws/usageReports/reservedInstances"
)

type (

	// Structure that allow to parse ES response for costs
	ResponseCost struct {
		Accounts struct {
			Buckets []struct {
				Key       string `json:"key"`
				Reservations struct {
					Buckets []struct {
						Key  string `json:"key"`
						Cost struct {
							Value float64 `json:"value"`
						} `json:"cost"`
					} `json:"buckets"`
				} `json:"reservations"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for ReservedReservations Monthly reservations
	ResponseReservedReservationsMonthly struct {
		Accounts struct {
			Buckets []struct {
				Reservations struct {
					Hits struct {
						Hits []struct {
							Reservation reservedInstances.ReservationReport `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"reservations"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// Structure that allow to parse ES response for ReservedReservations Daily reservations
	ResponseReservedReservationsDaily struct {
		Accounts struct {
			Buckets []struct {
				Dates struct {
					Buckets []struct {
						Time      string `json:"key_as_string"`
						Reservations struct {
							Hits struct {
								Hits []struct {
									Reservation reservedInstances.ReservationReport `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"reservations"`
					} `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		} `json:"accounts"`
	}

	// ReservationReport has all the information of an ReservedReservations reservation report
	ReservationReport struct {
		utils.ReportBase
		Reservation Reservation `json:"reservation"`
	}

	// Reservation contains the information of an ReservedReservations reservation
	Reservation struct {
		reservedInstances.ReservationBase
		Tags  map[string]string  `json:"tags"`
	}
)

func getReservedReservationsInstanceReportResponse(oldReservation reservedInstances.ReservationReport) ReservationReport {
	tags := make(map[string]string, 0)
	for _, tag := range oldReservation.Reservation.Tags {
		tags[tag.Key] = tag.Value
	}
	newReservation := ReservationReport{
		ReportBase: oldReservation.ReportBase,
		Reservation: Reservation{
			ReservationBase: oldReservation.Reservation.ReservationBase,
			Tags:         tags,
		},
	}
	return newReservation
}

// prepareResponseReservedReservationsDaily parses the results from elasticsearch and returns an array of ReservedReservations daily reservations report
func prepareResponseReservedReservationsDaily(ctx context.Context, resReservedReservations *elastic.SearchResult, resCost *elastic.SearchResult) ([]ReservationReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedReservedReservations ResponseReservedReservationsDaily
	var parsedCost ResponseCost
	reservations := make([]ReservationReport, 0)
	err := json.Unmarshal(*resReservedReservations.Aggregations["accounts"], &parsedReservedReservations.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES ReservedReservations response", err)
		return nil, err
	}
	if resCost != nil {
		err = json.Unmarshal(*resCost.Aggregations["accounts"], &parsedCost.Accounts)
		if err != nil {
			logger.Error("Error while unmarshaling ES cost response", err)
		}
	}
	for _, account := range parsedReservedReservations.Accounts.Buckets {
		var lastDate = ""
		for _, date := range account.Dates.Buckets {
			if date.Time > lastDate {
				lastDate = date.Time
			}
		}
		for _, date := range account.Dates.Buckets {
			if date.Time == lastDate {
				for _, reservation := range date.Reservations.Hits.Hits {
					reservations = append(reservations, getReservedReservationsInstanceReportResponse(reservation.Reservation))
				}
			}
		}
	}
	return reservations, nil
}

// prepareResponseReservedReservationsMonthly parses the results from elasticsearch and returns an array of ReservedReservations monthly reservations report
func prepareResponseReservedReservationsMonthly(ctx context.Context, resReservedReservations *elastic.SearchResult) ([]ReservationReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var response ResponseReservedReservationsMonthly
	reservations := make([]ReservationReport, 0)
	err := json.Unmarshal(*resReservedReservations.Aggregations["accounts"], &response.Accounts)
	if err != nil {
		logger.Error("Error while unmarshaling ES ReservedReservations response", err)
		return nil, errors.GetErrorMessage(ctx, err)
	}
	for _, account := range response.Accounts.Buckets {
		for _, reservation := range account.Reservations.Hits.Hits {
			reservations = append(reservations, getReservedReservationsInstanceReportResponse(reservation.Reservation))
		}
	}
	return reservations, nil
}
