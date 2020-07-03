// Package models contains the types for schema 'trackit'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// MostUsedTag represents a row from 'trackit.most_used_tags'.
type MostUsedTag struct {
	ID           int       `json:"id"`             // id
	ReportDate   time.Time `json:"report_date"`    // report_date
	AwsAccountID int       `json:"aws_account_id"` // aws_account_id
	Tags         string    `json:"tags"`           // tags

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the MostUsedTag exists in the database.
func (mut *MostUsedTag) Exists() bool {
	return mut._exists
}

// Deleted provides information if the MostUsedTag has been deleted from the database.
func (mut *MostUsedTag) Deleted() bool {
	return mut._deleted
}

// Insert inserts the MostUsedTag to the database.
func (mut *MostUsedTag) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if mut._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO trackit.most_used_tags (` +
		`report_date, aws_account_id, tags` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, mut.ReportDate, mut.AwsAccountID, mut.Tags)
	res, err := db.Exec(sqlstr, mut.ReportDate, mut.AwsAccountID, mut.Tags)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	mut.ID = int(id)
	mut._exists = true

	return nil
}

// Update updates the MostUsedTag in the database.
func (mut *MostUsedTag) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !mut._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if mut._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE trackit.most_used_tags SET ` +
		`report_date = ?, aws_account_id = ?, tags = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, mut.ReportDate, mut.AwsAccountID, mut.Tags, mut.ID)
	_, err = db.Exec(sqlstr, mut.ReportDate, mut.AwsAccountID, mut.Tags, mut.ID)
	return err
}

// Save saves the MostUsedTag to the database.
func (mut *MostUsedTag) Save(db XODB) error {
	if mut.Exists() {
		return mut.Update(db)
	}

	return mut.Insert(db)
}

// Delete deletes the MostUsedTag from the database.
func (mut *MostUsedTag) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !mut._exists {
		return nil
	}

	// if deleted, bail
	if mut._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM trackit.most_used_tags WHERE id = ?`

	// run query
	XOLog(sqlstr, mut.ID)
	_, err = db.Exec(sqlstr, mut.ID)
	if err != nil {
		return err
	}

	// set deleted
	mut._deleted = true

	return nil
}

// AwsAccount returns the AwsAccount associated with the MostUsedTag's AwsAccountID (aws_account_id).
//
// Generated from foreign key 'most_used_tags_ibfk_1'.
func (mut *MostUsedTag) AwsAccount(db XODB) (*AwsAccount, error) {
	return AwsAccountByID(db, mut.AwsAccountID)
}

// MostUsedTagsByAwsAccountID retrieves a row from 'trackit.most_used_tags' as a MostUsedTag.
//
// Generated from index 'foreign_aws_account'.
func MostUsedTagsByAwsAccountID(db XODB, awsAccountID int) ([]*MostUsedTag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, report_date, aws_account_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE aws_account_id = ?`

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
		err = q.Scan(&mut.ID, &mut.ReportDate, &mut.AwsAccountID, &mut.Tags)
		if err != nil {
			return nil, err
		}

		res = append(res, &mut)
	}

	return res, nil
}

// MostUsedTagByID retrieves a row from 'trackit.most_used_tags' as a MostUsedTag.
//
// Generated from index 'most_used_tags_id_pkey'.
func MostUsedTagByID(db XODB, id int) (*MostUsedTag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, report_date, aws_account_id, tags ` +
		`FROM trackit.most_used_tags ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	mut := MostUsedTag{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&mut.ID, &mut.ReportDate, &mut.AwsAccountID, &mut.Tags)
	if err != nil {
		return nil, err
	}

	return &mut, nil
}