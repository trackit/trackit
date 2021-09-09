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

import (
	"fmt"
	"strings"
)

// GetLatestAccountsBillRepositoriesStatus retrieves the most recent job row from 'trackit.aws_bill_update_job' as a AwsAccountStatus.
func GetLatestAccountsBillRepositoriesStatus(db XODB, billRepositoriesIds []int) (accounts map[int]AwsAccountStatus, err error) {
	// sql query
	const sqlstr = `SELECT ` +
		`aws_bill_repository_id, completed, error ` +
		`FROM trackit.aws_account_status ` +
		`WHERE aws_bill_repository_id IN (?)`

	formattedIds := strings.Trim(strings.Replace(fmt.Sprint(billRepositoriesIds), " ", ",", -1), "[]")

	// run query
	XOLog(sqlstr, formattedIds)

	q, err := db.Query(sqlstr, formattedIds)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := q.Close(); err == nil {
			err = closeErr
		}
	}()
	accounts = make(map[int]AwsAccountStatus, 0)
	for q.Next() {
		var account AwsAccountStatus
		var id int
		err = q.Scan(
			&id,
			&account.Completed,
			&account.Error)
		if err != nil {
			return nil, err
		}
		accounts[id] = account
	}
	return
}
