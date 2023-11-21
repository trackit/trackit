package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"time"
)

// TagbotUser represents a row from 'trackit.tagbot_user'.
type TagbotUser struct {
	ID                            int           `json:"id"`                               // id
	UserID                        int           `json:"user_id"`                          // user_id
	AwsCustomerIdentifier         string        `json:"aws_customer_identifier"`          // aws_customer_identifier
	AwsCustomerEntitlement        bool          `json:"aws_customer_entitlement"`         // aws_customer_entitlement
	StripeCustomerIdentifier      string        `json:"stripe_customer_identifier"`       // stripe_customer_identifier
	StripeCustomerEntitlement     bool          `json:"stripe_customer_entitlement"`      // stripe_customer_entitlement
	StripeSubscriptionIdentifier  string        `json:"stripe_subscription_identifier"`   // stripe_subscription_identifier
	StripePaymentMethodIdentifier string        `json:"stripe_payment_method_identifier"` // stripe_payment_method_identifier
	FreeTierEndAt                 time.Time     `json:"free_tier_end_at"`                 // free_tier_end_at
	DiscountCodeID                sql.NullInt64 `json:"discount_code_id"`                 // discount_code_id
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [TagbotUser] exists in the database.
func (tu *TagbotUser) Exists() bool {
	return tu._exists
}

// Deleted returns true when the [TagbotUser] has been marked for deletion
// from the database.
func (tu *TagbotUser) Deleted() bool {
	return tu._deleted
}

// Insert inserts the [TagbotUser] to the database.
func (tu *TagbotUser) Insert(db DB) error {
	switch {
	case tu._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case tu._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO trackit.tagbot_user (` +
		`user_id, aws_customer_identifier, aws_customer_entitlement, stripe_customer_identifier, stripe_customer_entitlement, stripe_subscription_identifier, stripe_payment_method_identifier, free_tier_end_at, discount_code_id` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, tu.UserID, tu.AwsCustomerIdentifier, tu.AwsCustomerEntitlement, tu.StripeCustomerIdentifier, tu.StripeCustomerEntitlement, tu.StripeSubscriptionIdentifier, tu.StripePaymentMethodIdentifier, tu.FreeTierEndAt, tu.DiscountCodeID)
	res, err := db.Exec(sqlstr, tu.UserID, tu.AwsCustomerIdentifier, tu.AwsCustomerEntitlement, tu.StripeCustomerIdentifier, tu.StripeCustomerEntitlement, tu.StripeSubscriptionIdentifier, tu.StripePaymentMethodIdentifier, tu.FreeTierEndAt, tu.DiscountCodeID)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	tu.ID = int(id)
	// set exists
	tu._exists = true
	return nil
}

// Update updates a [TagbotUser] in the database.
func (tu *TagbotUser) Update(db DB) error {
	switch {
	case !tu._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case tu._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE trackit.tagbot_user SET ` +
		`user_id = ?, aws_customer_identifier = ?, aws_customer_entitlement = ?, stripe_customer_identifier = ?, stripe_customer_entitlement = ?, stripe_subscription_identifier = ?, stripe_payment_method_identifier = ?, free_tier_end_at = ?, discount_code_id = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, tu.UserID, tu.AwsCustomerIdentifier, tu.AwsCustomerEntitlement, tu.StripeCustomerIdentifier, tu.StripeCustomerEntitlement, tu.StripeSubscriptionIdentifier, tu.StripePaymentMethodIdentifier, tu.FreeTierEndAt, tu.DiscountCodeID, tu.ID)
	if _, err := db.Exec(sqlstr, tu.UserID, tu.AwsCustomerIdentifier, tu.AwsCustomerEntitlement, tu.StripeCustomerIdentifier, tu.StripeCustomerEntitlement, tu.StripeSubscriptionIdentifier, tu.StripePaymentMethodIdentifier, tu.FreeTierEndAt, tu.DiscountCodeID, tu.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [TagbotUser] to the database.
func (tu *TagbotUser) Save(db DB) error {
	if tu.Exists() {
		return tu.Update(db)
	}
	return tu.Insert(db)
}

// Upsert performs an upsert for [TagbotUser].
func (tu *TagbotUser) Upsert(db DB) error {
	switch {
	case tu._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO trackit.tagbot_user (` +
		`id, user_id, aws_customer_identifier, aws_customer_entitlement, stripe_customer_identifier, stripe_customer_entitlement, stripe_subscription_identifier, stripe_payment_method_identifier, free_tier_end_at, discount_code_id` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`user_id = VALUES(user_id), aws_customer_identifier = VALUES(aws_customer_identifier), aws_customer_entitlement = VALUES(aws_customer_entitlement), stripe_customer_identifier = VALUES(stripe_customer_identifier), stripe_customer_entitlement = VALUES(stripe_customer_entitlement), stripe_subscription_identifier = VALUES(stripe_subscription_identifier), stripe_payment_method_identifier = VALUES(stripe_payment_method_identifier), free_tier_end_at = VALUES(free_tier_end_at), discount_code_id = VALUES(discount_code_id)`
	// run
	logf(sqlstr, tu.ID, tu.UserID, tu.AwsCustomerIdentifier, tu.AwsCustomerEntitlement, tu.StripeCustomerIdentifier, tu.StripeCustomerEntitlement, tu.StripeSubscriptionIdentifier, tu.StripePaymentMethodIdentifier, tu.FreeTierEndAt, tu.DiscountCodeID)
	if _, err := db.Exec(sqlstr, tu.ID, tu.UserID, tu.AwsCustomerIdentifier, tu.AwsCustomerEntitlement, tu.StripeCustomerIdentifier, tu.StripeCustomerEntitlement, tu.StripeSubscriptionIdentifier, tu.StripePaymentMethodIdentifier, tu.FreeTierEndAt, tu.DiscountCodeID); err != nil {
		return logerror(err)
	}
	// set exists
	tu._exists = true
	return nil
}

// Delete deletes the [TagbotUser] from the database.
func (tu *TagbotUser) Delete(db DB) error {
	switch {
	case !tu._exists: // doesn't exist
		return nil
	case tu._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM trackit.tagbot_user ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, tu.ID)
	if _, err := db.Exec(sqlstr, tu.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	tu._deleted = true
	return nil
}

// TagbotUserByDiscountCodeID retrieves a row from 'trackit.tagbot_user' as a [TagbotUser].
//
// Generated from index 'foreign_discount_code'.
func TagbotUserByDiscountCodeID(db DB, discountCodeID sql.NullInt64) ([]*TagbotUser, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, user_id, aws_customer_identifier, aws_customer_entitlement, stripe_customer_identifier, stripe_customer_entitlement, stripe_subscription_identifier, stripe_payment_method_identifier, free_tier_end_at, discount_code_id ` +
		`FROM trackit.tagbot_user ` +
		`WHERE discount_code_id = ?`
	// run
	logf(sqlstr, discountCodeID)
	rows, err := db.Query(sqlstr, discountCodeID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*TagbotUser
	for rows.Next() {
		tu := TagbotUser{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&tu.ID, &tu.UserID, &tu.AwsCustomerIdentifier, &tu.AwsCustomerEntitlement, &tu.StripeCustomerIdentifier, &tu.StripeCustomerEntitlement, &tu.StripeSubscriptionIdentifier, &tu.StripePaymentMethodIdentifier, &tu.FreeTierEndAt, &tu.DiscountCodeID); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &tu)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// TagbotUserByID retrieves a row from 'trackit.tagbot_user' as a [TagbotUser].
//
// Generated from index 'tagbot_user_id_pkey'.
func TagbotUserByID(db DB, id int) (*TagbotUser, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, user_id, aws_customer_identifier, aws_customer_entitlement, stripe_customer_identifier, stripe_customer_entitlement, stripe_subscription_identifier, stripe_payment_method_identifier, free_tier_end_at, discount_code_id ` +
		`FROM trackit.tagbot_user ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	tu := TagbotUser{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, id).Scan(&tu.ID, &tu.UserID, &tu.AwsCustomerIdentifier, &tu.AwsCustomerEntitlement, &tu.StripeCustomerIdentifier, &tu.StripeCustomerEntitlement, &tu.StripeSubscriptionIdentifier, &tu.StripePaymentMethodIdentifier, &tu.FreeTierEndAt, &tu.DiscountCodeID); err != nil {
		return nil, logerror(err)
	}
	return &tu, nil
}

// TagbotUserByUserID retrieves a row from 'trackit.tagbot_user' as a [TagbotUser].
//
// Generated from index 'user_id'.
func TagbotUserByUserID(db DB, userID int) (*TagbotUser, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, user_id, aws_customer_identifier, aws_customer_entitlement, stripe_customer_identifier, stripe_customer_entitlement, stripe_subscription_identifier, stripe_payment_method_identifier, free_tier_end_at, discount_code_id ` +
		`FROM trackit.tagbot_user ` +
		`WHERE user_id = ?`
	// run
	logf(sqlstr, userID)
	tu := TagbotUser{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, userID).Scan(&tu.ID, &tu.UserID, &tu.AwsCustomerIdentifier, &tu.AwsCustomerEntitlement, &tu.StripeCustomerIdentifier, &tu.StripeCustomerEntitlement, &tu.StripeSubscriptionIdentifier, &tu.StripePaymentMethodIdentifier, &tu.FreeTierEndAt, &tu.DiscountCodeID); err != nil {
		return nil, logerror(err)
	}
	return &tu, nil
}

// User returns the User associated with the [TagbotUser]'s (UserID).
//
// Generated from foreign key 'tagbot_user_ibfk_1'.
func (tu *TagbotUser) User(db DB) (*User, error) {
	return UserByID(db, tu.UserID)
}

// TagbotDiscountCode returns the TagbotDiscountCode associated with the [TagbotUser]'s (DiscountCodeID).
//
// Generated from foreign key 'tagbot_user_ibfk_2'.
func (tu *TagbotUser) TagbotDiscountCode(db DB) (*TagbotDiscountCode, error) {
	return TagbotDiscountCodeByID(db, int(tu.DiscountCodeID.Int64))
}
