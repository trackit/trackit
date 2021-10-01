package models

// Code generated by xo. DO NOT EDIT.

// AwsPricing represents a row from 'trackit.aws_pricing'.
type AwsPricing struct {
	ID      int    `json:"id"`      // id
	Product string `json:"product"` // product
	Pricing []byte `json:"pricing"` // pricing
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the AwsPricing exists in the database.
func (ap *AwsPricing) Exists() bool {
	return ap._exists
}

// Deleted returns true when the AwsPricing has been marked for deletion from
// the database.
func (ap *AwsPricing) Deleted() bool {
	return ap._deleted
}

// Insert inserts the AwsPricing to the database.
func (ap *AwsPricing) Insert(db DB) error {
	switch {
	case ap._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case ap._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO trackit.aws_pricing (` +
		`product, pricing` +
		`) VALUES (` +
		`?, ?` +
		`)`
	// run
	logf(sqlstr, ap.Product, ap.Pricing)
	res, err := db.Exec(sqlstr, ap.Product, ap.Pricing)
	if err != nil {
		return err
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	} // set primary key
	ap.ID = int(id)
	// set exists
	ap._exists = true
	return nil
}

// Update updates a AwsPricing in the database.
func (ap *AwsPricing) Update(db DB) error {
	switch {
	case !ap._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case ap._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE trackit.aws_pricing SET ` +
		`product = ?, pricing = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, ap.Product, ap.Pricing, ap.ID)
	if _, err := db.Exec(sqlstr, ap.Product, ap.Pricing, ap.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the AwsPricing to the database.
func (ap *AwsPricing) Save(db DB) error {
	if ap.Exists() {
		return ap.Update(db)
	}
	return ap.Insert(db)
}

// Upsert performs an upsert for AwsPricing.
func (ap *AwsPricing) Upsert(db DB) error {
	switch {
	case ap._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO trackit.aws_pricing (` +
		`id, product, pricing` +
		`) VALUES (` +
		`?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`product = VALUES(product), pricing = VALUES(pricing)`
	// run
	logf(sqlstr, ap.ID, ap.Product, ap.Pricing)
	if _, err := db.Exec(sqlstr, ap.ID, ap.Product, ap.Pricing); err != nil {
		return err
	}
	// set exists
	ap._exists = true
	return nil
}

// Delete deletes the AwsPricing from the database.
func (ap *AwsPricing) Delete(db DB) error {
	switch {
	case !ap._exists: // doesn't exist
		return nil
	case ap._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM trackit.aws_pricing ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, ap.ID)
	if _, err := db.Exec(sqlstr, ap.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	ap._deleted = true
	return nil
}

// AwsPricingByID retrieves a row from 'trackit.aws_pricing' as a AwsPricing.
//
// Generated from index 'aws_pricing_id_pkey'.
func AwsPricingByID(db DB, id int) (*AwsPricing, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, product, pricing ` +
		`FROM trackit.aws_pricing ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	ap := AwsPricing{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, id).Scan(&ap.ID, &ap.Product, &ap.Pricing); err != nil {
		return nil, logerror(err)
	}
	return &ap, nil
}

// AwsPricingByProduct retrieves a row from 'trackit.aws_pricing' as a AwsPricing.
//
// Generated from index 'product'.
func AwsPricingByProduct(db DB, product string) (*AwsPricing, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, product, pricing ` +
		`FROM trackit.aws_pricing ` +
		`WHERE product = ?`
	// run
	logf(sqlstr, product)
	ap := AwsPricing{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, product).Scan(&ap.ID, &ap.Product, &ap.Pricing); err != nil {
		return nil, logerror(err)
	}
	return &ap, nil
}
