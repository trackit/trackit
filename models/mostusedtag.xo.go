package models

// Code generated by xo. DO NOT EDIT.

import (
	"time"
)

// MostUsedTag represents a row from 'trackit.most_used_tags'.
type MostUsedTag struct {
	ID         int       `json:"id"`          // id
	ReportDate time.Time `json:"report_date"` // report_date
	UserID     int       `json:"user_id"`     // user_id
	Tags       string    `json:"tags"`        // tags
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [MostUsedTag] exists in the database.
func (mut *MostUsedTag) Exists() bool {
	return mut._exists
}

// Deleted returns true when the [MostUsedTag] has been marked for deletion
// from the database.
func (mut *MostUsedTag) Deleted() bool {
	return mut._deleted
}

// Insert inserts the [MostUsedTag] to the database.
func (mut *MostUsedTag) Insert(db DB) error {
	switch {
	case mut._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case mut._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO trackit.most_used_tags (` +
		`report_date, user_id, tags` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`
	// run
	logf(sqlstr, mut.ReportDate, mut.UserID, mut.Tags)
	res, err := db.Exec(sqlstr, mut.ReportDate, mut.UserID, mut.Tags)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	mut.ID = int(id)
	// set exists
	mut._exists = true
	return nil
}

// Update updates a [MostUsedTag] in the database.
func (mut *MostUsedTag) Update(db DB) error {
	switch {
	case !mut._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case mut._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE trackit.most_used_tags SET ` +
		`report_date = ?, user_id = ?, tags = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, mut.ReportDate, mut.UserID, mut.Tags, mut.ID)
	if _, err := db.Exec(sqlstr, mut.ReportDate, mut.UserID, mut.Tags, mut.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [MostUsedTag] to the database.
func (mut *MostUsedTag) Save(db DB) error {
	if mut.Exists() {
		return mut.Update(db)
	}
	return mut.Insert(db)
}

// Upsert performs an upsert for [MostUsedTag].
func (mut *MostUsedTag) Upsert(db DB) error {
	switch {
	case mut._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO trackit.most_used_tags (` +
		`id, report_date, user_id, tags` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`report_date = VALUES(report_date), user_id = VALUES(user_id), tags = VALUES(tags)`
	// run
	logf(sqlstr, mut.ID, mut.ReportDate, mut.UserID, mut.Tags)
	if _, err := db.Exec(sqlstr, mut.ID, mut.ReportDate, mut.UserID, mut.Tags); err != nil {
		return logerror(err)
	}
	// set exists
	mut._exists = true
	return nil
}

// Delete deletes the [MostUsedTag] from the database.
func (mut *MostUsedTag) Delete(db DB) error {
	switch {
	case !mut._exists: // doesn't exist
		return nil
	case mut._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM trackit.most_used_tags ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, mut.ID)
	if _, err := db.Exec(sqlstr, mut.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	mut._deleted = true
	return nil
}

// MostUsedTagsByUserID retrieves a row from 'trackit.most_used_tags' as a [MostUsedTag].
//
// Generated from index 'foreign_user'.
func MostUsedTagsByUserID(db DB, userID int) ([]*MostUsedTag, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, report_date, user_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE user_id = ?`
	// run
	logf(sqlstr, userID)
	rows, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*MostUsedTag
	for rows.Next() {
		mut := MostUsedTag{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&mut.ID, &mut.ReportDate, &mut.UserID, &mut.Tags); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &mut)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// MostUsedTagByID retrieves a row from 'trackit.most_used_tags' as a [MostUsedTag].
//
// Generated from index 'most_used_tags_id_pkey'.
func MostUsedTagByID(db DB, id int) (*MostUsedTag, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, report_date, user_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	mut := MostUsedTag{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, id).Scan(&mut.ID, &mut.ReportDate, &mut.UserID, &mut.Tags); err != nil {
		return nil, logerror(err)
	}
	return &mut, nil
}

// User returns the User associated with the [MostUsedTag]'s (UserID).
//
// Generated from foreign key 'most_used_tags_ibfk_1'.
func (mut *MostUsedTag) User(db DB) (*User, error) {
	return UserByID(db, mut.UserID)
}
