//   Copyright 2017 MSolution.IO
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

package main

import (
	"context"
	"errors"
	"flag"
	"strconv"
	"time"

	"database/sql"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/reports"
)

// taskSpreadsheet generates Spreadsheet with reports for a given AwsAccount.
func taskSpreadsheet(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'Spreadsheet'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("taskSpreadsheet requires an integer argument")
	} else if aaId, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else {
		return generateReport(ctx, aaId)
	}
}

func generateReport(ctx context.Context, aaId int) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	var updateId int64
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if updateId, err = registerAccountReportGeneration(db.Db, aa); err != nil {
	} else {
		errs := reports.GenerateReport(ctx, aa)
		updateAccountReportGenerationCompletion(ctx, aaId, db.Db, updateId, nil, errs)
	}
	if err != nil {
		logger.Error("Error while generating spreadsheet report.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
		updateAccountReportGenerationCompletion(ctx, aaId, db.Db, updateId, err, nil)
	}
	return
}

func registerAccountReportGeneration(db *sql.DB, aa aws.AwsAccount) (int64, error) {
	const sqlstr = `INSERT INTO aws_account_reports_job(
		aws_account_id,
		worker_id
	) VALUES (?, ?)`
	res, err := db.Exec(sqlstr, aa.Id, backendId)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateAccountReportGenerationCompletion(ctx context.Context, aaId int, db *sql.DB, updateId int64, jobErr error, errs map[string]error) {
	rErr := registerAccountReportGenerationCompletion(db, updateId, jobErr, errs)
	if rErr != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to register account processing completion.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        rErr.Error(),
			"updateId":     updateId,
		})
	}
}

func registerAccountReportGenerationCompletion(db *sql.DB, updateId int64, jobErr error, errs map[string]error) error {
	const sqlstr = `UPDATE aws_account_reports_job SET
		completed=?,
		jobError=?,
		spreadsheetError=?,
		costDiffError=?,
		ec2UsageReportError=?,
		rdsUsageReportError=?
	WHERE id=?`
	_, err := db.Exec(sqlstr, time.Now(), errToStr(jobErr),
		errToStr(errs["speadsheetError"]), errToStr(errs["costDiffError"]),
		errToStr(errs["ec2UsageReportError"]), errToStr(errs["rdsUsageReportError"]),
		updateId)
	return err
}
