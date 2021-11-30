package models

import "time"

// MostUsedTagsByUserInRange returns most used tags of a user in a specified range.
func MostUsedTagsByUserInRange(db DB, userId int, begin time.Time, end time.Time) (res []*MostUsedTag, err error) {
	// sql query
	sqlstr := `SELECT ` +
		`id, report_date, user_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE user_id = ? ` +
		`AND report_date >= '` + begin.String() + `' ` +
		`AND report_date < '` + end.String() + `'`

	// run query
	logf(sqlstr, userId)
	q, err := db.Query(sqlstr, userId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := q.Close(); err == nil {
			err = closeErr
		}
	}()

	// load results
	res = []*MostUsedTag{}
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
func MostUsedTagsInUseByUser(db DB, userID int) (*MostUsedTag, error) {
	var err error

	// sql query
	sqlstr := `SELECT ` +
		`id, report_date, user_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE user_id = ? ` +
		`ORDER BY report_date DESC LIMIT 1`

	// run query
	logf(sqlstr, userID)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := q.Close(); err == nil {
			err = closeErr
		}
	}()

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

// MostUsedTagsHistoryByUser returns the currently used most used tags history of a user
func MostUsedTagsHistoryByUser(db DB, userID int) ([]*MostUsedTag, error) {
	var err error

	// sql query
	sqlstr := `SELECT ` +
		`id, report_date, user_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE user_id = ? ` +
		`ORDER BY report_date`

	// run query
	logf(sqlstr, userID)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	res := []*MostUsedTag{}
	for q.Next() {
		mut := MostUsedTag{
			_exists: true,
		}
		err = q.Scan(&mut.ID, &mut.ReportDate, &mut.UserID, &mut.Tags)
		if err != nil {
			return nil, err
		}
		res = append(res, &mut)
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res, nil
}
