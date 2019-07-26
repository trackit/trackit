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
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/usageReports/riEc2"
	"github.com/trackit/trackit/users"
)

const riEc2ReportSheetName = "Reserved Instance Report"

var riEc2ReportModule = module{
	Name:          "EC2  Report",
	SheetName:     riEc2ReportSheetName,
	ErrorName:     "riec2ReportError",
	GenerateSheet: generateRiEc2ReportSheet,
}

// generateEc2ReportSheet will generate a sheet with EC2 usage report
// It will get data for given AWS account and for a given date
func generateRiEc2ReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return riec2ReportGenerateSheet(ctx, aas, date, tx, file)
}

func riec2ReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := riec2ReportGetData(ctx, aas, date, tx)
	if err == nil {
		return riec2ReportInsertDataInSheet(aas, file, data)
	}
	return
}

func riec2ReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []riEc2.ReservationReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := riEc2.ReservedInstancesQueryParams{
		AccountList: identities,
		Date:        date,
	}
	logger.Debug("Getting EC2  Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = riEc2.GetReservedInstancesData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an EC2  Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func riec2ReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []riEc2.ReservationReport) (err error) {
	file.NewSheet(riEc2ReportSheetName)
	riec2ReportGenerateHeader(file)
	line := 4
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		instance := report.Reservation
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(instance.Id, "B"+strconv.Itoa(line)),
			newCell(instance.Type, "C"+strconv.Itoa(line)),
			newCell(instance.Region, "D"+strconv.Itoa(line)),
			newCell(instance.State, "E"+strconv.Itoa(line)),
			newCell(instance.OfferingClass, "F"+strconv.Itoa(line)),
			newCell(instance.OfferingType, "G"+strconv.Itoa(line)),
			newCell(instance.InstanceCount, "H"+strconv.Itoa(line)),
			newCell(instance.Start.Format("2006-01-02T15:04:05"), "I"+strconv.Itoa(line)),
			newCell(instance.End.Format("2006-01-02T15:04:05"), "J"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, riEc2ReportSheetName)
		line++
	}
	return
}

func riec2ReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A3"),
		newCell("Reservation", "B1").mergeTo("L1"),
		newCell("ID", "B2").mergeTo("B3"),
		newCell("Type", "C2").mergeTo("C3"),
		newCell("Region", "D2").mergeTo("D3"),
		newCell("State", "E2").mergeTo("E3"),
		newCell("Offering Class", "F2").mergeTo("F3"),
		newCell("Offering Type", "G2").mergeTo("G3"),
		newCell("Amount", "H2").mergeTo("H3"),
		newCell("Date Reservations", "I2").mergeTo("J2"),
		newCell("Start", "I3"),
		newCell("End", "J3"),
		newCell("Recurring Charges", "K2").mergeTo("L2"),
		newCell("Amount", "K3"),
		newCell("Frequency", "L3"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, riEc2ReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 37),
		newColumnWidth("C", 15).toColumn("D"),
		newColumnWidth("E", 12.5).toColumn("G"),
		newColumnWidth("H", 7),
		newColumnWidth("I", 25).toColumn("J"),
	}
	columns.setValues(file, riEc2ReportSheetName)
	return
}
