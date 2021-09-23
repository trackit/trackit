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
