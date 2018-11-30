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
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/trackit/trackit-server/es"
)

const MonitorReservationStsSessionName = "monitor-reservation"

type (
	// ReservationReport is saved in ES to have all the information of an ReservedInstances reservation
	ReservationReport struct {
		utils.ReportBase
		Reservation Reservation `json:"reservation"`
	}

	// ReservationBase contains basics information of an ReservedInstances reservation
	ReservationBase struct {
		Id              string    `json:"id"`
		Region          string    `json:"region"`
		Type            string    `json:"type"`
		FixedPrice      float64   `json:"fixed_price"`
		UsagePrice      float64   `json:"usage_price"`
		UsageDuration   int64     `json:"usage_duration"`
		Start           time.Time `json:"start"`
		End             time.Time `json:"end"`
		InstanceCount   int64     `json:"instance_count"`
		InstanceTenancy string    `json:"instance_tenancy"`
	}

	// Reservation contains all the information of an ReservedInstances reservation
	Reservation struct {
		ReservationBase
		Tags  []utils.Tag        `json:"tags"`
	}
)

// importReservationsToEs imports ReservedInstances reservations in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importReservationsToEs(ctx context.Context, aa taws.AwsAccount, reservations []ReservationReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating ReservedInstances reservations for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixReservedInstancesReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, reservation := range reservations {
		id, err := generateId(reservation)
		if err != nil {
			logger.Error("Error when marshaling reservation var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, reservation, TypeReservedInstancesReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put ReservedInstances reservations in ES", err.Error())
		return err
	}
	logger.Info("ReservedInstances reservations put in ES", nil)
	return nil
}

func generateId(reservation ReservationReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		reservation.Account,
		reservation.ReportDate,
		reservation.Reservation.Id,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func merge(cs ...<-chan Reservation) <-chan Reservation {
	var wg sync.WaitGroup
	out := make(chan Reservation)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Reservation) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
