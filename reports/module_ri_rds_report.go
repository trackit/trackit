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
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/usageReports/riRDS"
	"github.com/trackit/trackit/users"
)

const riRDSReportSheetName = "Reserved Instances RDS Report"

var riRDSReportModule = module{
	Name:          "Reserved Instances RDS Report",
	SheetName:     riRDSReportSheetName,
	ErrorName:     "riRDSReportError",
	GenerateSheet: generateRiRDSReportSheet,
}

// generateRiRDSReportSheet will generate a sheet with Ri RDS usage report
// It will get data for given AWS account and for a given date
func generateRiRDSReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return riRDSReportGenerateSheet(ctx, aas, date, tx, file)
}

func riRDSReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := riRDSReportGetData(ctx, aas, date, tx)
	if err == nil {
		return riRDSReportInsertDataInSheet(aas, file, data)
	}
	return
}

func riRDSReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []riRDS.ReservationReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := riRDS.ReservedInstancesQueryParams{
		AccountList: identities,
		Date:        date,
	}
	logger.Debug("Getting Ri RDS Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = riRDS.GetReservedInstancesData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an Ri RDS Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func riRDSReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []riRDS.ReservationReport) (err error) {
	file.NewSheet(riRDSReportSheetName)
	riRDSReportGenerateHeader(file)
	line := 4
	toLine := 4
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		instance := report.Reservation
		recurringCells := cells{
			newCell(0, "M"+strconv.Itoa(line)).addStyles("price"),
			newCell(0, "N"+strconv.Itoa(line)),
		}
		recurringCells.addStyles("borders", "centerText").setValues(file, riRDSReportSheetName)
		for currentLine, recurringCharge := range instance.RecurringCharges {
			recurringCells := cells{
				newCell(recurringCharge.Amount, "M"+strconv.Itoa(currentLine+line)).addStyles("price"),
				newCell(recurringCharge.Frequency, "N"+strconv.Itoa(currentLine+line)),
			}
			recurringCells.addStyles("borders", "centerText").setValues(file, riRDSReportSheetName)
			toLine = currentLine + line
		}
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)).mergeTo("A"+strconv.Itoa(toLine)),
			newCell(instance.DBInstanceIdentifier, "B"+strconv.Itoa(line)).mergeTo("B"+strconv.Itoa(toLine)),
			newCell(instance.DBInstanceOfferingId, "C"+strconv.Itoa(line)).mergeTo("C"+strconv.Itoa(toLine)),
			newCell(instance.AvailabilityZone, "D"+strconv.Itoa(line)).mergeTo("D"+strconv.Itoa(toLine)),
			newCell(instance.DBInstanceClass, "E"+strconv.Itoa(line)).mergeTo("E"+strconv.Itoa(toLine)),
			newCell(instance.OfferingType, "F"+strconv.Itoa(line)).mergeTo("F"+strconv.Itoa(toLine)),
			newCell(instance.DBInstanceCount, "G"+strconv.Itoa(line)).mergeTo("G"+strconv.Itoa(toLine)),
			newCell(instance.MultiAZ, "H"+strconv.Itoa(line)).mergeTo("H"+strconv.Itoa(toLine)),
			newCell(instance.State, "I"+strconv.Itoa(line)).mergeTo("I"+strconv.Itoa(toLine)),
			newCell(instance.Duration / 60 / 60 / 24 / 365, "J"+strconv.Itoa(line)).mergeTo("J"+strconv.Itoa(toLine)),
			newCell(instance.StartTime.Format("2006-01-02T15:04:05"), "K"+strconv.Itoa(line)).mergeTo("K"+strconv.Itoa(toLine)),
			newCell(instance.EndTime.Format("2006-01-02T15:04:05"), "L"+strconv.Itoa(line)).mergeTo("L"+strconv.Itoa(toLine)),
		}
		cells.addStyles("borders", "centerText").setValues(file, riRDSReportSheetName)
		line++
	}
	return
}

func riRDSReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A3"),
		newCell("Reservation", "B1").mergeTo("N1"),
		newCell("ID", "B2").mergeTo("B3"),
		newCell("Offering ID", "C2").mergeTo("C3"),
		newCell("Region", "D2").mergeTo("D3"),
		newCell("Class", "E2").mergeTo("E3"),
		newCell("Type", "F2").mergeTo("F3"),
		newCell("Count", "G2").mergeTo("G3"),
		newCell("MultiAZ", "H2").mergeTo("H3"),
		newCell("State", "I2").mergeTo("I3"),
		newCell("Duration (Year)", "J2").mergeTo("J3"),
		newCell("Start Date", "K2").mergeTo("K3"),
		newCell("End Date", "L2").mergeTo("L3"),
		newCell("Recurring Charges", "M2").mergeTo("N2"),
		newCell("Amount", "M3"),
		newCell("Frequency", "N3"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, riRDSReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 30),
		newColumnWidth("C", 40),
		newColumnWidth("D", 15).toColumn("E"),
		newColumnWidth("F", 20),
		newColumnWidth("H", 10),
		newColumnWidth("J", 16),
		newColumnWidth("K", 25),
		newColumnWidth("L", 25),
		newColumnWidth("M", 10),
		newColumnWidth("N", 13),
	}
	columns.setValues(file, riRDSReportSheetName)
	return
}
