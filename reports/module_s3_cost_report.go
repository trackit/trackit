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
	"fmt"
	"context"
	"database/sql"
	"strconv"
	"strings"
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

func generateS3CostReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, _ *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return s3CostReportGenerateSheet(ctx, aas, date, file)
}

func s3CostReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, file *excelize.File) (err error) {
	data, err := s3CostReportGetData(ctx, aas, date)
	if err == nil {
		return s3CostReportInsertDataInSheet(ctx, file, data)
	} else {
		return
	}
}

func s3CostReportGetData(ctx context.Context, aas []aws.AwsAccount, dateBegin time.Time) (reports map[aws.AwsAccount]costs.BucketsInfo, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	reports = make(map[aws.AwsAccount]costs.BucketsInfo, len(aas))
	dateEnd := time.Date(dateBegin.Year(), dateBegin.Month()+1, 0, 23, 59, 59, 999999999, dateBegin.Location()).UTC()
	for _, aa := range aas {
		parameters := costs.S3QueryParams{
			AccountList: []string{aa.AwsIdentity},
			DateBegin:   dateBegin,
			DateEnd:     dateEnd,
		}
		logger.Debug("Getting S3 Cost Report for accounts", map[string]interface{}{
			"accounts": aa,
			"date":     dateBegin,
		})
		_, reports[aa], err = costs.GetS3CostData(ctx, parameters)
		if err != nil {
			logger.Error("An error occurred while generating an S3 Cost Report", map[string]interface{}{
				"error":    err,
				"accounts": aa,
				"date":     dateBegin,
			})
			return reports, err
		}
	}
	return
}

func s3CostReportInsertDataInSheet(_ context.Context, file *excelize.File, data map[aws.AwsAccount]costs.BucketsInfo) (err error) {
	file.NewSheet(s3CostReportSheetName)
	s3CostReportGenerateHeader(file)
	line := 3
	for acc, report := range data {
		for name, values := range report {
			costCols := []string{
				"D" + strconv.Itoa(line),
				"E" + strconv.Itoa(line),
				"F" + strconv.Itoa(line),
			}
			totalCostFormula := fmt.Sprintf("SUM(%s)", strings.Join(costCols, ","))
			cells := cells{
				newCell(formatAwsAccount(acc), "A" + strconv.Itoa(line)),
				newCell(name, "B" + strconv.Itoa(line)),
				newCell(values.GbMonth, "C" + strconv.Itoa(line)),
				newCell(values.StorageCost, "D" + strconv.Itoa(line)).addStyles("price"),
				newCell(values.BandwidthCost, "E" + strconv.Itoa(line)).addStyles("price"),
				newCell(values.RequestsCost, "F" + strconv.Itoa(line)).addStyles("price"),
				newFormula(totalCostFormula, "G" + strconv.Itoa(line)).addStyles("price"),
				newCell(values.DataIn, "H" + strconv.Itoa(line)),
				newCell(values.DataOut, "I" + strconv.Itoa(line)),
			}
			cells.addStyles("borders", "centerText").setValues(file, s3CostReportSheetName)
			line++
		}
	}
	return
}

func s3CostReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("Name", "B1").mergeTo("B2"),
		newCell("Billable Size (GigaBytes)", "C1").mergeTo("C2"),
		newCell("Cost", "D1").mergeTo("G1"),
		newCell("Storage", "D2"),
		newCell("Bandwidth", "E2"),
		newCell("Requests", "F2"),
		newCell("Total", "G2"),
		newCell("Data Transfers (GigaBytes)", "H1").mergeTo("I1"),
		newCell("In", "H2"),
		newCell("Out", "I2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, s3CostReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 50),
		newColumnWidth("C", 20),
		newColumnWidth("D", 12.5).toColumn("G"),
		newColumnWidth("H", 20).toColumn("I"),
	}
	columns.setValues(file, s3CostReportSheetName)
	return
}
