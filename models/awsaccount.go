//   Copyright 2018 MSolution.IO
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
func AwsAccounts(db XODB) ([]*AwsAccount, error) {
	var err error
	const sqlstr = `SELECT ` +
		`id, user_id, pretty, role_arn, external, next_update, grace_update ` +
		`FROM trackit.aws_account`
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var res []*AwsAccount
	for q.Next() {
		aa := AwsAccount{
			_exists: true,
		}
		err = q.Scan(&aa.ID, &aa.UserID, &aa.Pretty, &aa.RoleArn, &aa.External, &aa.NextUpdate, &aa.GraceUpdate)
		if err != nil {
			return nil, err
		}
		res = append(res, &aa)
	}
	return res, nil
}
