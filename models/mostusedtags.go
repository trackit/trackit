package models

import "time"

// MostUsedTagsByUserInRange returns most used tags of a user in a specified range.
func MostUsedTagsByUserInRange(db XODB, userId int, begin time.Time, end time.Time) ([]*MostUsedTag, error) {
	var err error

	// sql query
	sqlstr := `SELECT ` +
		`id, report_date, user_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE user_id = ? ` +
		`AND report_date >= '` + begin.String() + `' ` +
		`AND report_date < '` + end.String() + `'`

	// run query
	XOLog(sqlstr, userId)
	q, err := db.Query(sqlstr, userId)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*MostUsedTag{}
	for q.Next() {
		mut := MostUsedTag{
			_exists: true,
		}

		// scan
		err = q.Scan(&mut.ID, &mut.ReportDate, &mut.UserID, &mut.Tags)
		if err != nil {
			return nil, err
		}

		res = append(res, &mut)
	}

	return res, nil
}

// MostUsedTagsInUseByUser returns the currently used most used tags of a user
func MostUsedTagsInUseByUser(db XODB, awsAccountID int) (*MostUsedTag, error) {
	var err error

	// sql query
	sqlstr := `SELECT ` +
		`id, report_date, user_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE user_id = ? ` +
		`ORDER BY report_date DESC LIMIT 1`

	// run query
	XOLog(sqlstr, awsAccountID)
	q, err := db.Query(sqlstr, awsAccountID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*MostUsedTag{}
	for q.Next() {
		mut := MostUsedTag{
			_exists: true,
		}

		// scan
		err = q.Scan(&mut.ID, &mut.ReportDate, &mut.UserID, &mut.Tags)
		if err != nil {
			return nil, err
		}

		res = append(res, &mut)
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res[0], nil
}
