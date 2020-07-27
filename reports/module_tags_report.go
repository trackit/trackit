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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/costs/tags"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/lineItems"
	"github.com/trackit/trackit/users"
)

const tagsUsageReportSheetName = "Tags Report"

var tagsUsageReportModule = module{
	Name:          "Tags Usage Report",
	SheetName:     tagsUsageReportSheetName,
	ErrorName:     "tagsUsageReportError",
	GenerateSheet: generateTagsUsageReportSheet,
}

// generateTagsUsageReportSheet will generate a sheet with Tags usage report
// It will get data for given AWS account and for a given date
func generateTagsUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return tagsUsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func tagsUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := tagsUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		for key, report := range data {
			if err = tagsUsageReportInsertDataInSheet(ctx, file, key, report); err != nil {
				return
			}
		}
	}
	return
}

func getTagsKey(ctx context.Context, parsedParams tags.TagsValuesQueryParams) ([]string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	parsedParamsKeys := tags.TagsKeysQueryParams{
		AccountList: []string{},
		IndexList:   []string{},
		DateBegin:   parsedParams.DateBegin,
		DateEnd:     parsedParams.DateEnd,
	}
	_, keys, err := tags.GetTagsKeysWithParsedParams(ctx, parsedParamsKeys)
	if err != nil {
		logger.Error("Error when getting Tags Keys", err)
		return nil, err
	}
	return keys, nil
}

func tagsUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports tags.TagsValuesResponse, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parsedParams := tags.TagsValuesQueryParams{
		AccountList: []string{},
		IndexList:   []string{},
		DateBegin:   date,
		DateEnd:     time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location()).UTC(),
		TagsKeys:    []string{},
		By:          "product",
		Detailed:    true,
	}
	parsedParams.AccountList = identities
	accountsAndIndexes, _, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, lineItems.IndexSuffix)
	if err != nil {
		return nil, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	keys, err := getTagsKey(ctx, parsedParams)
	if err != nil {
		return nil, err
	}
	parsedParams.TagsKeys = keys
	logger.Debug("Getting Tags Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = tags.GetTagsValuesWithParsedParams(ctx, parsedParams)
	return
}

func tagsUsageReportInsertDataInSheet(ctx context.Context, file *excelize.File, key string, data []tags.TagsValues) (err error) {
	file.NewSheet(key)
	tagsUsageReportGenerateHeader(file, key)
	column := 2
	maxColumn := 2
	var totalCostCells []string
	var valueProductExist bool
	charColumn := 3
	for _, tag := range data {
		if tag.Tag == "" {
			tag.Tag = "no tag"
		}
		maxColumn, totalCostCells, valueProductExist = putProductDataInSheet(file, key, tag, maxColumn)
		if valueProductExist {
			cellsTagsCosts := cells{
				newCell(tag.Tag, "A"+strconv.Itoa(column)).mergeTo("A" + strconv.Itoa(maxColumn-1)),
				newFormula(fmt.Sprintf("SUM(%s)", strings.Join(totalCostCells, ",")), "F"+strconv.Itoa(column)).mergeTo("F" + strconv.Itoa(maxColumn-1)).addStyles("price"),
			}
			cellsTagsCosts.addStyles("borders", "centerText").setValues(file, key)
			cellsChart := cells{
				newCell(file.GetCellValue(key, fmt.Sprintf("A%s", strconv.Itoa(column))), "H"+strconv.Itoa(charColumn)),
				newFormula(fmt.Sprintf("=%s", "F"+strconv.Itoa(column)), "I"+strconv.Itoa(charColumn)).addStyles("price"),
			}
			cellsChart.addStyles("borders", "centerText").setValues(file, key)
			column = maxColumn
			charColumn++
		}
	}
	chartData := fmt.Sprintf(`{"name":"Costs","categories":"%s!$H$3:$H$%s","values":"%s!$I$3:$I$%s"}`, key, strconv.Itoa(charColumn-1), key, strconv.Itoa(charColumn-1))
	if err := generateLinearChart(ctx, file, key, chartData); err != nil {
		return err
	}
	return
}

func putProductDataInSheet(file *excelize.File, sheetName string, tag tags.TagsValues, maxColumn int) (int, []string, bool) {
	productColumn := maxColumn
	totalCostCells := make([]string, 0)
	valueProductExist := false
	for _, product := range tag.Items {
		productCostCells := make([]string, 0)
		valueTypeExist := false
		usageTypes := sortUsageTypesCosts(product.UsageTypes)
		for _, usageType := range usageTypes {
			if usageType.Cost != 0 {
				usageTypeCells := cells{
					newCell(usageType.UsageType, "C"+strconv.Itoa(maxColumn)),
					newCell(usageType.Cost, "D"+strconv.Itoa(maxColumn)).addStyles("price"),
				}
				usageTypeCells.addStyles("borders", "centerText").setValues(file, sheetName)
				maxColumn++
				productCostCells = append(productCostCells, "D"+strconv.Itoa(maxColumn-1))
				valueTypeExist = true
			}
		}
		if valueTypeExist {
			productCells := cells{
				newCell(product.Item, "B"+strconv.Itoa(productColumn)).mergeTo("B" + strconv.Itoa(maxColumn-1)),
				newFormula(fmt.Sprintf("SUM(%s)", strings.Join(productCostCells, ",")), "E"+strconv.Itoa(productColumn)).mergeTo("E" + strconv.Itoa(maxColumn-1)).addStyles("price"),
			}
			productCells.addStyles("borders", "centerText").setValues(file, sheetName)
			totalCostCells = append(totalCostCells, "E"+strconv.Itoa(productColumn))
			productColumn = maxColumn
			valueProductExist = true
		}
	}
	return maxColumn, totalCostCells, valueProductExist
}

func tagsUsageReportGenerateHeader(file *excelize.File, key string) {
	header := cells{
		newCell("Tags", "A1"),
		newCell("Products", "B1"),
		newCell("UsageTypes", "C1"),
		newCell("Costs", "D1"),
		newCell("Product Cost", "E1"),
		newCell("Total Cost", "F1"),
		newCell("Resume", "H1").mergeTo("I1"),
		newCell("Tags", "H2"),
		newCell("Total Cost", "I2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, key)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 20),
		newColumnWidth("C", 45),
		newColumnWidth("D", 15).toColumn("F"),
		newColumnWidth("H", 30),
	}
	columns.setValues(file, key)
	return
}

func generateLinearChart(ctx context.Context, file *excelize.File, sheetName, data string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	err := file.AddChart(sheetName, "J1", fmt.Sprintf(
		`{"type":"bar",
				"series":[%s],
				"format":{
					"x_scale":2.0,"y_scale":2.0,
					"x_offset":15,
					"y_offset":10,
					"print_obj":true,
					"lock_aspect_ratio":false,
					"locked":false
				},
				"legend":{
					"position":"bottom",
					"show_legend_key":false
				},
				"title":{"name":"Cost per tag"},
				"plotarea":{
					"show_bubble_size":true,
					"show_cat_name":false,
					"show_leader_lines":false,
					"show_percent":true,
					"show_series_name":false,
					"show_val":false
				},
				"show_blanks_as":"gap"}`, data))
	if err != nil {
		logger.Error("Error when generating the chart", err)
	}
	return nil
}

func sortUsageTypesCosts(usageTypes []tags.ValueDetailed) (sortTypes []tags.ValueDetailed) {
	sortTypes = make([]tags.ValueDetailed, 0)
	toSort := make([]float64, 0)
	for _, usageType := range usageTypes {
		toSort = append(toSort, usageType.Cost)
	}
	sort.Float64s(toSort)
	sort.Sort(sort.Reverse(sort.Float64Slice(toSort)))
	for _, numberSort := range toSort {
		for _, usageType := range usageTypes {
			if alreadySort := checkTypeIsSort(usageType.UsageType, sortTypes); !alreadySort && usageType.Cost == numberSort {
				sortTypes = append(sortTypes, usageType)
				break
			}
		}
	}
	return
}

func checkTypeIsSort(usageType string, sortTypes []tags.ValueDetailed) bool {
	for _, sortType := range sortTypes {
		if sortType.UsageType == usageType {
			return true
		}
	}
	return false
}
