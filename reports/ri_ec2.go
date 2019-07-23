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

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/riEc2"
	"github.com/trackit/trackit-server/users"
)

var riEc2InstanceFormat = [][]cell{{
	newCell("", 7).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Offering", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Recurring Charges", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Reservation", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 1).addStyle(textCenter, textBold, backgroundGrey),
}, {
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
	newCell("ID").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Region").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Description").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Count").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Tenancy").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Class").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Amount").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Frequency").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Start").addStyle(textCenter, textBold, backgroundGrey),
	newCell("End").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Tags").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatRiEc2Instance(report riEc2.ReservationReport) []cell {
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
		newCell(reservation.Id),
		newCell(reservation.Type),
		newCell(reservation.Region),
		newCell(reservation.ProductDescription),
		newCell(reservation.InstanceCount),
		newCell(reservation.Tenancy),
		newCell(reservation.OfferingClass),
		newCell(reservation.OfferingType),
		newCell(strings.Join(chargesAmount, ";")),
		newCell(strings.Join(chargesFrequency, ";")),
		newCell(reservation.Start),
		newCell(reservation.End),
		newCell(strings.Join(tags, ";")),
	}
}

func getRiEc2Report(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0, len(riEc2InstanceFormat))
	for _, headerRow := range riEc2InstanceFormat {
		data = append(data, headerRow)
	}

	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}

	if len(aas) < 1 {
		err = errors.New("Missing AWS Account for Reserved Instances EC2 Report")
		return
	}

	identities := getIdentities(aas)

	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}

	parameters := riEc2.ReservedInstancesQueryParams{
		AccountList: identities,
		Date:        date,
	}

	logger.Debug("Getting Reserved Instances EC2 Report for accounts", map[string]interface{}{
		"accounts": aas,
	})
	_, reports, err := riEc2.GetReservedInstancesData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatRiEc2Instance(report)
			data = append(data, row)
		}
	}
	return
}
