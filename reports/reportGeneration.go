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
	Function  func(context.Context, aws.AwsAccount, *sql.Tx) ([][]string, error)
	ErrorName string
}

var modules = []module{
	{
		Name: "EC2 Usage Report",
		Function: getEc2UsageReport,
		ErrorName: "ec2UsageReportError",
	},
	{
		Name: "RDS Usage Report",
		Function: getRdsUsageReport,
		ErrorName: "rdsUsageReportError",
	},
	{
		Name: "Cost Differentiator Report",
		Function: getCostDiff,
		ErrorName: "CostDifferentiatorError",
	},
}

func GenerateReport(ctx context.Context, aa aws.AwsAccount) (errs map[string]error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Generating spreadsheet for account", aa)
	now := time.Now()
	date := fmt.Sprintf("%s%s", (now.Month()-1).String(), strconv.Itoa(now.Year()))
	errs = make(map[string]error)

	if tx, err := db.Db.BeginTx(ctx, nil); err == nil {
		sheets, errs := getSpreadsheetData(ctx, aa, tx)
		if len(errs) > 0 {
			logger.Error("Error while getting spreadsheet data", errs)
		}

		file, spreadsheetErrors := generateSpreadsheet(ctx, aa, date, sheets)
		if len(spreadsheetErrors) > 0 {
			logger.Error("Error while generating spreadsheet", spreadsheetErrors)
			for errorKey := range spreadsheetErrors {
				if _, ok := errs[errorKey]; !ok {
					errs[errorKey] = spreadsheetErrors[errorKey]
				}
			}
		}
		errs["speadsheetError"] = saveSpreadsheet(ctx, file)
	} else {
		errs["speadsheetError"] = err
	}
	return
}

func getSpreadsheetData(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) ([]Sheet, map[string]error) {
	sheets := make([]Sheet, 0)
	errors := make(map[string]error)

	for _, module := range modules {
		data, err := module.Function(ctx, aa, tx)
		if err != nil {
			errors[module.ErrorName] = err
		} else {
			sheets = append(sheets, Sheet{Name: module.Name, Data: data})
		}
	}

	return sheets, errors
}
