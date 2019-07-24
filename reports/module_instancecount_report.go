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
	"github.com/trackit/trackit-server/costs/diff"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/instanceCount"
	"github.com/trackit/trackit-server/users"
)

const instanceCountReportDetailledSheetName = "Instance Count Report (Detailed)"

var instanceCountUsageReportModule = module{
	Name:          "Instance Count Report",
	SheetName:     instanceCountReportDetailledSheetName,
	ErrorName:     "instanceCountUsageReportError",
	GenerateSheet: generateInstanceCountUsageReportSheet,
}

func generateInstanceCountUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	var dateRange diff.DateRange
	if date.IsZero() {
		dateRange.Begin, dateRange.End = history.GetHistoryDate()
	} else {
		dateRange = diff.DateRange{
			Begin: date,
			End:   time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location()).UTC(),
		}
	}
	return instanceCountUsageReportGenerateSheet(ctx, aas, dateRange, tx, file)
}

func instanceCountUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date diff.DateRange, tx *sql.Tx, file *excelize.File) (err error) {
	data, dates, err := instanceCountUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return instanceCountUsageReportInsertDataInSheet(ctx, aas, file, data, dates)
	} else {
		return
	}
}

func instanceCountUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date diff.DateRange, tx *sql.Tx) (reports []instanceCount.InstanceCountReport, dates map[time.Time]int, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := instanceCount.InstanceCountQueryParams{
		AccountList: identities,
		Date:        date.Begin,
	}
	logger.Debug("Getting InstanceCount Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date.Begin,
	})
	_, reports, err = instanceCount.GetInstanceCountData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an InstanceCount Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date.Begin,
		})
	}
	dates = getAllDateInstanceCount(date)
	return
}

// getAllDateInstanceCount get all the hours and the column position for a DateRange and put it in a map[time.Time]int
func getAllDateInstanceCount(date diff.DateRange) map[time.Time]int {
	hour := date.Begin
	dates := make(map[time.Time]int)
	column := 3
	for ; hour.Month() == date.End.Month(); column++ {
		dates[hour] = column
		if hour.Hour() == 23 {
			column++
		}
		hour = time.Date(hour.Year(), hour.Month(), hour.Day(), hour.Hour() + 1, hour.Minute(), hour.Second(), hour.Nanosecond(), hour.Location())
	}
	return dates
}

func instanceCountUsageReportInsertDataInSheet(_ context.Context, aas []aws.AwsAccount, file *excelize.File, data []instanceCount.InstanceCountReport, dates map[time.Time]int) (err error) {
	file.NewSheet(instanceCountReportDetailledSheetName)
	lastColumn := instanceCountUsageReportGenerateHeader(file, dates)
	line := 4
	reportsCells := make(cells, 0)
	for _, report := range data {
		reportCells := make(cells, 0, 3)
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		reportCells = append(reportCells, cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(report.InstanceCount.Region, "B"+strconv.Itoa(line)),
			newCell(report.InstanceCount.Type, "C"+strconv.Itoa(line)),
		}...)
		reportCells = instanceCountReportDatesAmountInSheet(file, dates, report, lastColumn, line, reportCells)
		reportsCells = append(reportsCells, reportCells...)
		line++
	}
	reportsCells.addStyles("borders", "centerText").setValues(file, instanceCountReportDetailledSheetName)
	return
}

//instanceCountReportDatesAmountInSheet put the Amount in terms of dates in the sheet
func instanceCountReportDatesAmountInSheet(file *excelize.File, dates map[time.Time]int, report instanceCount.InstanceCountReport,
lastColumn int, line int, reportCells cells) cells {
	totalColumnPosition := make([]string, 0)
	for date, column := range dates {
		reportCells = append(reportCells, newCell(0, excelize.ToAlphaString(column) + strconv.Itoa(line)))
		if date.Hour() == 0 {
			totalColumnPosition = append(totalColumnPosition, excelize.ToAlphaString(column + 24) + strconv.Itoa(line))
			formula := fmt.Sprintf("SUM(%s:%s)", excelize.ToAlphaString(column) + strconv.Itoa(line), excelize.ToAlphaString(column + 23) + strconv.Itoa(line))
			formulaLocation := excelize.ToAlphaString(column + 24) + strconv.Itoa(line)
			reportCells = append(reportCells, newFormula(formula, formulaLocation))
		}
	}
	for _, instanceDate := range report.InstanceCount.Hours {
		column, ok := dates[instanceDate.Hour]
		if ok {
			reportCells = append(reportCells, newCell(int(instanceDate.Count), excelize.ToAlphaString(column) + strconv.Itoa(line)))
		}
		file.SetColVisible(instanceCountReportDetailledSheetName, excelize.ToAlphaString(column), false)
	}
	formula := fmt.Sprintf("SUM(%s)", strings.Join(totalColumnPosition, ","))
	formulaLocation := excelize.ToAlphaString(lastColumn + 1) + strconv.Itoa(line)
	reportCells = append(reportCells, newFormula(formula, formulaLocation))
	return reportCells
}

func instanceCountUsageReportGenerateHeader(file *excelize.File, dates map[time.Time]int) int {
	lastColumn, totalColumnsWidth := instanceCountDatesHeader(file, dates)
	header := cells{
		newCell("Account", "A1").mergeTo("A3"),
		newCell("Region", "B1").mergeTo("B3"),
		newCell("Type", "C1").mergeTo("C3"),
		newCell("Dates", "D1").mergeTo(excelize.ToAlphaString(lastColumn) + "1"),
		newCell("Total", excelize.ToAlphaString(lastColumn + 1) + "1").mergeTo(excelize.ToAlphaString(lastColumn + 1) + "3"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, instanceCountReportDetailledSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 20),
		newColumnWidth("C", 30),
		newColumnWidth("D", 6).toColumn(excelize.ToAlphaString(lastColumn)),
		newColumnWidth(excelize.ToAlphaString(lastColumn + 1), 20),
	}
	columns = append(columns, totalColumnsWidth...)
	columns.setValues(file, instanceCountReportDetailledSheetName)
	return lastColumn
}

//instanceCountDatesHeader generate header for all the dates and total
func instanceCountDatesHeader(file *excelize.File, dates map[time.Time]int) (int, columnsWidth){
	widthColumn := 0
	lastColumn := 3
	totalColumnsWidth := make(columnsWidth, 0)
	for date, column := range dates {
		cellsDate := cells{
			newCell(date.Format("15:04"), excelize.ToAlphaString(column) + "3"),
		}
		cellsDate.addStyles("borders", "centerText").setValues(file, instanceCountReportDetailledSheetName)
		if column > widthColumn {
			widthColumn = column
			lastColumn = column + 1
		}
		if date.Hour() == 0 {
			cellDay := cells{
				newCell(date.Format("2006-01-02"), excelize.ToAlphaString(column) + "2").mergeTo(excelize.ToAlphaString(column + 24) + "2"),
				newCell("Total", excelize.ToAlphaString(column + 24) + "3"),
			}
			totalColumnsWidth = append(totalColumnsWidth, newColumnWidth(excelize.ToAlphaString(column + 24), 10))
			cellDay.addStyles("borders", "bold", "centerText").setValues(file, instanceCountReportDetailledSheetName)
		}
	}
	return lastColumn, totalColumnsWidth
}