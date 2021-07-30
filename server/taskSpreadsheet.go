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
	"strconv"
	"time"

	"database/sql"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/reports"
)

// taskSpreadsheet generates Spreadsheet with reports for a given AwsAccount.
func taskSpreadsheet(ctx context.Context) error {
	args := paramsFromContextOrArgs(ctx)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'Spreadsheet'.", map[string]interface{}{
		"args": args,
	})

	aaId, date, err := checkArguments(args)
	if err != nil {
		return err
	} else {
		return generateReport(ctx, aaId, date)
	}
}

func checkArguments(args []string) (int, time.Time, error) {
	var aaId int
	var date time.Time
	if len(args) != 1 && len(args) != 3 {
		return -1, time.Time{}, errors.New("taskSpreadsheet requires an integer argument (AWS Account ID) or three integer arguments (AWS Account ID, month and year)")
	} else if id, err := strconv.Atoi(args[0]); err != nil {
		return -1, time.Time{}, err
	} else {
		aaId = id
	}
	if len(args) == 3 {
		if month, err := strconv.Atoi(args[1]); err != nil {
			return -1, time.Time{}, err
		} else if year, err := strconv.Atoi(args[2]); err != nil {
			return -1, time.Time{}, err
		} else {
			now := time.Now().UTC()
			formattedMonth := time.Month(month)
			date = time.Date(year, formattedMonth, 1, 0, 0, 0, 0, now.Location()).UTC()
		}
	}
	return aaId, date, nil
}

func generateReport(ctx context.Context, aaId int, date time.Time) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	var updateId int64
	var generation bool
	forceGeneration := !date.IsZero()
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

	var user *models.User // We can't use := because then there would be a new err which would shadow the returned value
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if user, err = models.UserByID(db.Db, aa.UserId); err != nil || user.AccountType != "trackit" {
		if err == nil {
			logger.Info("Task 'SpreadSheet' has been skipped because the user has the wrong account type.", map[string]interface{}{
				"userAccountType": user.AccountType,
				"requiredAccount": "trackit",
			})
		}
	} else if generation, err = checkReportGeneration(ctx, db.Db, aa, forceGeneration); err != nil || !generation {
	} else if updateId, err = registerAccountReportGeneration(db.Db, aa); err != nil {
	} else {
		errs := reports.GenerateReport(ctx, aa, nil, date)
		updateAccountReportGenerationCompletion(ctx, aaId, db.Db, updateId, nil, errs, forceGeneration)
	}
	if err != nil {
		logger.Error("Error while generating spreadsheet report.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
		updateAccountReportGenerationCompletion(ctx, aaId, db.Db, updateId, err, nil, forceGeneration)
	}
	return
}

func checkReportGeneration(ctx context.Context, db *sql.DB, aa aws.AwsAccount, forceGeneration bool) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	if forceGeneration {
		return true, nil
	}
	startDate, endDate := history.GetHistoryDate()
	logger.Info("Checking report generation conditions", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	complete, err := history.CheckBillingDataCompleted(ctx, startDate, endDate, aa)
	if err != nil {
		logger.Info("Error while checking if billing data are completed", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err,
		})
		return false, err
	} else if !complete {
		logger.Info("Billing data are not completed", map[string]interface{}{
			"awsAccountId": aa.Id,
		})
		return false, nil
	}
	dbAccount, err := models.AwsAccountByID(db, aa.Id)
	if err != nil {
		logger.Info("Error while getting AWS account", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err,
		})
		return false, err
	}
	if dbAccount.LastSpreadsheetReportGeneration.Before(endDate) {
		return true, nil
	}
	dbProcessAccountJobs, err := models.AwsAccountUpdateJobsByAwsAccountID(db, aa.Id)
	if err != nil {
		logger.Info("Error while getting process account job", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err,
		})
		return false, err
	}
	if len(dbProcessAccountJobs) == 0 {
		return false, nil
	}
	generation := false
	for _, job := range dbProcessAccountJobs {
		if job.MonthlyReportsGenerated {
			if dbAccount.LastSpreadsheetReportGeneration.Before(job.Completed) {
				generation = true
			}
		}
	}
	if generation {
		return true, nil
	} else {
		logger.Info("No new monthly reports", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err,
		})
		return false, nil
	}
}

func registerAccountReportGeneration(db *sql.DB, aa aws.AwsAccount) (int64, error) {
	dbReportGeneration := models.AwsAccountReportsJob{
		AwsAccountID: aa.Id,
		WorkerID:     backendId,
	}
	err := dbReportGeneration.Insert(db)
	if err != nil {
		return 0, err
	}
	return int64(dbReportGeneration.ID), err
}

func updateAccountReportGenerationCompletion(ctx context.Context, aaId int, db *sql.DB, updateId int64, jobErr error, errs map[string]error, forceGeneration bool) {
	rErr := registerAccountReportGenerationCompletion(db, aaId, updateId, jobErr, errs, forceGeneration)
	if rErr != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to register account processing completion.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        rErr.Error(),
			"updateId":     updateId,
		})
	}
}

func registerAccountReportGenerationCompletion(db *sql.DB, aaId int, updateId int64, jobErr error, errs map[string]error, forceGeneration bool) error {
	dbAccountReports, err := models.AwsAccountReportsJobByID(db, int(updateId))
	if err != nil {
		return err
	}
	date := time.Now()
	dbAccountReports.Completed = date
	dbAccountReports.Joberror = errToStr(jobErr)
	dbAccountReports.Spreadsheeterror = errToStr(errs["speadsheetError"])
	dbAccountReports.Costdifferror = errToStr(errs["costDiffError"])
	dbAccountReports.Ec2usagereporterror = errToStr(errs["ec2UsageReportError"])
	dbAccountReports.Rdsusagereporterror = errToStr(errs["rdsUsageReportError"])
	dbAccountReports.Esusagereporterror = errToStr(errs["esUsageReportError"])
	dbAccountReports.Elasticacheusagereporterror = errToStr(errs["elasticacheUsageReportError"])
	dbAccountReports.Lambdausagereporterror = errToStr(errs["lambdaUsageReportError"])
	err = dbAccountReports.Update(db)
	if err != nil {
		return err
	}
	if !forceGeneration {
		var dbAccount *models.AwsAccount // We can't use := because then there would be a new err which would shadow the returned value
		dbAccount, err = models.AwsAccountByID(db, aaId)
		if err != nil {
			return err
		}
		dbAccount.LastSpreadsheetReportGeneration = date
		err = dbAccount.Update(db)
	}
	return err
}
