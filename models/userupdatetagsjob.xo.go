package models

// Code generated by xo. DO NOT EDIT.

import (
	"time"
)

// UserUpdateTagsJob represents a row from 'trackit.user_update_tags_job'.
type UserUpdateTagsJob struct {
	ID        int       `json:"id"`        // id
	Created   time.Time `json:"created"`   // created
	UserID    int       `json:"user_id"`   // user_id
	Completed time.Time `json:"completed"` // completed
	WorkerID  string    `json:"worker_id"` // worker_id
	JobError  string    `json:"job_error"` // job_error
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [UserUpdateTagsJob] exists in the database.
func (uutj *UserUpdateTagsJob) Exists() bool {
	return uutj._exists
}

// Deleted returns true when the [UserUpdateTagsJob] has been marked for deletion
// from the database.
func (uutj *UserUpdateTagsJob) Deleted() bool {
	return uutj._deleted
}

// Insert inserts the [UserUpdateTagsJob] to the database.
func (uutj *UserUpdateTagsJob) Insert(db DB) error {
	switch {
	case uutj._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case uutj._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO trackit.user_update_tags_job (` +
		`created, user_id, completed, worker_id, job_error` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, uutj.Created, uutj.UserID, uutj.Completed, uutj.WorkerID, uutj.JobError)
	res, err := db.Exec(sqlstr, uutj.Created, uutj.UserID, uutj.Completed, uutj.WorkerID, uutj.JobError)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	uutj.ID = int(id)
	// set exists
	uutj._exists = true
	return nil
}

// Update updates a [UserUpdateTagsJob] in the database.
func (uutj *UserUpdateTagsJob) Update(db DB) error {
	switch {
	case !uutj._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case uutj._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE trackit.user_update_tags_job SET ` +
		`created = ?, user_id = ?, completed = ?, worker_id = ?, job_error = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, uutj.Created, uutj.UserID, uutj.Completed, uutj.WorkerID, uutj.JobError, uutj.ID)
	if _, err := db.Exec(sqlstr, uutj.Created, uutj.UserID, uutj.Completed, uutj.WorkerID, uutj.JobError, uutj.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [UserUpdateTagsJob] to the database.
func (uutj *UserUpdateTagsJob) Save(db DB) error {
	if uutj.Exists() {
		return uutj.Update(db)
	}
	return uutj.Insert(db)
}

// Upsert performs an upsert for [UserUpdateTagsJob].
func (uutj *UserUpdateTagsJob) Upsert(db DB) error {
	switch {
	case uutj._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO trackit.user_update_tags_job (` +
		`id, created, user_id, completed, worker_id, job_error` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`created = VALUES(created), user_id = VALUES(user_id), completed = VALUES(completed), worker_id = VALUES(worker_id), job_error = VALUES(job_error)`
	// run
	logf(sqlstr, uutj.ID, uutj.Created, uutj.UserID, uutj.Completed, uutj.WorkerID, uutj.JobError)
	if _, err := db.Exec(sqlstr, uutj.ID, uutj.Created, uutj.UserID, uutj.Completed, uutj.WorkerID, uutj.JobError); err != nil {
		return logerror(err)
	}
	// set exists
	uutj._exists = true
	return nil
}

// Delete deletes the [UserUpdateTagsJob] from the database.
func (uutj *UserUpdateTagsJob) Delete(db DB) error {
	switch {
	case !uutj._exists: // doesn't exist
		return nil
	case uutj._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM trackit.user_update_tags_job ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, uutj.ID)
	if _, err := db.Exec(sqlstr, uutj.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	uutj._deleted = true
	return nil
}

// UserUpdateTagsJobByUserID retrieves a row from 'trackit.user_update_tags_job' as a [UserUpdateTagsJob].
//
// Generated from index 'foreign_user'.
func UserUpdateTagsJobByUserID(db DB, userID int) ([]*UserUpdateTagsJob, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, created, user_id, completed, worker_id, job_error ` +
		`FROM trackit.user_update_tags_job ` +
		`WHERE user_id = ?`
	// run
	logf(sqlstr, userID)
	rows, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*UserUpdateTagsJob
	for rows.Next() {
		uutj := UserUpdateTagsJob{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&uutj.ID, &uutj.Created, &uutj.UserID, &uutj.Completed, &uutj.WorkerID, &uutj.JobError); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &uutj)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// UserUpdateTagsJobByID retrieves a row from 'trackit.user_update_tags_job' as a [UserUpdateTagsJob].
//
// Generated from index 'user_update_tags_job_id_pkey'.
func UserUpdateTagsJobByID(db DB, id int) (*UserUpdateTagsJob, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, created, user_id, completed, worker_id, job_error ` +
		`FROM trackit.user_update_tags_job ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	uutj := UserUpdateTagsJob{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, id).Scan(&uutj.ID, &uutj.Created, &uutj.UserID, &uutj.Completed, &uutj.WorkerID, &uutj.JobError); err != nil {
		return nil, logerror(err)
	}
	return &uutj, nil
}

// User returns the User associated with the [UserUpdateTagsJob]'s (UserID).
//
// Generated from foreign key 'user_update_tags_job_ibfk_1'.
func (uutj *UserUpdateTagsJob) User(db DB) (*User, error) {
	return UserByID(db, uutj.UserID)
}
