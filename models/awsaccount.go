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

// AwsAccounts returns the set of aws account
func AwsAccounts(db XODB) (res []*AwsAccount, err error) {
	const sqlstr = `SELECT ` +
		`id, user_id, pretty, role_arn, external, next_update, aws_identity ` +
		`FROM trackit.aws_account`
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := q.Close(); err == nil {
			err = closeErr
		}
	}()
	for q.Next() {
		aa := AwsAccount{
			_exists: true,
		}
		err = q.Scan(&aa.ID, &aa.UserID, &aa.Pretty, &aa.RoleArn, &aa.External, &aa.NextUpdate, &aa.AwsIdentity)
		if err != nil {
			return nil, err
		}
		res = append(res, &aa)
	}
	return res, nil
}

func AwsAccountsByParentId(db XODB, parentID int) (res []*AwsAccount, err error) {
	const sqlstr = `SELECT ` +
		`id, user_id, pretty, role_arn, external, next_update, aws_identity, last_spreadsheet_report_generation ` +
		`FROM trackit.aws_account ` +
		`WHERE parent_id = ?`
	XOLog(sqlstr)
	q, err := db.Query(sqlstr, parentID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := q.Close(); err == nil {
			err = closeErr
		}
	}()
	for q.Next() {
		aa := AwsAccount{
			_exists: true,
		}
		err = q.Scan(&aa.ID, &aa.UserID, &aa.Pretty, &aa.RoleArn, &aa.External, &aa.NextUpdate, &aa.AwsIdentity, &aa.LastSpreadsheetReportGeneration)
		if err != nil {
			return nil, err
		}
		res = append(res, &aa)
	}
	return res, nil
}
