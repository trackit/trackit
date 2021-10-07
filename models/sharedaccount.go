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

// SharedAccountWithRole represent a row from 'trackit.shared_account' joined with
// the role_arn from 'trackit.aws_account'
type SharedAccountWithRole struct {
	ID              int    `json:"id"`               // id
	AccountID       int    `json:"account_id"`       // account_id
	UserID          int    `json:"user_id"`          // user_id
	UserPermission  int    `json:"user_permission"`  // user_permission
	SharingAccepted bool   `json:"sharing_accepted"` // sharing_accepted
	RoleArn         string `json:"role_arn"`         // role_arn
	AwsIdentity     string `json:"aws_identity"`     // aws_identity
	OwnerID         int    `json:"owner_id"`         // owner_id
}

// SharedAccountsWithRoleByUserID returns all the shared accounts and their role arn
// for a user
func SharedAccountsWithRoleByUserID(db DB, userID int) ([]*SharedAccountWithRole, error) {
	var err error
	const sqlstr = `SELECT ` +
		`sa.id, sa.account_id, sa.user_id, sa.user_permission, sa.sharing_accepted, aa.role_arn, aa.aws_identity, aa.user_id ` +
		`FROM trackit.shared_account AS sa ` +
		`INNER JOIN trackit.aws_account AS aa ON sa.account_id=aa.id ` +
		`WHERE sa.user_id=?`
	logf(sqlstr)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	res := []*SharedAccountWithRole{}
	for q.Next() {
		sa := SharedAccountWithRole{}
		err = q.Scan(&sa.ID, &sa.AccountID, &sa.UserID, &sa.UserPermission, &sa.SharingAccepted, &sa.RoleArn, &sa.AwsIdentity, &sa.OwnerID)
		if err != nil {
			return nil, err
		}
		res = append(res, &sa)
	}
	return res, nil
}
