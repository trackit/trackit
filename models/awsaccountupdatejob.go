//   Copyright 2017 MSolution.IO
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

// GetLatestAccountUpdateJob retrieves the most recent completed row from 'trackit.aws_account_update_job' as a AwsAccountUpdateJob.
func GetLatestAccountUpdateJob(db XODB, accountId int) (*AwsAccountUpdateJob, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, completed, worker_id, jobError, rdsError, monthly_reports_generated ` +
		`FROM trackit.aws_account_update_job ` +
		`WHERE aws_account_id = ? ORDER BY completed DESC LIMIT 1`

	// run query
	XOLog(sqlstr, accountId)
	aauj := AwsAccountUpdateJob{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, accountId).Scan(&aauj.ID, &aauj.AwsAccountID, &aauj.Completed, &aauj.WorkerID, &aauj.Joberror, &aauj.Rdserror, &aauj.MonthlyReportsGenerated)
	if err != nil {
		return nil, err
	}

	return &aauj, nil
}
