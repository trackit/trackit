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
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	odRi "github.com/trackit/trackit/onDemandToRI/ec2"
	"github.com/trackit/trackit/users"
	"strconv"
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
	return OdToRiEc2UsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func OdToRiEc2UsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := getOdToRiEc2UsageReport(ctx, aas, date, tx)
	if err == nil {
		return odToRiEc2UsageReportInsertDataInSheet(aas, file, data)
	}
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
		err = errors.New("missing AWS Account for Od to Ri EC2 Usage Reports")
		return nil, err
	}

	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
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
	if err != nil || reports == nil {
		return nil, err
	}
	return reports, nil
}

func odToRiEc2UsageReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []odRi.OdToRiEc2Report) (err error) {
	file.NewSheet(odToRiEc2UsageReportSheetName)
	odToRiEc2UsageReportInsertHeaderInSheet(file)
	line := 5
	accountLine := 0
	for _, report := range data {
		accountLine = line
		for _, instance := range report.Instances {
			instanceCells := cells{
				newCell(instance.Region, "B" + strconv.Itoa(line)),
				newCell(instance.Type, "C" + strconv.Itoa(line)),
				newCell(instance.InstanceCount, "D" + strconv.Itoa(line)),
				newCell(instance.Platform, "E" + strconv.Itoa(line)),
				newCell(instance.OnDemand.Monthly.PerUnit, "F" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.Monthly.Total, "G" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.OneYear.PerUnit, "I" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.OneYear.Total, "J" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.ThreeYears.PerUnit, "L" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.ThreeYears.Total, "M" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Monthly.PerUnit, "O" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Monthly.Total, "P" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Global.PerUnit, "R" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Global.Total, "S" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Saving.PerUnit, "U" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Saving.Total, "V" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Monthly.PerUnit, "X" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Monthly.Total, "Y" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Global.PerUnit, "AA" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Global.Total, "AB" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Saving.PerUnit, "AD" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Saving.Total, "AE" + strconv.Itoa(line)).addStyles("price"),
			}
			instanceCells.addStyles("borders", "centerText").setValues(file, odToRiEc2UsageReportSheetName)
			line++
		}
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		accountCells := cells{
			newCell(formattedAccount, "A" + strconv.Itoa(accountLine)).mergeTo("A" + strconv.Itoa(line - 1)),
			newCell(report.OnDemand.MonthlyTotal, "H" + strconv.Itoa(accountLine)).mergeTo("H" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.OnDemand.OneYearTotal, "K" + strconv.Itoa(accountLine)).mergeTo("K" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.OnDemand.ThreeYearsTotal, "N" + strconv.Itoa(accountLine)).mergeTo("N" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.OneYear.MonthlyTotal, "Q" + strconv.Itoa(accountLine)).mergeTo("Q" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.OneYear.GlobalTotal, "T" + strconv.Itoa(accountLine)).mergeTo("T" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.OneYear.SavingTotal, "W" + strconv.Itoa(accountLine)).mergeTo("W" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.ThreeYear.MonthlyTotal, "Z" + strconv.Itoa(accountLine)).mergeTo("Z" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.ThreeYear.GlobalTotal, "AC" + strconv.Itoa(accountLine)).mergeTo("AC" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.ThreeYear.SavingTotal, "AF" + strconv.Itoa(accountLine)).mergeTo("AF" + strconv.Itoa(line - 1)).addStyles("price"),
		}
		accountCells.addStyles("borders", "centerText").setValues(file, odToRiEc2UsageReportSheetName)
	}
	return
}

func odToRiEc2UsageReportInsertHeaderInSheet(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A4"),
		newCell("Region", "B1").mergeTo("B4"),
		newCell("Type", "C1").mergeTo("C4"),
		newCell("Instance Count", "D1").mergeTo("D4"),
		newCell("Platform", "E1").mergeTo("E4"),
		newCell("On Demand Cost", "F1").mergeTo("M1"),
			newCell("Monthly", "F2").mergeTo("H2"),
				newCell("Per Unit", "F3").mergeTo("F4"),
				newCell("Total", "G3").mergeTo("G4"),
				newCell("Total account", "H3").mergeTo("H4"),
			newCell("One Year", "I2").mergeTo("K2"),
				newCell("Per Unit", "I3").mergeTo("I4"),
				newCell("Total", "J3").mergeTo("J4"),
				newCell("Total account", "K3").mergeTo("K4"),
			newCell("Three Years", "L2").mergeTo("N2"),
				newCell("Per Unit", "L3").mergeTo("L4"),
				newCell("Total", "M3").mergeTo("M4"),
				newCell("Total account", "N3").mergeTo("N4"),
		newCell("Reservation One Year Cost", "O1").mergeTo("W1"),
			newCell("Monthly", "O2").mergeTo("Q2"),
				newCell("Per Unit", "O3").mergeTo("O4"),
				newCell("Total", "P3").mergeTo("P4"),
				newCell("Total account", "Q3").mergeTo("Q4"),
			newCell("One Year", "R2").mergeTo("T2"),
				newCell("Per Unit", "R3").mergeTo("R4"),
				newCell("Total", "S3").mergeTo("S4"),
				newCell("Total account", "T3").mergeTo("T4"),
			newCell("Three Years", "U2").mergeTo("W2"),
				newCell("Per Unit", "U3").mergeTo("U4"),
				newCell("Total", "V3").mergeTo("V4"),
				newCell("Total account", "W3").mergeTo("W4"),
		newCell("Reservation Three Years Cost", "X1").mergeTo("AF1"),
			newCell("Monthly", "X2").mergeTo("Z2"),
				newCell("Per Unit", "X3").mergeTo("X4"),
				newCell("Total", "Y3").mergeTo("Y4"),
				newCell("Total account", "Z3").mergeTo("Z4"),
		newCell("One Year", "AA2").mergeTo("AC2"),
				newCell("Per Unit", "AA3").mergeTo("AA4"),
				newCell("Total", "AB3").mergeTo("AB4"),
				newCell("Total account", "AC3").mergeTo("AC4"),
		newCell("Three Years", "AD2").mergeTo("AF2"),
				newCell("Per Unit", "AD3").mergeTo("AD4"),
				newCell("Total", "AE3").mergeTo("AE4"),
				newCell("Total account", "AF3").mergeTo("AF4"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, odToRiEc2UsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 12).toColumn("C"),
		newColumnWidth("D", 15),
		newColumnWidth("E", 12),
	}
	columns.setValues(file, odToRiEc2UsageReportSheetName)
}