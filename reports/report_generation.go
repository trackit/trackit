//   Copyright 2018 MSolution.IO
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

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
)

type module struct {
	Name      string
	Function  func(context.Context, []aws.AwsAccount, time.Time, *sql.Tx) ([][]cell, error)
	ErrorName string
	Header    [][]cell
}

var modules = []module{
	{
		Name:      "EC2 Usage Report",
		Function:  getEc2UsageReport,
		ErrorName: "ec2UsageReportError",
		Header:    ec2InstanceFormat,
	},
	{
		Name:      "RDS Usage Report",
		Function:  getRdsUsageReport,
		ErrorName: "rdsUsageReportError",
		Header:    rdsInstanceFormat,
	},
	{
		Name:      "ElasticSearch Usage Report",
		Function:  getEsUsageReport,
		ErrorName: "esUsageReportError",
		Header:    esDomainFormat,
	},
	{
		Name:      "ElastiCache Usage Report",
		Function:  getElasticacheUsageReport,
		ErrorName: "elasticacheUsageReportError",
		Header:    elasticacheInstanceFormat,
	},
	{
		Name:      "Lambda Usage Report",
		Function:  getLambdaUsageReport,
		ErrorName: "lambdaUsageReportError",
		Header:    lambdaFunctionFormat,
	},
	{
		Name:      "Cost Differentiator Report",
		Function:  getCostDiff,
		ErrorName: "CostDifferentiatorError",
		Header:    costDiffHeader,
	},
}

func GenerateReport(ctx context.Context, aa aws.AwsAccount, date time.Time) (errs map[string]error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	now := time.Now()
	var reportDate string
	if date.IsZero() {
		reportDate = fmt.Sprintf("%s%s", (now.Month() - 1).String(), strconv.Itoa(now.Year()))
	} else {
		reportDate = fmt.Sprintf("%s%s", (date.Month()).String(), strconv.Itoa(date.Year()))
	}
	logger.Info("Generating spreadsheet for account", map[string]interface{}{
		"account": aa,
		"date": reportDate,
	})
	errs = make(map[string]error)

	if tx, err := db.Db.BeginTx(ctx, nil); err == nil {
		sheets, errs := getSpreadsheetData(ctx, aa, date, tx)
		if len(errs) > 0 {
			logger.Error("Error while getting spreadsheet data", errs)
		}

		file, spreadsheetErrors := generateSpreadsheet(ctx, aa, reportDate, sheets)
		if len(spreadsheetErrors) > 0 {
			logger.Error("Error while generating spreadsheet", spreadsheetErrors)
			for errorKey := range spreadsheetErrors {
				if _, ok := errs[errorKey]; !ok {
					errs[errorKey] = spreadsheetErrors[errorKey]
				}
			}
		}
		errs["speadsheetError"] = saveSpreadsheetLocally(ctx, file, false)
	} else {
		errs["speadsheetError"] = err
	}
	return
}

func GenerateMasterReport(ctx context.Context, aa aws.AwsAccount, aas []aws.AwsAccount, date time.Time) (errs map[string]error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	now := time.Now()
	var reportDate string
	if date.IsZero() {
		reportDate = fmt.Sprintf("%s%s", (now.Month() - 1).String(), strconv.Itoa(now.Year()))
	} else {
		reportDate = fmt.Sprintf("%s%s", (date.Month()).String(), strconv.Itoa(date.Year()))
	}
	logger.Info("Generating master spreadsheet for accounts", map[string]interface{}{
		"masterAccount": aa,
		"accounts": aas,
		"date": reportDate,
	})
	errs = make(map[string]error)

	if tx, err := db.Db.BeginTx(ctx, nil); err == nil {
		sheets, errs := getMasterSpreadsheetData(ctx, aas, date, tx)
		if len(errs) > 0 {
			logger.Error("Error while getting spreadsheet data", errs)
		}

		file, spreadsheetErrors := generateSpreadsheet(ctx, aa, reportDate, sheets)
		if len(spreadsheetErrors) > 0 {
			logger.Error("Error while generating spreadsheet", spreadsheetErrors)
			for errorKey := range spreadsheetErrors {
				if _, ok := errs[errorKey]; !ok {
					errs[errorKey] = spreadsheetErrors[errorKey]
				}
			}
		}
		errs["speadsheetError"] = saveSpreadsheetLocally(ctx, file, true)
	} else {
		errs["speadsheetError"] = err
	}
	return
/*
	for _, account := range aas {
		if file, err := loadSpreadsheetLocally(ctx, account, reportDate, false); err != nil {
			errs[account.Id] = err
		} else {
			files = append(files, file)
		}
	}

	masterFile, err := generateMasterSpreadsheet(ctx, aa, reportDate, files)
	if err != nil {
		logger.Error("Error while generating master spreadsheet", err)
		return err
	}
	return saveSpreadsheetLocally(ctx, masterFile, true)*/
}

func getSpreadsheetData(ctx context.Context, aa aws.AwsAccount, date time.Time, tx *sql.Tx) ([]sheet, map[string]error) {
	sheets := make([]sheet, 0)
	errors := make(map[string]error)

	for _, module := range modules {
		data, err := module.Function(ctx, []aws.AwsAccount{aa}, date, tx)
		if err != nil {
			errors[module.ErrorName] = err
		} else {
			sheets = append(sheets, sheet{name: module.Name, data: data})
		}
	}

	return sheets, errors
}

func getMasterSpreadsheetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) ([]sheet, map[string]error) {
	sheets := make([]sheet, 0)
	errors := make(map[string]error)

	for _, module := range modules {
		data, err := module.Function(ctx, aas, date, tx)
		if err != nil {
			errors[module.ErrorName] = err
		} else {
			sheets = append(sheets, sheet{name: module.Name, data: data})
		}
	}

	return sheets, errors
}
