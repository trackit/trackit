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
)

var zeroDate = time.Date(0001, 1, 1, 00, 00, 00, 00, time.UTC)

func taskOnboardTagbot(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'onboard-tagbot'.", map[string]interface{}{
		"args": args,
	})

	userId, err := checkOnboardTagbotArguments(args)
	if err != nil {
		logger.Error("Failed to execute task 'onboard-tagbot'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	}
	if user, err := models.UserByID(db.Db, userId); err != nil {
		logger.Error("Failed to execute task 'onboard-tagbot'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	} else if user.AccountType != "tagbot" {
		return nil
	}

	job, err := registerOnboardTagbotTask(db.Db, userId)
	err = onboardTagbotUser(ctx, userId)
	updateOnboardTagbotTask(db.Db, job, err)

	if err == nil {
		logger.Info("Task 'onboard-tagbot' done.", map[string]interface{}{
			"args": args,
		})
	}
	return err
}

func checkOnboardTagbotArguments(args []string) (int, error) {
	if len(args) < 1 {
		return invalidUserID, errors.New("Task 'onboard-tagbot' requires at least an integer argument as User ID")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func onboardTagbotUser(ctx context.Context, userId int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	accounts, err := models.AwsAccountsByUserID(db.Db, userId)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		logger.Info("Task 'onboard-tagbot', ingesting data for an AWS account.", map[string]interface{}{
			"id": account.ID,
		})
		err = ingestDataForAccount(ctx, account.ID, zeroDate)
		if err == nil {
			logger.Info("Task 'onboard-tagbot', ingested data for an AWS account.", map[string]interface{}{
				"id": account.ID,
			})
		}
	}

	return updateTagsForUser(ctx, userId)
}

func registerOnboardTagbotTask(db *sql.DB, userId int) (models.UserOnboardTagbotJob, error) {
	job := models.UserOnboardTagbotJob{
		UserID:   userId,
		WorkerID: backendId,
		Created:  time.Now(),
	}

	err := job.Insert(db)

	return job, err
}

func updateOnboardTagbotTask(db *sql.DB, job models.UserOnboardTagbotJob, jobError error) error {
	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
