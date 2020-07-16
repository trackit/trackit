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

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/tagging"
)

const invalidAccID = -1

func taskUpdateTags(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'update-tags'.", map[string]interface{}{
		"args": args,
	})

	aaId, err := checkUpdateTagsArguments(args)
	if err != nil {
		logger.Error("Failed to execute task 'update-tags'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	}

	err = updateTagsForAccount(ctx, aaId)

	if err == nil {
		logger.Info("Task 'update-tags' done.", map[string]interface{}{
			"args": args,
		})
	}
	return nil
}

func checkUpdateTagsArguments(args []string) (int, error) {
	if len(args) < 1 {
		return invalidAccID, errors.New("Task 'update-tags' requires at least an integer argument as AWS Account ID")
	}

	aaId, err := strconv.Atoi(args[0])
	if err != nil {
		return invalidAccID, err
	}

	return aaId, nil
}

func updateTagsForAccount(ctx context.Context, aaId int) (err error) {
	var tx *sql.Tx
	var awsAccount aws.AwsAccount
	var job models.AwsAccountUpdateTagsJob
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

	if tx, err := db.Db.BeginTx(ctx, nil); err != nil {
	} else if awsAccount, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if job, err = registerUpdateTagsTask(db.Db, aaId); err != nil {
	} else {
		err = tagging.UpdateTagsForAccount(ctx, awsAccount)
		updateUpdateTagsTask(db.Db, job, err)
	}
	if err != nil {
		logger.Error("Failed to process account data.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	return
}

func registerUpdateTagsTask(db *sql.DB, aaId int) (models.AwsAccountUpdateTagsJob, error) {
	job := models.AwsAccountUpdateTagsJob{
		AwsAccountID: aaId,
		WorkerID:     backendId,
		Created:      time.Now(),
	}

	err := job.Insert(db)

	return job, err
}

func updateUpdateTagsTask(db *sql.DB, job models.AwsAccountUpdateTagsJob, jobError error) error {

	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
