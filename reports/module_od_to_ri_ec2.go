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

package reports

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	odRi "github.com/trackit/trackit/onDemandToRI/ec2"
	"github.com/trackit/trackit/users"
	"time"
)

const odToRiEc2UsageReportSheetName = "Od To RI EC2 Usage Report"

var odToRiEc2UsageReportModule = module{
	Name:          "Od To Ri EC2 Usage Report",
	SheetName:     odToRiEc2UsageReportSheetName,
	ErrorName:     "odToRiEc2UsageReportError",
	GenerateSheet: generateOdToRiEc2UsageReportSheet,
}

// generateOdToRiEc2UsageReportSheet will generate a sheet with Od To Ri EC2 usage report
// It will get data for given AWS account and for a given date
func generateOdToRiEc2UsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	fmt.Printf("launching generate sheet for od to ri ec2 usage report\n")
	return OdToRiEc2UsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func OdToRiEc2UsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := getOdToRiEc2UsageReport(ctx, aas, date, tx)
	if err == nil {
		return odToRiEc2UsageReportInsertDataInSheet(aas, file, data)
	}
	fmt.Printf("error when getting data\n")
	return
}

func getOdToRiEc2UsageReport(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []odRi.OdToRiEc2Report, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	var dateBegin, dateEnd time.Time
	if date.IsZero() {
		dateBegin, dateEnd = history.GetHistoryDate()
	} else {
		dateBegin = date
		dateEnd = time.Date(dateBegin.Year(), dateBegin.Month()+1, 0, 23, 59, 59, 999999999, dateBegin.Location()).UTC()
	}

	if len(aas) < 1 {
		fmt.Printf("error: missing aws account, aas = %v\n", aas)
		err = errors.New("missing AWS Account for Od to Ri EC2 Usage Reports")
		return nil, err
	}

	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		fmt.Printf("error: getting user\n")
		return nil, err
	}

	parameters := odRi.RiEc2QueryParams{
		AccountList: identities,
		DateBegin:   dateBegin,
		DateEnd:     dateEnd,
	}

	logger.Debug("Getting odToRiEc2 Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
	})
	_, reports, err = odRi.GetRiEc2Report(ctx, parameters, user, tx)
	if err != nil {
		fmt.Printf("error: getting data\n")
		return nil, err
	}
	fmt.Printf("reports = %v\n", reports)
	/*if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			for _, instance := range report.Instances {
				row := formatOdToRiEc2Instance(report.Account, instance)
				data = append(data, row)
			}
		}
	}*/
	return reports, nil
}

func odToRiEc2UsageReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []odRi.OdToRiEc2Report) (err error) {
	for _, report := range data {
		fmt.Printf("report = %v", report)
	}
	return
}

/*
var odToRiEc2InstanceFormat = [][]cell{{
	newCell("", 5).addStyle(textCenter, backgroundGrey),
	newCell("", 6).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Reservation", 13).addStyle(textCenter, textBold, backgroundGrey),
}, {
	newCell("", 5).addStyle(textCenter, backgroundGrey),

	newCell("On Demand", 6).addStyle(textCenter, textBold, backgroundGrey),
	newCell("").addStyle(textCenter, textBold, backgroundGrey),
	newCell("One Year", 6).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Three Years", 6).addStyle(textCenter, textBold, backgroundGrey),
}, {
	newCell("", 5).addStyle(textCenter, backgroundGrey),
	newCell("Monthly", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("One Year", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Three Years", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Monthly", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Global", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Saving", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Monthly", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Global", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Saving", 2).addStyle(textCenter, textBold, backgroundGrey),
}, {
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Region").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Platform").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Count").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Per Unit").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatOdToRiEc2Cost(cost ec2.Cost) []cell {
	return []cell{
		newCell(cost.PerUnit),
		newCell(cost.Total),
	}
}

func formatOdToRiEc2Instance(account string, instance ec2.InstancesSpecs) []cell {
	cells := []cell{
		newCell(account),
		newCell(instance.Type),
		newCell(instance.Region),
		newCell(instance.Platform),
		newCell(instance.InstanceCount),
	}
	onDemand := make([]cell, 0, 6)
	onDemand = append(onDemand, formatOdToRiEc2Cost(instance.OnDemand.Monthly)...)
	onDemand = append(onDemand, formatOdToRiEc2Cost(instance.OnDemand.OneYear)...)
	onDemand = append(onDemand, formatOdToRiEc2Cost(instance.OnDemand.ThreeYears)...)

	cells = append(cells, onDemand...)

	reservation := make([]cell, 1, 13)
	reservation[0] = newCell(instance.Reservation.Type)
	reservation = append(reservation, formatOdToRiEc2Cost(instance.Reservation.OneYear.Monthly)...)
	reservation = append(reservation, formatOdToRiEc2Cost(instance.Reservation.OneYear.Global)...)
	reservation = append(reservation, formatOdToRiEc2Cost(instance.Reservation.OneYear.Saving)...)
	reservation = append(reservation, formatOdToRiEc2Cost(instance.Reservation.ThreeYear.Monthly)...)
	reservation = append(reservation, formatOdToRiEc2Cost(instance.Reservation.ThreeYear.Global)...)
	reservation = append(reservation, formatOdToRiEc2Cost(instance.Reservation.ThreeYear.Saving)...)
	cells = append(cells, reservation...)

	return cells
}
*/

