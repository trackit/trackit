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
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/reports"
)

// taskMasterSpreadsheet generates Spreadsheet with reports for a master AwsAccount including subaccounts.
func taskMasterSpreadsheet(ctx context.Context) error {
	args := paramsFromContextOrArgs(ctx)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'Master Spreadsheet'.", map[string]interface{}{
		"args": args,
	})

	aaId, date, err := checkArguments(args)
	if err != nil {
		return err
	} else {
		return generateMasterReport(ctx, aaId, date)
	}
}

func generateMasterReport(ctx context.Context, aaId int, date time.Time) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	var updateId int64
	var generation bool
	forceGeneration := !date.IsZero()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer utilsUsualTxFinalize(&tx, &err, &logger, "generate-master-spreadsheet")

	var user *models.User // We can't use := because then there would be a new err which would shadow the returned value
	var dbAccounts []*models.AwsAccount
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if user, err = models.UserByID(db.Db, aa.UserId); err != nil || user.AccountType != "trackit" {
		if err == nil {
			logger.Info("Task 'MasterSpreadSheet' has been skipped because the user has the wrong account type.", map[string]interface{}{
				"userAccountType": user.AccountType,
				"requiredAccount": "trackit",
			})
		}
	} else if generation, err = checkMasterReportGeneration(ctx, db.Db, aa, forceGeneration); err != nil || !generation {
	} else if updateId, err = registerMasterAccountReportGeneration(db.Db, aa); err != nil {
	} else if dbAccounts, err = getAccounts(ctx, db.Db, aa); err != nil {
	} else {
		accounts := make([]aws.AwsAccount, 0)
		for _, dbAccount := range dbAccounts {
			account := aws.AwsAccountFromDbAwsAccount(*dbAccount)
			accounts = append(accounts, account)
		}
		errs := reports.GenerateReport(ctx, aa, accounts, date)
		updateMasterAccountReportGenerationCompletion(ctx, aaId, db.Db, updateId, nil, errs, forceGeneration)
	}
	if err != nil {
		logger.Error("Error while generating spreadsheet report.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
		updateMasterAccountReportGenerationCompletion(ctx, aaId, db.Db, updateId, err, nil, forceGeneration)
	}
	return
}

func getAccounts(ctx context.Context, db *sql.DB, aa aws.AwsAccount) ([]*models.AwsAccount, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAllAccounts := make([]*models.AwsAccount, 0)
	dbMasterAccount, err := models.AwsAccountByID(db, aa.Id)
	if err != nil {
		logger.Info("Error while getting AWS account", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err,
		})
		return dbAllAccounts, err
	}
	dbAllAccounts = append(dbAllAccounts, dbMasterAccount)
	dbSubAccounts, err := models.AwsAccountsByParentId(db, aa.Id)
	if err != nil {
		logger.Info("Error while getting AWS sub account", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err,
		})
		return dbAllAccounts, err
	}
	dbAllAccounts = append(dbAllAccounts, dbSubAccounts...)
	return dbAllAccounts, nil
}

func checkMasterReportGeneration(ctx context.Context, db *sql.DB, aa aws.AwsAccount, forceGeneration bool) (bool, error) {
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
	dbAllAccounts, err := getAccounts(ctx, db, aa)
	if err != nil {
		return false, err
	}
	dbMasterAccount, err := models.AwsAccountByID(db, aa.Id)
	if err != nil {
		logger.Info("Error while getting AWS account", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err,
		})
		return false, err
	}
	if dbMasterAccount.LastMasterSpreadsheetReportGeneration.Before(endDate) {
		return true, nil
	}
	generation := false
	for _, account := range dbAllAccounts {
		if dbMasterAccount.LastMasterSpreadsheetReportGeneration.Before(account.LastSpreadsheetReportGeneration) {
			generation = true
		}
	}
	return generation, nil
}

func registerMasterAccountReportGeneration(db *sql.DB, aa aws.AwsAccount) (int64, error) {
	dbReportGeneration := models.AwsAccountMasterReportsJob{
		AwsAccountID: aa.Id,
		WorkerID:     backendId,
	}
	err := dbReportGeneration.Insert(db)
	if err != nil {
		return 0, err
	}
	return int64(dbReportGeneration.ID), err
}

func updateMasterAccountReportGenerationCompletion(ctx context.Context, aaId int, db *sql.DB, updateId int64, jobErr error, errs map[string]error, forceGeneration bool) {
	rErr := registerMasterAccountReportGenerationCompletion(db, aaId, updateId, jobErr, errs, forceGeneration)
	if rErr != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to register account processing completion.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        rErr.Error(),
			"updateId":     updateId,
		})
	}
}

func registerMasterAccountReportGenerationCompletion(db *sql.DB, aaId int, updateId int64, jobErr error, errs map[string]error, forceGeneration bool) error {
	dbAccountReports, err := models.AwsAccountMasterReportsJobByID(db, int(updateId))
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
		dbAccount.LastMasterSpreadsheetReportGeneration = date
		err = dbAccount.Update(db)
	}
	return err
}
