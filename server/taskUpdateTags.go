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

const invalidUserID = -1

func taskUpdateTags(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'update-tags'.", map[string]interface{}{
		"args": args,
	})

	userId, err := checkUpdateTagsArguments(args)
	if err != nil {
		logger.Error("Failed to execute task 'update-tags'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	}
	if user, err := models.UserByID(db.Db, userId); err != nil {
		logger.Error("Failed to execute task 'update-tags'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	} else if user.AccountType != "tagbot" {
		return nil
	}

	err = updateTagsForUser(ctx, userId)

	if err == nil {
		logger.Info("Task 'update-tags' done.", map[string]interface{}{
			"args": args,
		})
	}
	return err
}

func checkUpdateTagsArguments(args []string) (int, error) {
	if len(args) < 1 {
		return invalidUserID, errors.New("Task 'update-tags' requires at least an integer argument as User ID")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func updateTagsForUser(ctx context.Context, userId int) (err error) {
	var job models.UserUpdateTagsJob
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	if job, err = registerUpdateTagsTask(db.Db, userId); err != nil {
	} else if err = tagging.UpdateTagsForUser(ctx, userId); err != nil {
	} else if err = tagging.UpdateMostUsedTagsForUser(ctx, userId); err != nil {
	} else if err = tagging.UpdateTaggingComplianceForUser(ctx, userId); err != nil {
	}
	updateUpdateTagsTask(db.Db, job, err)
	if err != nil {
		logger.Error("Failed to execute task 'update-tags'.", map[string]interface{}{
			"userId": userId,
			"error":  err.Error(),
		})
	}
	return
}

func registerUpdateTagsTask(db *sql.DB, userId int) (models.UserUpdateTagsJob, error) {
	job := models.UserUpdateTagsJob{
		UserID:   userId,
		WorkerID: backendId,
		Created:  time.Now(),
	}

	err := job.Insert(db)

	return job, err
}

func updateUpdateTagsTask(db *sql.DB, job models.UserUpdateTagsJob, jobError error) error {
	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
