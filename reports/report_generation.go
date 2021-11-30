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
	"fmt"
	"strconv"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/db"
)

// GenerateReport will generate a spreadsheet report for a given AWS account and for a given month
// It will iterate over available modules and generate a sheet for each module.
// Note: First sheet is removed since it is unused (Created by excelize)
// Report is then uploaded to an S3 bucket
// Note: File can be saved locally by using `saveSpreadsheetLocally` instead of `saveSpreadsheet`
func GenerateReport(ctx context.Context, aa aws.AwsAccount, aas []aws.AwsAccount, date time.Time) (errs map[string]error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	now := time.Now()
	var reportType spreadsheetType
	var reportDate string
	if date.IsZero() {
		if now.Month() != time.January {
			reportDate = fmt.Sprintf("%s%s", (now.Month() - 1).String(), strconv.Itoa(now.Year()))
		} else {
			reportDate = fmt.Sprintf("%s%s", time.December.String(), strconv.Itoa(now.Year()-1))
		}
	} else {
		reportDate = fmt.Sprintf("%s%s", (date.Month()).String(), strconv.Itoa(date.Year()))
	}
	if aas == nil {
		aas = []aws.AwsAccount{aa}
		reportType = regularReport
		logger.Info("Generating spreadsheet for account", map[string]interface{}{
			"account": aa,
			"date":    reportDate,
		})
	} else {
		reportType = masterReport
		logger.Info("Generating spreadsheet for accounts", map[string]interface{}{
			"masterAccount": aa,
			"accounts":      aas,
			"date":          reportDate,
		})
	}
	errs = make(map[string]error)
	file := createSpreadsheet(aa, reportDate)
	if tx, err := db.Db.BeginTx(ctx, nil); err != nil {
		errs["speadsheetError"] = err
	} else {
		for _, module := range modules {
			err = module.GenerateSheet(ctx, aas, date, tx, file.File)
			if err != nil {
				errs[module.ErrorName] = err
			}
		}
		file.File.DeleteSheet(file.File.GetSheetName(1))
		errs["speadsheetError"] = saveSpreadsheet(ctx, file, reportType)
	}
	return
}

// GenerateTagsReport will generate a spreadsheet tags report for a given AWS account and for a given month
// It will iterate over available modules and generate a sheet for each module.
// Note: First sheet is removed since it is unused (Created by excelize)
// Report is then uploaded to an S3 bucket
// Note: File can be saved locally by using `saveSpreadsheetLocally` instead of `saveSpreadsheet`
func GenerateTagsReport(ctx context.Context, aa aws.AwsAccount, aas []aws.AwsAccount, date time.Time) (errs map[string]error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	now := time.Now()
	var reportDate string
	if date.IsZero() {
		if now.Month() != time.January {
			reportDate = fmt.Sprintf("%s%s", (now.Month() - 1).String(), strconv.Itoa(now.Year()))
		} else {
			reportDate = fmt.Sprintf("%s%s", time.December.String(), strconv.Itoa(now.Year()-1))
		}
	} else {
		reportDate = fmt.Sprintf("%s%s", (date.Month()).String(), strconv.Itoa(date.Year()))
	}
	if aas == nil {
		aas = []aws.AwsAccount{aa}
		logger.Info("Generating spreadsheet tags for account", map[string]interface{}{
			"account": aa,
			"date":    reportDate,
		})
	} else {
		logger.Info("Generating spreadsheet tags for accounts", map[string]interface{}{
			"masterAccount": aa,
			"accounts":      aas,
			"date":          reportDate,
		})
	}
	errs = make(map[string]error)
	file := createSpreadsheet(aa, reportDate)
	if tx, err := db.Db.BeginTx(ctx, nil); err != nil {
		errs["speadsheetError"] = err
	} else {
		err = generateTagsUsageReportSheet(ctx, aas, date, tx, file.File)
		if err != nil {
			errs["tagsError"] = err
		}
		file.File.DeleteSheet(file.File.GetSheetName(1))
		errs["speadsheetError"] = saveSpreadsheet(ctx, file, tagsReport)
	}
	return
}
