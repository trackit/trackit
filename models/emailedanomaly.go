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

import "time"

// IsAnomalyAlreadyEmailed checks if an anomaly has already been sent.
func IsAnomalyAlreadyEmailed(db XODB, awsAccountId int, product string, date time.Time) (result bool, err error) {
	const sqlstr = `SELECT ` +
		`id, aws_account_id, product, recipient, date ` +
		`FROM trackit.emailed_anomaly ` +
		`WHERE aws_account_id = ? AND product = ? AND date = ?`
	XOLog(sqlstr, awsAccountId, product, date)
	q, err := db.Query(sqlstr, awsAccountId, product, date)
	if err != nil {
		return false, err
	}
	defer func() {
		if closeErr := q.Close(); err == nil {
			err = closeErr
		}
	}()
	if q.Next() {
		return true, nil
	}
	return false, nil
}
