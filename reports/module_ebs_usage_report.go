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
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/usageReports/ebs"
	"github.com/trackit/trackit/users"
)

const ebsUsageReportSheetName = "EBS Usage Report"

var ebsUsageReportModule = module{
	Name:          "EBS Usage Report",
	SheetName:     ebsUsageReportSheetName,
	ErrorName:     "ebsUsageReportError",
	GenerateSheet: generateEbsUsageReportSheet,
}

func generateEbsUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return ebsUsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func ebsUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := ebsUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return ebsUsageReportInsertDataInSheet(ctx, aas, file, data)
	} else {
		return
	}
}

func ebsUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []ebs.SnapshotReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := ebs.EbsQueryParams{
		AccountList: identities,
		Date:        date,
	}
	logger.Debug("Getting EBS Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = ebs.GetEbsData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an EBS Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func ebsUsageReportInsertDataInSheet(_ context.Context, aas []aws.AwsAccount, file *excelize.File, data []ebs.SnapshotReport) (err error) {
	file.NewSheet(ebsUsageReportSheetName)
	ebsUsageReportGenerateHeader(file)
	line := 3
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		snapshot := report.Snapshot
		date := snapshot.StartTime.Format("2006-01-02T15:04:05Z")
		tags := formatTags(snapshot.Tags)
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(snapshot.Id, "B"+strconv.Itoa(line)),
			newCell(date, "C"+strconv.Itoa(line)),
			newCell(snapshot.Region, "D"+strconv.Itoa(line)),
			newCell(snapshot.Cost, "E"+strconv.Itoa(line)).addStyles("price"),
			newCell(snapshot.Volume.Id, "F"+strconv.Itoa(line)),
			newCell(snapshot.Volume.Size, "G"+strconv.Itoa(line)),
			newCell(strings.Join(tags, ";"), "H"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, ebsUsageReportSheetName)
		line++
	}
	return
}

func ebsUsageReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("ID", "B1").mergeTo("B2"),
		newCell("Date", "C1").mergeTo("C2"),
		newCell("Region", "D1").mergeTo("D2"),
		newCell("Cost", "E1").mergeTo("E2"),
		newCell("Volume", "F1").mergeTo("G1"),
		newCell("ID", "F2"),
		newCell("Size (GigaBytes)", "G2"),
		newCell("Tags", "H1").mergeTo("H2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, ebsUsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35),
		newColumnWidth("C", 30),
		newColumnWidth("D", 20).toColumn("E"),
		newColumnWidth("F", 30),
		newColumnWidth("G", 20),
		newColumnWidth("H", 30),
	}
	columns.setValues(file, ebsUsageReportSheetName)
	return
}
