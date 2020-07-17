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

func taskUpdateMostUsedTags(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'update-most-used-tags'.", map[string]interface{}{
		"args": args,
	})

	userId, err := checkUpdateMostUsedTagsArguments(args)
	if err != nil {
		logger.Error("Failed to execute task 'update-most-used-tags'.", map[string]interface{}{
			"err": err.Error(),
		})
		return err
	}

	err = updateMostUsedTagsForUser(ctx, userId)

	if err != nil {
		logger.Info("Task 'update-most-used-tags' done.", map[string]interface{}{
			"args": args,
		})
	}
	return err
}

func checkUpdateMostUsedTagsArguments(args []string) (int, error) {
	if len(args) < 1 {
		return invalidUserID, errors.New("Task 'update-most-used-tags' requires at least an integer argument as User ID")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return invalidUserID, err
	}

	return userId, nil
}

func updateMostUsedTagsForUser(ctx context.Context, userId int) (err error) {
	var job models.UserUpdateMostUsedTagsJob
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	if job, err = registerUpdateMostUsedTagsTask(db.Db, userId); err != nil {
	} else {
		err = tagging.UpdateMostUsedTagsForUser(ctx, userId)
		updateUpdateMostUsedTagsTask(db.Db, job, err)
	}
	if err != nil {
		logger.Error("Failed to execute task 'update-most-used-tags'.", map[string]interface{}{
			"userId": userId,
			"error":  err.Error(),
		})
	}
	return
}

func registerUpdateMostUsedTagsTask(db *sql.DB, userId int) (models.UserUpdateMostUsedTagsJob, error) {
	job := models.UserUpdateMostUsedTagsJob{
		UserID:   userId,
		WorkerID: backendId,
		Created:  time.Now(),
	}

	err := job.Insert(db)

	return job, err
}

func updateUpdateMostUsedTagsTask(db *sql.DB, job models.UserUpdateMostUsedTagsJob, jobError error) error {
	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
