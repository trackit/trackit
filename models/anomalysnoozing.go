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

func AnomalySnoozingsByUserID(db XODB, userID int) ([]*AnomalySnoozing, error) {
	var err error
	const sqlstr = `SELECT ` +
		`id, user_id, anomaly_id ` +
		`FROM trackit.anomaly_snoozing ` +
		`WHERE user_id = ?`
	XOLog(sqlstr)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var res []*AnomalySnoozing
	for q.Next() {
		as := AnomalySnoozing{
			_exists: true,
		}
		err = q.Scan(&as.ID, &as.UserID, &as.AnomalyID)
		if err != nil {
			return nil, err
		}
		res = append(res, &as)
	}
	return res, nil
}
