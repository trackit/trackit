//   Copyright 2021 MSolution.IO
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

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	unusedaccounts "github.com/trackit/trackit/unusedAccounts"
)

// taskCheckUnusedAccounts checks if there are mails to send or data to delete because accounts are unused
func taskCheckUnusedAccounts(ctx context.Context) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Running task 'check-unused-accounts'.", nil)

	err = checkUnusedAccounts(ctx)
	if err != nil {
		logger.Error("Failed to execute task 'check-unused-accounts'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	}

	logger.Info("Task 'check-unused-accounts' done.", nil)
	return nil
}

func checkUnusedAccounts(ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	tx, err := db.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer utilsUsualTxFinalize(&tx, &err, &logger, "check-unused-accounts")

	job, err := registerCheckUnusedAccountsTask(db.Db)
	if err != nil {
		return err
	}

	err = unusedaccounts.CheckUnusedAccounts(ctx)
	if err != nil {
		logger.Error("Failed to execute task 'check-unused-accounts'.", map[string]interface{}{
			"err": err.Error(),
		})
	}

	return updateCheckUnusedAccountsTask(db.Db, job, err)
}

func registerCheckUnusedAccountsTask(db *sql.DB) (models.CheckUnusedAccountsJob, error) {
	job := models.CheckUnusedAccountsJob{
		WorkerID: backendId,
	}

	err := job.Insert(db)

	return job, err
}

func updateCheckUnusedAccountsTask(db *sql.DB, job models.CheckUnusedAccountsJob, jobError error) error {

	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
