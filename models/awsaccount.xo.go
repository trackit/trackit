// Package models contains the types for schema 'trackit'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// AwsAccount represents a row from 'trackit.aws_account'.
type AwsAccount struct {
	ID          int       `json:"id"`           // id
	UserID      int       `json:"user_id"`      // user_id
	Pretty      string    `json:"pretty"`       // pretty
	RoleArn     string    `json:"role_arn"`     // role_arn
	External    string    `json:"external"`     // external
	NextUpdate  time.Time `json:"next_update"`  // next_update
	GraceUpdate time.Time `json:"grace_update"` // grace_update
	Payer       bool      `json:"payer"`        // payer

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the AwsAccount exists in the database.
func (aa *AwsAccount) Exists() bool {
	return aa._exists
}

// Deleted provides information if the AwsAccount has been deleted from the database.
func (aa *AwsAccount) Deleted() bool {
	return aa._deleted
}

// Insert inserts the AwsAccount to the database.
func (aa *AwsAccount) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if aa._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO trackit.aws_account (` +
		`user_id, pretty, role_arn, external, next_update, grace_update, payer` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, aa.UserID, aa.Pretty, aa.RoleArn, aa.External, aa.NextUpdate, aa.GraceUpdate, aa.Payer)
	res, err := db.Exec(sqlstr, aa.UserID, aa.Pretty, aa.RoleArn, aa.External, aa.NextUpdate, aa.GraceUpdate, aa.Payer)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	aa.ID = int(id)
	aa._exists = true

	return nil
}

// Update updates the AwsAccount in the database.
func (aa *AwsAccount) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aa._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if aa._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE trackit.aws_account SET ` +
		`user_id = ?, pretty = ?, role_arn = ?, external = ?, next_update = ?, grace_update = ?, payer = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, aa.UserID, aa.Pretty, aa.RoleArn, aa.External, aa.NextUpdate, aa.GraceUpdate, aa.Payer, aa.ID)
	_, err = db.Exec(sqlstr, aa.UserID, aa.Pretty, aa.RoleArn, aa.External, aa.NextUpdate, aa.GraceUpdate, aa.Payer, aa.ID)
	return err
}

// Save saves the AwsAccount to the database.
func (aa *AwsAccount) Save(db XODB) error {
	if aa.Exists() {
		return aa.Update(db)
	}

	return aa.Insert(db)
}

// Delete deletes the AwsAccount from the database.
func (aa *AwsAccount) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aa._exists {
		return nil
	}

	// if deleted, bail
	if aa._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM trackit.aws_account WHERE id = ?`

	// run query
	XOLog(sqlstr, aa.ID)
	_, err = db.Exec(sqlstr, aa.ID)
	if err != nil {
		return err
	}

	// set deleted
	aa._deleted = true

	return nil
}

// User returns the User associated with the AwsAccount's UserID (user_id).
//
// Generated from foreign key 'aws_account_ibfk_1'.
func (aa *AwsAccount) User(db XODB) (*User, error) {
	return UserByID(db, aa.UserID)
}

// AwsAccountByID retrieves a row from 'trackit.aws_account' as a AwsAccount.
//
// Generated from index 'aws_account_id_pkey'.
func AwsAccountByID(db XODB, id int) (*AwsAccount, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, pretty, role_arn, external, next_update, grace_update, payer ` +
		`FROM trackit.aws_account ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	aa := AwsAccount{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&aa.ID, &aa.UserID, &aa.Pretty, &aa.RoleArn, &aa.External, &aa.NextUpdate, &aa.GraceUpdate, &aa.Payer)
	if err != nil {
		return nil, err
	}

	return &aa, nil
}

// AwsAccountsByUserID retrieves a row from 'trackit.aws_account' as a AwsAccount.
//
// Generated from index 'foreign_user'.
func AwsAccountsByUserID(db XODB, userID int) ([]*AwsAccount, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, pretty, role_arn, external, next_update, grace_update, payer ` +
		`FROM trackit.aws_account ` +
		`WHERE user_id = ?`

	// run query
	XOLog(sqlstr, userID)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*AwsAccount{}
	for q.Next() {
		aa := AwsAccount{
			_exists: true,
		}

		// scan
		err = q.Scan(&aa.ID, &aa.UserID, &aa.Pretty, &aa.RoleArn, &aa.External, &aa.NextUpdate, &aa.GraceUpdate, &aa.Payer)
		if err != nil {
			return nil, err
		}

		res = append(res, &aa)
	}

	return res, nil
}
