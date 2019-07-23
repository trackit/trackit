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

package reports

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/usageReports/riRds"
	"github.com/trackit/trackit/users"
)

var riRdsInstanceFormat = [][]cell{{
	newCell("", 6).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Offering", 1).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Recurring Charges", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Reservation", 1).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Usage", 1).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 1).addStyle(textCenter, textBold, backgroundGrey),
}, {
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
	newCell("ID").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Availability Zone").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Description").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Count").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Amount").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Frequency").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Start").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Tags").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatRiRdsInstance(report riRds.ReservationReport) []cell {
	reservation := report.Reservation
	tags := formatTags(reservation.Tags)
	chargesAmount := make([]string, len(reservation.RecurringCharges))
	chargesFrequency := make([]string, len(reservation.RecurringCharges))
	for i, charge := range reservation.RecurringCharges {
		amount := fmt.Sprintf("%f", charge.Amount)
		chargesAmount[i] = amount
		chargesFrequency[i] = charge.Frequency
	}
	return []cell{
		newCell(report.Account),
		newCell(reservation.DBInstanceIdentifier),
		newCell(reservation.DBInstanceClass),
		newCell(reservation.AvailabilityZone),
		newCell(reservation.ProductDescription),
		newCell(reservation.DBInstanceCount),
		newCell(reservation.OfferingType),
		newCell(strings.Join(chargesAmount, ";")),
		newCell(strings.Join(chargesFrequency, ";")),
		newCell(reservation.StartTime),
		newCell(strings.Join(tags, ";")),
	}
}

func getRiRdsReport(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0, len(riRdsInstanceFormat))
	for _, headerRow := range riRdsInstanceFormat {
		data = append(data, headerRow)
	}

	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}

	if len(aas) < 1 {
		err = errors.New("Missing AWS Account for Reserved Instances RDS Report")
		return
	}

	identities := getIdentities(aas)

	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}

	parameters := riRds.ReservedInstancesQueryParams{
		AccountList: identities,
		Date:        date,
	}

	logger.Debug("Getting Reserved Instances RDS Report for accounts", map[string]interface{}{
		"accounts": aas,
	})
	_, reports, err := riRds.GetReservedInstancesData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatRiRdsInstance(report)
			data = append(data, row)
		}
	}
	return
}
