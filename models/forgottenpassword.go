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
	"time"
)

// DeleteExpiredForgottenPassword deletes the ForgottenPassword that are older than the date parameter
func DeleteExpiredForgottenPassword(db DB, date time.Time) error {
	var err error

	// sql query
	const sqlstr = `DELETE FROM trackit.forgotten_password WHERE created < ?`

	// run query
	logf(sqlstr, date)
	_, err = db.Exec(sqlstr, date)
	if err != nil {
		return err
	}

	return nil
}
