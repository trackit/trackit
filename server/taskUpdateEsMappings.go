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
	"flag"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es/indexes"
	"github.com/trackit/trackit/models"
)

func taskUpdateEsMappings(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'update-es-mappings'.", nil)

	err := updateEsMappings(ctx)
	if err == nil {
		logger.Info("Task 'update-es-mappings' done.", map[string]interface{}{
			"args": args,
		})
	} else {
		logger.Error("Failed to execute task 'update-es-mappings'.", map[string]interface{}{
			"error": err.Error(),
		})
	}
	return err
}

func updateEsMappings(ctx context.Context) (err error) {
	var job models.UpdateEsMappingsJob

	if job, err = registerUpdateEsMappingsTask(db.Db); err != nil {
	} else if err = indexes.UpdateEsMappings(ctx); err != nil {
	}
	updateUpdateEsMappingsTask(db.Db, job, err)
	return
}

func registerUpdateEsMappingsTask(db *sql.DB) (models.UpdateEsMappingsJob, error) {
	job := models.UpdateEsMappingsJob{
		WorkerID: backendId,
		Created:  time.Now(),
	}

	err := job.Insert(db)

	return job, err
}

func updateUpdateEsMappingsTask(db *sql.DB, job models.UpdateEsMappingsJob, jobError error) error {
	job.Completed = time.Now()

	if jobError != nil {
		job.JobError = jobError.Error()
	}

	return job.Update(db)
}
