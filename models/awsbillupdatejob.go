//   Copyright 2019 MSolution.IO
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

// Package models contains the types for schema 'trackit'.
package models

// LastAwsBillUpdateJobsByAwsBillRepositoryID returns the last bill update job.
func LastAwsBillUpdateJobsByAwsBillRepositoryID(db DB, awsBillRepositoryID int) (*AwsBillUpdateJob, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, aws_bill_repository_id, expired, completed, worker_id, error ` +
		`FROM trackit.aws_bill_update_job ` +
		`WHERE aws_bill_repository_id = ? ` +
		`ORDER BY id DESC LIMIT 1`

	// run query
	logf(sqlstr, awsBillRepositoryID)
	abuj := AwsBillUpdateJob{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, awsBillRepositoryID).Scan(&abuj.ID, &abuj.AwsBillRepositoryID, &abuj.Expired, &abuj.Completed, &abuj.WorkerID, &abuj.Error)
	if err != nil {
		return nil, err
	}

	return &abuj, nil
}
