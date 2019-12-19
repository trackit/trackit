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
				newCell(instance.Region, "K" + strconv.Itoa(line)),
				newCell(instance.Region, "L" + strconv.Itoa(line)),
				newCell(instance.Region, "M" + strconv.Itoa(line)),
				newCell(instance.Region, "N" + strconv.Itoa(line)),
				newCell(instance.OnDemand.Monthly.PerUnit, "O" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.Monthly.Total, "P" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.OneYear.PerUnit, "Q" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.OneYear.Total, "R" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.ThreeYears.PerUnit, "S" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.OnDemand.ThreeYears.Total, "T" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Monthly.PerUnit, "U" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Monthly.Total, "V" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Global.PerUnit, "W" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Global.Total, "X" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Saving.PerUnit, "Y" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.OneYear.Saving.Total, "Z" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Monthly.PerUnit, "AA" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Monthly.Total, "AB" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Global.PerUnit, "AC" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Global.Total, "AD" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Saving.PerUnit, "AE" + strconv.Itoa(line)).addStyles("price"),
				newCell(instance.Reservation.ThreeYear.Saving.Total, "AF" + strconv.Itoa(line)).addStyles("price"),
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
			newCell(report.OnDemand.MonthlyTotal, "B" + strconv.Itoa(accountLine)).mergeTo("B" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.OnDemand.OneYearTotal, "C" + strconv.Itoa(accountLine)).mergeTo("C" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.OnDemand.ThreeYearsTotal, "D" + strconv.Itoa(accountLine)).mergeTo("D" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.OneYear.MonthlyTotal, "E" + strconv.Itoa(accountLine)).mergeTo("E" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.OneYear.GlobalTotal, "F" + strconv.Itoa(accountLine)).mergeTo("F" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.OneYear.SavingTotal, "G" + strconv.Itoa(accountLine)).mergeTo("G" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.ThreeYear.MonthlyTotal, "H" + strconv.Itoa(accountLine)).mergeTo("H" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.ThreeYear.GlobalTotal, "I" + strconv.Itoa(accountLine)).mergeTo("I" + strconv.Itoa(line - 1)).addStyles("price"),
			newCell(report.Reservation.ThreeYear.SavingTotal, "J" + strconv.Itoa(accountLine)).mergeTo("J" + strconv.Itoa(line - 1)).addStyles("price"),
		}
		accountCells.addStyles("borders", "centerText").setValues(file, odToRiEc2UsageReportSheetName)
	}
	return
}

func odToRiEc2UsageReportInsertHeaderInSheet(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A4"),
		newCell("On Demand Total Cost", "B1").mergeTo("D1"),
		newCell("Monthly", "B2").mergeTo("B4"),
		newCell("One Year", "C2").mergeTo("C4"),
		newCell("Three Years", "D2").mergeTo("D4"),
		newCell("Reservation", "E1").mergeTo("J1"),
		newCell("One Year", "E2").mergeTo("G2"),
		newCell("Monthly", "E3").mergeTo("E4"),
		newCell("Global", "F3").mergeTo("F4"),
		newCell("Monthly", "G3").mergeTo("G4"),
		newCell("Three Years", "H2").mergeTo("J2"),
		newCell("Monthly", "H3").mergeTo("H4"),
		newCell("Global", "I3").mergeTo("I4"),
		newCell("Monthly", "J3").mergeTo("J4"),
		newCell("Region", "K1").mergeTo("K4"),
		newCell("Type", "L1").mergeTo("L4"),
		newCell("Instance Count", "M1").mergeTo("M4"),
		newCell("Platform", "N1").mergeTo("N4"),
		newCell("On Demand Cost", "O1").mergeTo("T1"),
		newCell("Monthly", "O2").mergeTo("P2"),
		newCell("Per Unit", "O3").mergeTo("O4"),
		newCell("Total", "P3").mergeTo("P4"),
		newCell("One Year", "Q2").mergeTo("R2"),
		newCell("Per Unit", "Q3").mergeTo("Q4"),
		newCell("Total", "R3").mergeTo("R4"),
		newCell("Three Years", "S2").mergeTo("T2"),
		newCell("Per Unit", "S3").mergeTo("S4"),
		newCell("Total", "T3").mergeTo("T4"),
		newCell("Reservation One Year Cost", "U1").mergeTo("Z1"),
		newCell("Monthly", "U2").mergeTo("V2"),
		newCell("Per Unit", "U3").mergeTo("U4"),
		newCell("Total", "V3").mergeTo("V4"),
		newCell("One Year", "W2").mergeTo("X2"),
		newCell("Per Unit", "W3").mergeTo("W4"),
		newCell("Total", "X3").mergeTo("X4"),
		newCell("Three Years", "Y2").mergeTo("Z2"),
		newCell("Per Unit", "Y3").mergeTo("Y4"),
		newCell("Total", "Z3").mergeTo("Z4"),
		newCell("Reservation Three Years Cost", "AA1").mergeTo("AF1"),
		newCell("Monthly", "AA2").mergeTo("AB2"),
		newCell("Per Unit", "AA3").mergeTo("AA4"),
		newCell("Total", "AB3").mergeTo("AB4"),
		newCell("One Year", "AC2").mergeTo("AD2"),
		newCell("Per Unit", "AC3").mergeTo("AC4"),
		newCell("Total", "AD3").mergeTo("AD4"),
		newCell("Three Years", "AE2").mergeTo("AF2"),
		newCell("Per Unit", "AE3").mergeTo("AE4"),
		newCell("Total", "AF3").mergeTo("AF4"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, odToRiEc2UsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 15).toColumn("J"),
		newColumnWidth("K", 10).toColumn("L"),
		newColumnWidth("N", 10),
	}
	columns.setValues(file, odToRiEc2UsageReportSheetName)
}