//   Copyright 2020 MSolution.IO
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

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/tagging"
)

func taskUpdateTaggingCompliance(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'update-tagging-compliance'.", map[string]interface{}{
		"args": args,
	})

	userId, err := checkUpdateTaggingComplianceArguments(args)
	if err != nil {
		logger.Error("Failed to execute task 'update-tagging-compliance'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	}

	err = updateTaggingComplianceForAccount(ctx, userId)

	if err == nil {
		logger.Info("Task 'update-tagging-compliance' done.", map[string]interface{}{
			"args": args,
		})
	}
	return err
}

func checkUpdateTaggingComplianceArguments(args []string) (int, error) {
	if len(args) < 1 {
		return invalidAccID, errors.New("Task 'update-tagging-compliance' requires at least an integer argument as Account ID")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return invalidAccID, err
	}

	return userId, nil
}

func updateTaggingComplianceForAccount(ctx context.Context, userId int) (err error) {
	var job models.AwsAccountUpdateTaggingComplianceJob
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	if job, err = registerUpdateTaggingComplianceTask(db.Db, userId); err != nil {
	} else {
		err = tagging.UpdateTaggingComplianceForAccount(ctx, userId)
		updateUpdateTaggingComplianceTask(db.Db, job, err)
	}

	if err != nil {
		logger.Error("Failed to execute task 'update-tagging-compliance'.", map[string]interface{}{
			"userId": userId,
			"error":  err.Error(),
		})
	}

	return
}

func registerUpdateTaggingComplianceTask(db *sql.DB, userId int) (models.AwsAccountUpdateTaggingComplianceJob, error) {
	job := models.AwsAccountUpdateTaggingComplianceJob{
		AwsAccountID: userId,
		WorkerID:     backendId,
		Created:      time.Now(),
	}

	err := job.Insert(db)

	return job, err
}

func updateUpdateTaggingComplianceTask(db *sql.DB, job models.AwsAccountUpdateTaggingComplianceJob, jobError error) error {
	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
