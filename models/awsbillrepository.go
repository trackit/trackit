// Package models contains the types for schema 'trackit'.
package models

// AwsBillRepositoriesWithDueUpdate returns the set of bill repositories with a
// due update.
func AwsBillRepositoriesWithDueUpdate(db XODB) ([]*AwsBillRepository, error) {
	var err error
	const sqlstr = `SELECT ` +
		`id, aws_account_id, bucket, prefix, last_imported_period, next_update ` +
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
		err = q.Scan(&abr.ID, &abr.AwsAccountID, &abr.Bucket, &abr.Prefix, &abr.LastImportedPeriod, &abr.NextUpdate)
		if err != nil {
			return nil, err
		}
		res = append(res, &abr)
	}
	return res, nil
}
