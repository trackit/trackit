// Package models contains the types for schema 'trackit'.
package models

// AwsBillRepositoriesWithDueUpdate returns the set of bill repositories with a
// due update.
func AwsBillRepositoriesWithDueUpdate(db XODB) ([]*AwsBillRepository, error) {
	var err error
	const sqlstr = `SELECT ` +
		`id, aws_account_id, bucket, prefix, last_imported_manifest, next_update ` +
		`FROM trackit.aws_bill_repository ` +
		`WHERE next_update <= NOW()`
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	res := []*AwsBillRepository{}
	for q.Next() {
		abr := AwsBillRepository{
			_exists: true,
		}
		err = q.Scan(&abr.ID, &abr.AwsAccountID, &abr.Bucket, &abr.Prefix, &abr.LastImportedManifest, &abr.NextUpdate)
		if err != nil {
			return nil, err
		}
		res = append(res, &abr)
	}
	return res, nil
}

// UpdateUnsafe updates the BillRepository but doesn't do XO's usual checks.
func (abr *AwsBillRepository) UpdateUnsafe(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE trackit.aws_bill_repository SET ` +
		`aws_account_id = ?, bucket = ?, prefix = ?, last_imported_manifest = ?, next_update = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, abr.AwsAccountID, abr.Bucket, abr.Prefix, abr.LastImportedManifest, abr.NextUpdate, abr.ID)
	_, err = db.Exec(sqlstr, abr.AwsAccountID, abr.Bucket, abr.Prefix, abr.LastImportedManifest, abr.NextUpdate, abr.ID)
	return err
}
