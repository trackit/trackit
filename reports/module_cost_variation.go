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
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/costs/diff"
)

type (
	costVariationProduct   map[time.Time]diff.PricePoint
	costVariationReport    map[string]costVariationProduct
	costVariationFrequency struct {
		SheetName   string
		Title       string
		Aggregation string
		DateFormat  string
		Dates       []time.Time
	}
)

const (
	costVariationLastMonthSheetName   = "Cost Variations (Last Month)"
	costVariationLast6MonthsSheetName = "Cost Variations (Last 6 Months)"
)

var costVariationLastMonth = module{
	Name:          "Cost Variations (Last Month)",
	SheetName:     costVariationLastMonthSheetName,
	ErrorName:     "costVariationLastMonthError",
	GenerateSheet: costVariationGenerateLastMonth,
}

var costVariationLast6Months = module{
	Name:          "Cost Variations (Last 6 Months)",
	SheetName:     costVariationLast6MonthsSheetName,
	ErrorName:     "costVariationLastSixMonthsError",
	GenerateSheet: costVariationGenerateLast6Months,
}

// costVariationGenerateLastMonth will generate a sheet with daily cost variation data for last month
// It will get cost data for given AWS account and for a given date
func costVariationGenerateLastMonth(ctx context.Context, aas []aws.AwsAccount, date time.Time, _ *sql.Tx, file *excelize.File) (err error) {
	var dateRange diff.DateRange
	if date.IsZero() {
		dateRange.Begin, dateRange.End = history.GetHistoryDate()
	} else {
		dateRange = diff.DateRange{
			Begin: date,
			End:   time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location()).UTC(),
		}
	}
	dates := make([]time.Time, dateRange.End.Day())
	for index := range dates {
		dates[index] = dateRange.Begin.AddDate(0, 0, index)
	}
	frequency := costVariationFrequency{costVariationLastMonthSheetName, "Daily Cost", "day", "2006-01-02", dates}
	return costVariationGenerateSheet(ctx, aas, dateRange, frequency, file)
}

// costVariationGenerateLast6Months will generate a sheet with monthly cost variation data for last 6 months
// It will get cost data for given AWS account and for a given date
func costVariationGenerateLast6Months(ctx context.Context, aas []aws.AwsAccount, date time.Time, _ *sql.Tx, file *excelize.File) (err error) {
	var dateRange diff.DateRange
	if date.IsZero() {
		_, dateRange.End = history.GetHistoryDate()
	} else {
		dateRange.End = time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location()).UTC()
	}
	dateRange.Begin = time.Date(dateRange.End.Year(), dateRange.End.Month()-5, 1, 0, 0, 0, 0, dateRange.End.Location()).UTC()
	dates := make([]time.Time, 6)
	for index := range dates {
		dates[index] = dateRange.Begin.AddDate(0, index, 0)
	}
	frequency := costVariationFrequency{costVariationLast6MonthsSheetName, "Monthly Cost", "month", "2006-01", dates}
	return costVariationGenerateSheet(ctx, aas, dateRange, frequency, file)
}

func costVariationGenerateSheet(ctx context.Context, aas []aws.AwsAccount, dateRange diff.DateRange, frequency costVariationFrequency, file *excelize.File) (err error) {
	data, err := costVariationGetData(ctx, aas, dateRange, frequency)
	if err == nil {
		return costVariationInsertDataInSheet(frequency, file, data)
	}
	return
}

func costVariationGetData(ctx context.Context,
	aas []aws.AwsAccount,
	dateRange diff.DateRange,
	frequency costVariationFrequency) (
	data map[aws.AwsAccount]costVariationReport,
	err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Getting Cost Variation Report for accounts", map[string]interface{}{
		"accounts":    aas,
		"dateStart":   dateRange.Begin,
		"dateEnd":     dateRange.End,
		"aggregation": frequency.Aggregation,
	})
	data = make(map[aws.AwsAccount]costVariationReport, len(aas))
	for _, account := range aas {
		report, err := diff.TaskDiffData(ctx, account, dateRange, frequency.Aggregation)
		if err != nil {
			logger.Error("An error occurred while generating a Cost Variation Report", map[string]interface{}{
				"error":     err,
				"account":   account,
				"dateStart": dateRange.Begin,
				"dateEnd":   dateRange.End,
			})
			return data, err
		}
		data[account] = make(costVariationReport, len(report))
		for product, values := range report {
			data[account][product], err = costVariationFormatCostDiff(values)
			if err != nil {
				logger.Error("An error occurred while parsing timestamp", map[string]interface{}{
					"error":     err,
					"account":   account,
					"values":    values,
					"dateStart": dateRange.Begin,
					"dateEnd":   dateRange.End,
				})
				return data, err
			}
		}
	}
	return
}

func costVariationInsertDataInSheet(
	frequency costVariationFrequency,
	file *excelize.File,
	data map[aws.AwsAccount]costVariationReport) (err error) {
	file.NewSheet(frequency.SheetName)
	costVariationGenerateHeader(file, frequency.SheetName, frequency.Dates, frequency)
	line := 4
	for account, report := range data {
		for product, values := range report {
			cells := make(cells, 0, len(frequency.Dates)*2+2)
			cells = append(cells, newCell(formatAwsAccount(account), "A"+strconv.Itoa(line)),
				newCell(product, "B"+strconv.Itoa(line)))
			totalNeededCols := make([]string, 0, len(frequency.Dates))
			for index, date := range frequency.Dates {
				value := costVariationGetValueForDate(values, date)
				if index == 0 {
					cells = append(cells, newCell(value.Cost, "C"+strconv.Itoa(line)).addStyles("price"))
				} else {
					costPos := excelize.ToAlphaString(index*2+2) + strconv.Itoa(line)
					previousCol := excelize.ToAlphaString(index*2) + strconv.Itoa(line)
					posVariation := excelize.ToAlphaString(index*2+1) + strconv.Itoa(line)
					formula := fmt.Sprintf(`IF(%s=0,"",%s/%s-1)`, previousCol, costPos, previousCol)
					variation := newFormula(formula, posVariation).addStyles("percentage")
					variation = variation.addConditionalFormat("negative", "green", "borders")
					variation = variation.addConditionalFormat("positive", "red", "borders")
					variation = variation.addConditionalFormat("zero", "red", "borders")
					cells = append(cells, newCell(value.Cost, costPos).addStyles("price"), variation)
					totalNeededCols = append(totalNeededCols, costPos)
				}
			}
			totalCol := len(frequency.Dates)*2 + 1
			formula := fmt.Sprintf("SUM(%s)", strings.Join(totalNeededCols, ","))
			cells = append(cells, newFormula(formula, excelize.ToAlphaString(totalCol)+strconv.Itoa(line)).addStyles("price"))
			cells.addStyles("borders", "centerText").setValues(file, frequency.SheetName)
			line++
		}
	}
	return
}

func costVariationGenerateHeader(file *excelize.File, sheetName string, dates []time.Time, frequency costVariationFrequency) {
	header := make(cells, 0, len(dates)*3+2)
	totalCol := excelize.ToAlphaString(len(dates)*2 + 1)
	header = append(header, newCell("Account", "A1").mergeTo("A3"),
		newCell("Usage type", "B1").mergeTo("B3"),
		newCell(frequency.Title, "C1").mergeTo(excelize.ToAlphaString(len(dates)*2)+"1"),
		newCell("Total", totalCol+"1").mergeTo(totalCol+"3"))
	for index, date := range dates {
		if index == 0 {
			header = append(header, newCell(date.Format(frequency.DateFormat), "C2"),
				newCell("Cost", "C3"))
		} else {
			col1 := excelize.ToAlphaString(index*2 + 1)
			col2 := excelize.ToAlphaString(index*2 + 2)
			header = append(header, newCell(date.Format(frequency.DateFormat), col1+"2").mergeTo(col2+"2"),
				newCell("Variation", col1+"3"),
				newCell("Cost", col2+"3"))
		}
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, sheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35),
		newColumnWidth("C", 12.5).toColumn(excelize.ToAlphaString(len(dates) * 2)),
		newColumnWidth(totalCol, 15),
	}
	columns.setValues(file, sheetName)
}

func costVariationFormatCostDiff(data []diff.PricePoint) (values costVariationProduct, err error) {
	values = make(costVariationProduct)
	for _, value := range data {
		date, err := time.Parse("2006-01-02T15:04:05.000Z", value.Date)
		if err != nil {
			return values, err
		}
		values[date] = value
	}
	return
}

func costVariationGetValueForDate(values costVariationProduct, date time.Time) *diff.PricePoint {
	if value, ok := values[date]; ok {
		return &value
	}
	return nil
}
