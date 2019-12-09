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
	"database/sql"
	"errors"
	"flag"
	"strconv"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/db"
)

func checkArgumentsForElementalProcessAccount(args []string) (int, time.Time, error) {
	if len(args) < 1 {
		return invalidAccId, invalidTime, errors.New("elemental process account task requires at least an integer argument as AWS Account UD")
	}
	accId, err := strconv.Atoi(args[0])
	if err != nil {
		return invalidAccId, invalidTime, err
	}
	now := time.Now().UTC()
	if len(args) == 3 {
		if month, err := strconv.Atoi(args[1]); err != nil {
			return invalidAccId, invalidTime, err
		} else if year, err := strconv.Atoi(args[2]); err != nil {
			return invalidAccId, invalidTime, err
		} else {
			formattedMonth := time.Month(month)
			date := time.Date(year, formattedMonth, 1, 0, 0, 0, 0, now.Location()).UTC()
			return accId, date, nil
		}
	}
	return accId, invalidTime, nil
}

// taskElementalProcessAccount processes an AwsAccount to retrieve data from the AWS Elemental products api.
func taskElementalProcessAccount(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'elemental-process-account'.", map[string]interface{}{
		"args": args,
	})
	aaId, date, err := checkArgumentsForElementalProcessAccount(args)
	if err != nil {
		return err
	} else {
		return ingestElementalDataForAccount(ctx, aaId, date)
	}
}

// ingestElementalDataForAccount ingests the AWS api data for an AwsAccount.
func ingestElementalDataForAccount(ctx context.Context, aaId int, date time.Time) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	//var updateId int64
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
	} else if _, err = registerAccountElementalProcessing(db.Db, aa); err != nil { // _, err = updateId
	} else {
		if date.IsZero() {
			// TODO: elemental products ingest
		}
		_, _ = elementalProcessAccountHistory(ctx, aa, date) // _, _ = historyCreated, historyErr
		//updateAccountElementalProcessingCompletion(ctx, aaId, db.Db, updateId, nil, rdsErr, ec2Err, esErr, elastiCacheErr, lambdaErr, riEc2Err, riRdsErr, odToRiEc2Err, historyErr, ebsErr, historyCreated)
	}
	if err != nil {
		//updateAccountElementalProcessingCompletion(ctx, aaId, db.Db, updateId, err, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, false)
		logger.Error("Failed to process account elemental data.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	return
}

func registerAccountElementalProcessing(db *sql.DB, aa aws.AwsAccount) (int64, error) {
	const sqlstr = `INSERT INTO aws_account_elemental_update_job(
		aws_account_elemental_id,
		worker_id
	) VALUES (?, ?)`
	res, err := db.Exec(sqlstr, aa.Id, backendId)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateAccountElementalProcessingCompletion(ctx context.Context, aaId int, db *sql.DB, updateId int64, jobErr, historyErr error, historyCreated bool) {
	updateNextElementalUpdateAccount(db, aaId)
	rErr := registerAccountElementalProcessingCompletion(db, updateId, jobErr, historyErr, historyCreated)
	if rErr != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to register account processing elemental completion.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        rErr.Error(),
			"updateId":     updateId,
		})
	}
}

func updateNextElementalUpdateAccount(db *sql.DB, aaId int) error {
	const sqlstr = `UPDATE aws_account SET
		next_elemental_update=?
	WHERE id=?`
	_, err := db.Exec(sqlstr, time.Now().AddDate(0, 0, 1), aaId)
	return err
}

func registerAccountElementalProcessingCompletion(db *sql.DB, updateId int64, jobErr, historyErr error, historyCreated bool) error {
	const sqlstr = `UPDATE aws_account_elemental_update_job SET
		completed=?,
		jobError=?,
		historyError=?,
		monthly_reports_generated=?
	WHERE id=?`
	_, err := db.Exec(sqlstr, time.Now(), errToStr(jobErr), errToStr(historyErr), historyCreated, updateId)
	return err
}

// elementalProcessAccountHistory processes Elemental Products data with billing data for an AwsAccount
func elementalProcessAccountHistory(ctx context.Context, aa aws.AwsAccount, date time.Time) (bool, error) {
	status, err := history.FetchHistoryInfos(ctx, aa, date)
	if err != nil && err != history.ErrBillingDataIncomplete {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest History elemental data.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return status, err
}
