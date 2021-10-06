//   Copyright 2021 MSolution.IO
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
package models

import "time"

// GetUnusedAccounts returns all user not seen for the duration passed in arguments
func GetUnusedAccounts(db DB, unseenFor time.Duration) (res []*User, err error) {
	limitTime := time.Now().Add(-unseenFor)

	// sql query
	const sqlstr = `SELECT ` +
		`id, email, auth, next_external, parent_user_id, aws_customer_identifier, aws_customer_entitlement, next_update_entitlement, anomalies_filters, last_seen, last_unused_reminder, account_type ` +
		`FROM trackit.user ` +
		`WHERE last_seen < ?`

	// run query
	logf(sqlstr, limitTime.Format(time.RFC3339))
	q, err := db.Query(sqlstr, limitTime.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := q.Close(); err == nil {
			err = closeErr
		}
	}()

	// load results
	res = []*User{}
	for q.Next() {
		u := User{
			_exists: true,
		}

		// scan
		err = q.Scan(&u.ID, &u.Email, &u.Auth, &u.NextExternal, &u.ParentUserID, &u.AwsCustomerIdentifier, &u.AwsCustomerEntitlement, &u.NextUpdateEntitlement, &u.AnomaliesFilters, &u.LastSeen, &u.LastUnusedReminder, &u.AccountType)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}
