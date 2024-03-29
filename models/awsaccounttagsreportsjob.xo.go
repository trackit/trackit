package models

// Code generated by xo. DO NOT EDIT.

import (
	"time"
)

// AwsAccountTagsReportsJob represents a row from 'trackit.aws_account_tags_reports_job'.
type AwsAccountTagsReportsJob struct {
	ID               int       `json:"id"`               // id
	AwsAccountID     int       `json:"aws_account_id"`   // aws_account_id
	Completed        time.Time `json:"completed"`        // completed
	WorkerID         string    `json:"worker_id"`        // worker_id
	JobError         string    `json:"jobError"`         // jobError
	SpreadsheetError string    `json:"spreadsheetError"` // spreadsheetError
	TagsReportError  string    `json:"tagsReportError"`  // tagsReportError
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the AwsAccountTagsReportsJob exists in the database.
func (aatrj *AwsAccountTagsReportsJob) Exists() bool {
	return aatrj._exists
}

// Deleted returns true when the AwsAccountTagsReportsJob has been marked for deletion from
// the database.
func (aatrj *AwsAccountTagsReportsJob) Deleted() bool {
	return aatrj._deleted
}

// Insert inserts the AwsAccountTagsReportsJob to the database.
func (aatrj *AwsAccountTagsReportsJob) Insert(db DB) error {
	switch {
	case aatrj._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case aatrj._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO trackit.aws_account_tags_reports_job (` +
		`aws_account_id, completed, worker_id, jobError, spreadsheetError, tagsReportError` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.JobError, aatrj.SpreadsheetError, aatrj.TagsReportError)
	res, err := db.Exec(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.JobError, aatrj.SpreadsheetError, aatrj.TagsReportError)
	if err != nil {
		return err
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	} // set primary key
	aatrj.ID = int(id)
	// set exists
	aatrj._exists = true
	return nil
}

// Update updates a AwsAccountTagsReportsJob in the database.
func (aatrj *AwsAccountTagsReportsJob) Update(db DB) error {
	switch {
	case !aatrj._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case aatrj._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE trackit.aws_account_tags_reports_job SET ` +
		`aws_account_id = ?, completed = ?, worker_id = ?, jobError = ?, spreadsheetError = ?, tagsReportError = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.JobError, aatrj.SpreadsheetError, aatrj.TagsReportError, aatrj.ID)
	if _, err := db.Exec(sqlstr, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.JobError, aatrj.SpreadsheetError, aatrj.TagsReportError, aatrj.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the AwsAccountTagsReportsJob to the database.
func (aatrj *AwsAccountTagsReportsJob) Save(db DB) error {
	if aatrj.Exists() {
		return aatrj.Update(db)
	}
	return aatrj.Insert(db)
}

// Upsert performs an upsert for AwsAccountTagsReportsJob.
func (aatrj *AwsAccountTagsReportsJob) Upsert(db DB) error {
	switch {
	case aatrj._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO trackit.aws_account_tags_reports_job (` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, tagsReportError` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`aws_account_id = VALUES(aws_account_id), completed = VALUES(completed), worker_id = VALUES(worker_id), jobError = VALUES(jobError), spreadsheetError = VALUES(spreadsheetError), tagsReportError = VALUES(tagsReportError)`
	// run
	logf(sqlstr, aatrj.ID, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.JobError, aatrj.SpreadsheetError, aatrj.TagsReportError)
	if _, err := db.Exec(sqlstr, aatrj.ID, aatrj.AwsAccountID, aatrj.Completed, aatrj.WorkerID, aatrj.JobError, aatrj.SpreadsheetError, aatrj.TagsReportError); err != nil {
		return err
	}
	// set exists
	aatrj._exists = true
	return nil
}

// Delete deletes the AwsAccountTagsReportsJob from the database.
func (aatrj *AwsAccountTagsReportsJob) Delete(db DB) error {
	switch {
	case !aatrj._exists: // doesn't exist
		return nil
	case aatrj._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM trackit.aws_account_tags_reports_job ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, aatrj.ID)
	if _, err := db.Exec(sqlstr, aatrj.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	aatrj._deleted = true
	return nil
}

// AwsAccountTagsReportsJobByID retrieves a row from 'trackit.aws_account_tags_reports_job' as a AwsAccountTagsReportsJob.
//
// Generated from index 'aws_account_tags_reports_job_id_pkey'.
func AwsAccountTagsReportsJobByID(db DB, id int) (*AwsAccountTagsReportsJob, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, tagsReportError ` +
		`FROM trackit.aws_account_tags_reports_job ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	aatrj := AwsAccountTagsReportsJob{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, id).Scan(&aatrj.ID, &aatrj.AwsAccountID, &aatrj.Completed, &aatrj.WorkerID, &aatrj.JobError, &aatrj.SpreadsheetError, &aatrj.TagsReportError); err != nil {
		return nil, logerror(err)
	}
	return &aatrj, nil
}

// AwsAccountTagsReportsJobByAwsAccountID retrieves a row from 'trackit.aws_account_tags_reports_job' as a AwsAccountTagsReportsJob.
//
// Generated from index 'foreign_aws_account'.
func AwsAccountTagsReportsJobByAwsAccountID(db DB, awsAccountID int) ([]*AwsAccountTagsReportsJob, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, tagsReportError ` +
		`FROM trackit.aws_account_tags_reports_job ` +
		`WHERE aws_account_id = ?`
	// run
	logf(sqlstr, awsAccountID)
	rows, err := db.Query(sqlstr, awsAccountID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*AwsAccountTagsReportsJob
	for rows.Next() {
		aatrj := AwsAccountTagsReportsJob{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&aatrj.ID, &aatrj.AwsAccountID, &aatrj.Completed, &aatrj.WorkerID, &aatrj.JobError, &aatrj.SpreadsheetError, &aatrj.TagsReportError); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &aatrj)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// AwsAccount returns the AwsAccount associated with the AwsAccountTagsReportsJob's (AwsAccountID).
//
// Generated from foreign key 'aws_account_tags_reports_job_ibfk_1'.
func (aatrj *AwsAccountTagsReportsJob) AwsAccount(db DB) (*AwsAccount, error) {
	return AwsAccountByID(db, aatrj.AwsAccountID)
}
