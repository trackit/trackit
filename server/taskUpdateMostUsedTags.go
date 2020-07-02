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

func taskUpdateMostUsedTags(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'update-most-used-tags'.", map[string]interface{}{
		"args": args,
	})

	amazonAccountID, err := checkUpdateMostUsedTagsArguments(args)
	if err != nil {
		return err
	}

	err = updateMostUsedTagsForAccount(ctx, amazonAccountID)
	logger.Info("Task 'update-most-used-tags' done.", map[string]interface{}{
		"args": args,
	})

	return err
}

func checkUpdateMostUsedTagsArguments(args []string) (int, error) {
	if len(args) < 1 {
		return invalidAccID, errors.New("Task 'update-most-used-tags' requires at least an integer argument as AWS Account ID")
	}

	amazonAccountID, err := strconv.Atoi(args[0])
	if err != nil {
		return invalidAccID, err
	}

	return amazonAccountID, nil
}

func updateMostUsedTagsForAccount(ctx context.Context, accountID int) error {
	tx, err := db.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()

	awsAccount, err := aws.GetAwsAccountWithId(accountID, tx)
	if err != nil {
		return nil
	}

	job, err := registerUpdateMostUsedTagsTask(db.Db, accountID)
	if err != nil {
		return err
	}

	err = tagging.UpdateMostUsedTagsForAccount(ctx, accountID, awsAccount.AwsIdentity)

	return updateUpdateMostUsedTagsTask(db.Db, job, err)
}

func registerUpdateMostUsedTagsTask(db *sql.DB, accountID int) (models.AwsAccountUpdateMostUsedTagsJob, error) {
	job := models.AwsAccountUpdateMostUsedTagsJob{
		AwsAccountID: accountID,
		WorkerID:     backendId,
		Created:      time.Now(),
	}

	err := job.Insert(db)

	return job, err
}

func updateUpdateMostUsedTagsTask(db *sql.DB, job models.AwsAccountUpdateMostUsedTagsJob, jobError error) error {

	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
