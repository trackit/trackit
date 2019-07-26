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
	"github.com/trackit/trackit/usageReports/riRds"
	"github.com/trackit/trackit/users"
)

const riRdsReportSheetName = "EC2 Usage Report"

var riRdsReportModule = module{
	Name:          "EC2 Usage Report",
	SheetName:     riRdsReportSheetName,
	ErrorName:     "riRdsReportError",
	GenerateSheet: generateRiRdsReportSheet,
}

// generateRiRdsReportSheet will generate a sheet with EC2 usage report
// It will get data for given AWS account and for a given date
func generateRiRdsReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return riRdsReportGenerateSheet(ctx, aas, date, tx, file)
}

func riRdsReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := riRdsReportGetData(ctx, aas, date, tx)
	if err == nil {
		return riRdsReportInsertDataInSheet(aas, file, data)
	}
	return
}

func riRdsReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []riRds.ReservationReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := riRds.ReservedInstancesQueryParams{
		AccountList: identities,
		Date:        date,
	}
	logger.Debug("Getting EC2 Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = riRds.GetReservedInstancesData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an EC2 Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func riRdsReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []riRds.ReservationReport) (err error) {
	file.NewSheet(riRdsReportSheetName)
	riRdsReportGenerateHeader(file)
	line := 3
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		instance := report.Reservation
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(instance.DBInstanceIdentifier, "B"+strconv.Itoa(line)),
			newCell(instance.DBInstanceClass, "C"+strconv.Itoa(line)),
			newCell(instance.AvailabilityZone, "D"+strconv.Itoa(line)),
			newCell(instance.MultiAZ, "E"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, riRdsReportSheetName)
		line++
	}
	return
}

func riRdsReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("ID", "B1").mergeTo("B2"),
		newCell("Type", "C1").mergeTo("C2"),
		newCell("Region", "D1").mergeTo("D2"),
		newCell("MultiAZ", "E1").mergeTo("E2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, riRdsReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35).toColumn("C"),
		newColumnWidth("D", 15).toColumn("E"),
	}
	columns.setValues(file, riRdsReportSheetName)
	return
}
