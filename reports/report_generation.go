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

func GenerateReport(ctx context.Context, aa aws.AwsAccount, aas []aws.AwsAccount, date time.Time) (errs map[string]error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	now := time.Now()
	var masterReport bool
	var reportDate string
	if date.IsZero() {
		reportDate = fmt.Sprintf("%s%s", (now.Month() - 1).String(), strconv.Itoa(now.Year()))
	} else {
		reportDate = fmt.Sprintf("%s%s", (date.Month()).String(), strconv.Itoa(date.Year()))
	}
	if aas == nil {
		aas = []aws.AwsAccount{aa}
		masterReport = false
		logger.Info("Generating spreadsheet for account", map[string]interface{}{
			"account":      aa,
			"date":          reportDate,
		})
	} else {
		masterReport = true
		logger.Info("Generating spreadsheet for accounts", map[string]interface{}{
			"masterAccount": aa,
			"accounts":      aas,
			"date":          reportDate,
		})
	}
	errs = make(map[string]error)
	file := createSpreadsheet(aa, reportDate)

	if tx, err := db.Db.BeginTx(ctx, nil); err == nil {
		for _, module := range modules {
			err := module.GenerateSheet(ctx, aas, date, tx, file.File)
			if err != nil {
				errs[module.ErrorName] = err
			}
		}
		errs["speadsheetError"] = saveSpreadsheet(ctx, file, masterReport)
	} else {
		errs["speadsheetError"] = err
	}
	return
}
