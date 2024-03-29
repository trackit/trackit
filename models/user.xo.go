package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"time"
)

// User represents a row from 'trackit.user'.
type User struct {
	ID                     int            `json:"id"`                       // id
	Email                  string         `json:"email"`                    // email
	Auth                   string         `json:"auth"`                     // auth
	NextExternal           sql.NullString `json:"next_external"`            // next_external
	ParentUserID           sql.NullInt64  `json:"parent_user_id"`           // parent_user_id
	AwsCustomerIdentifier  string         `json:"aws_customer_identifier"`  // aws_customer_identifier
	AwsCustomerEntitlement bool           `json:"aws_customer_entitlement"` // aws_customer_entitlement
	NextUpdateEntitlement  time.Time      `json:"next_update_entitlement"`  // next_update_entitlement
	AnomaliesFilters       []byte         `json:"anomalies_filters"`        // anomalies_filters
	NextUpdateTags         time.Time      `json:"next_update_tags"`         // next_update_tags
	LastSeen               time.Time      `json:"last_seen"`                // last_seen
	LastUnusedReminder     time.Time      `json:"last_unused_reminder"`     // last_unused_reminder
	LastUnusedSlack        time.Time      `json:"last_unused_slack"`        // last_unused_slack
	AccountType            string         `json:"account_type"`             // account_type
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the User exists in the database.
func (u *User) Exists() bool {
	return u._exists
}

// Deleted returns true when the User has been marked for deletion from
// the database.
func (u *User) Deleted() bool {
	return u._deleted
}

// Insert inserts the User to the database.
func (u *User) Insert(db DB) error {
	switch {
	case u._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case u._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO trackit.user (` +
		`email, auth, next_external, parent_user_id, aws_customer_identifier, aws_customer_entitlement, next_update_entitlement, anomalies_filters, next_update_tags, last_seen, last_unused_reminder, last_unused_slack, account_type` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, u.Email, u.Auth, u.NextExternal, u.ParentUserID, u.AwsCustomerIdentifier, u.AwsCustomerEntitlement, u.NextUpdateEntitlement, u.AnomaliesFilters, u.NextUpdateTags, u.LastSeen, u.LastUnusedReminder, u.LastUnusedSlack, u.AccountType)
	res, err := db.Exec(sqlstr, u.Email, u.Auth, u.NextExternal, u.ParentUserID, u.AwsCustomerIdentifier, u.AwsCustomerEntitlement, u.NextUpdateEntitlement, u.AnomaliesFilters, u.NextUpdateTags, u.LastSeen, u.LastUnusedReminder, u.LastUnusedSlack, u.AccountType)
	if err != nil {
		return err
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	} // set primary key
	u.ID = int(id)
	// set exists
	u._exists = true
	return nil
}

// Update updates a User in the database.
func (u *User) Update(db DB) error {
	switch {
	case !u._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case u._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE trackit.user SET ` +
		`email = ?, auth = ?, next_external = ?, parent_user_id = ?, aws_customer_identifier = ?, aws_customer_entitlement = ?, next_update_entitlement = ?, anomalies_filters = ?, next_update_tags = ?, last_seen = ?, last_unused_reminder = ?, last_unused_slack = ?, account_type = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, u.Email, u.Auth, u.NextExternal, u.ParentUserID, u.AwsCustomerIdentifier, u.AwsCustomerEntitlement, u.NextUpdateEntitlement, u.AnomaliesFilters, u.NextUpdateTags, u.LastSeen, u.LastUnusedReminder, u.LastUnusedSlack, u.AccountType, u.ID)
	if _, err := db.Exec(sqlstr, u.Email, u.Auth, u.NextExternal, u.ParentUserID, u.AwsCustomerIdentifier, u.AwsCustomerEntitlement, u.NextUpdateEntitlement, u.AnomaliesFilters, u.NextUpdateTags, u.LastSeen, u.LastUnusedReminder, u.LastUnusedSlack, u.AccountType, u.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the User to the database.
func (u *User) Save(db DB) error {
	if u.Exists() {
		return u.Update(db)
	}
	return u.Insert(db)
}

// Upsert performs an upsert for User.
func (u *User) Upsert(db DB) error {
	switch {
	case u._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO trackit.user (` +
		`id, email, auth, next_external, parent_user_id, aws_customer_identifier, aws_customer_entitlement, next_update_entitlement, anomalies_filters, next_update_tags, last_seen, last_unused_reminder, last_unused_slack, account_type` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`email = VALUES(email), auth = VALUES(auth), next_external = VALUES(next_external), parent_user_id = VALUES(parent_user_id), aws_customer_identifier = VALUES(aws_customer_identifier), aws_customer_entitlement = VALUES(aws_customer_entitlement), next_update_entitlement = VALUES(next_update_entitlement), anomalies_filters = VALUES(anomalies_filters), next_update_tags = VALUES(next_update_tags), last_seen = VALUES(last_seen), last_unused_reminder = VALUES(last_unused_reminder), last_unused_slack = VALUES(last_unused_slack), account_type = VALUES(account_type)`
	// run
	logf(sqlstr, u.ID, u.Email, u.Auth, u.NextExternal, u.ParentUserID, u.AwsCustomerIdentifier, u.AwsCustomerEntitlement, u.NextUpdateEntitlement, u.AnomaliesFilters, u.NextUpdateTags, u.LastSeen, u.LastUnusedReminder, u.LastUnusedSlack, u.AccountType)
	if _, err := db.Exec(sqlstr, u.ID, u.Email, u.Auth, u.NextExternal, u.ParentUserID, u.AwsCustomerIdentifier, u.AwsCustomerEntitlement, u.NextUpdateEntitlement, u.AnomaliesFilters, u.NextUpdateTags, u.LastSeen, u.LastUnusedReminder, u.LastUnusedSlack, u.AccountType); err != nil {
		return err
	}
	// set exists
	u._exists = true
	return nil
}

// Delete deletes the User from the database.
func (u *User) Delete(db DB) error {
	switch {
	case !u._exists: // doesn't exist
		return nil
	case u._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM trackit.user ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, u.ID)
	if _, err := db.Exec(sqlstr, u.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	u._deleted = true
	return nil
}

// UserByParentUserID retrieves a row from 'trackit.user' as a User.
//
// Generated from index 'parent_user'.
func UserByParentUserID(db DB, parentUserID sql.NullInt64) ([]*User, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, email, auth, next_external, parent_user_id, aws_customer_identifier, aws_customer_entitlement, next_update_entitlement, anomalies_filters, next_update_tags, last_seen, last_unused_reminder, last_unused_slack, account_type ` +
		`FROM trackit.user ` +
		`WHERE parent_user_id = ?`
	// run
	logf(sqlstr, parentUserID)
	rows, err := db.Query(sqlstr, parentUserID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*User
	for rows.Next() {
		u := User{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&u.ID, &u.Email, &u.Auth, &u.NextExternal, &u.ParentUserID, &u.AwsCustomerIdentifier, &u.AwsCustomerEntitlement, &u.NextUpdateEntitlement, &u.AnomaliesFilters, &u.NextUpdateTags, &u.LastSeen, &u.LastUnusedReminder, &u.LastUnusedSlack, &u.AccountType); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// UserByEmailAccountType retrieves a row from 'trackit.user' as a User.
//
// Generated from index 'unique_email_account_type'.
func UserByEmailAccountType(db DB, email, accountType string) (*User, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, email, auth, next_external, parent_user_id, aws_customer_identifier, aws_customer_entitlement, next_update_entitlement, anomalies_filters, next_update_tags, last_seen, last_unused_reminder, last_unused_slack, account_type ` +
		`FROM trackit.user ` +
		`WHERE email = ? AND account_type = ?`
	// run
	logf(sqlstr, email, accountType)
	u := User{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, email, accountType).Scan(&u.ID, &u.Email, &u.Auth, &u.NextExternal, &u.ParentUserID, &u.AwsCustomerIdentifier, &u.AwsCustomerEntitlement, &u.NextUpdateEntitlement, &u.AnomaliesFilters, &u.NextUpdateTags, &u.LastSeen, &u.LastUnusedReminder, &u.LastUnusedSlack, &u.AccountType); err != nil {
		return nil, logerror(err)
	}
	return &u, nil
}

// UserByID retrieves a row from 'trackit.user' as a User.
//
// Generated from index 'user_id_pkey'.
func UserByID(db DB, id int) (*User, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, email, auth, next_external, parent_user_id, aws_customer_identifier, aws_customer_entitlement, next_update_entitlement, anomalies_filters, next_update_tags, last_seen, last_unused_reminder, last_unused_slack, account_type ` +
		`FROM trackit.user ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	u := User{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, id).Scan(&u.ID, &u.Email, &u.Auth, &u.NextExternal, &u.ParentUserID, &u.AwsCustomerIdentifier, &u.AwsCustomerEntitlement, &u.NextUpdateEntitlement, &u.AnomaliesFilters, &u.NextUpdateTags, &u.LastSeen, &u.LastUnusedReminder, &u.LastUnusedSlack, &u.AccountType); err != nil {
		return nil, logerror(err)
	}
	return &u, nil
}

// User returns the User associated with the User's (ParentUserID).
//
// Generated from foreign key 'parent_user'.
func (u *User) User(db DB) (*User, error) {
	return UserByID(db, int(u.ParentUserID.Int64))
}
