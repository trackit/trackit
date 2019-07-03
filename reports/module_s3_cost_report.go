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

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/s3/costs"
)

const s3CostReportSheetName = "S3 Cost Report"

var s3CostReportModule = module{
	Name:          "S3 Cost Report",
	SheetName:     s3CostReportSheetName,
	ErrorName:     "s3CostReportError",
	GenerateSheet: generateS3CostReportSheet,
}

func generateS3CostReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return s3CostReportGenerateSheet(ctx, aas, date, tx, file)
}

func s3CostReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := s3CostReportGetData(ctx, aas, date, tx)
	if err == nil {
		return s3CostReportInsertDataInSheet(ctx, aas, file, data)
	} else {
		return
	}
}

func s3CostReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports costs.BucketsInfo, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	identities := getAwsIdentities(aas)

	/*user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {												//user isn't use right now
		return
	}*/

	parameters := costs.S3QueryParams{
		AccountList: identities,
		DateBegin:        date,
		DateEnd:          time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location()).UTC(),
	}

	logger.Debug("Getting S3 Cost Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = costs.GetS3CostData(ctx, parameters)
	if err != nil {
		logger.Error("An error occurred while generating an S3 Cost Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func s3CostReportInsertDataInSheet(_ context.Context, aas []aws.AwsAccount, file *excelize.File, data costs.BucketsInfo) (err error) {
	file.NewSheet(s3CostReportSheetName)
	s3CostReportGenerateHeader(file)
	line := 3
	for idx, report := range data {
		account := getAwsAccount(idx, aas)
		formattedAccount := idx
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		//instance := report.Instance
		/*name := ""
		if value, ok := instance.Tags["Name"]; ok {
			name = value
		}*/
		//tags := formatTags(instance.Tags)
		total := make(map[string]float64)
		total["Storage"] = report.StorageCost
		total["BandwidthCost"] = report.BandwidthCost
		total["RequestCost"] = report.RequestsCost
		cells := cells{
			//newCell(instance.Id, "B"+strconv.Itoa(line)), // ACCOUNT
			newCell(formattedAccount, "B"+strconv.Itoa(line)),
			newCell(report.GbMonth, "C"+strconv.Itoa(line)).addStyles("Gb"),
			newCell(report.StorageCost, "D"+strconv.Itoa(line)).addStyles("price"),
			newCell(report.BandwidthCost, "E"+strconv.Itoa(line)).addStyles("price"),
			newCell(report.RequestsCost, "F"+strconv.Itoa(line)).addStyles("price"),
			newCell(getTotal(total), "G"+strconv.Itoa(line)).addStyles("price"),
			newCell(report.DataIn, "H"+strconv.Itoa(line)),
			newCell(report.DataOut, "I"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, s3CostReportSheetName)
		line++
	}
	return
}

func s3CostReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("Name", "B1").mergeTo("B2"),
		newCell("Billable Size (GB)", "C1").mergeTo("C2"),
		newCell("Cost", "D1").mergeTo("G1"),
		newCell("Storage", "D2") ,
		newCell("Bandwidth", "E2"),
		newCell("Requests", "F2"),
		newCell("Total", "G2"),
		newCell("Data transfers", "H1").mergeTo("I1"),
		newCell("In (GB)", "H2"),
		newCell("Out (GB)", "I2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, s3CostReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 50),
		newColumnWidth("C", 20).toColumn("I"),
	}
	columns.setValues(file, s3CostReportSheetName)
	return
}
