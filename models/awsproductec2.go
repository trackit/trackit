// Package models contains the types for schema 'trackit'.
package models

// AwsProductEc2Purge purges the AwsProductEc2 table from the database.
func AwsProductEc2Purge(db XODB) error {
	// sql query
	const sqlstr = `DELETE FROM trackit.aws_product_ec2`

	// run query
	XOLog(sqlstr)
	_, err := db.Exec(sqlstr)
	return err
}
