package models

import "time"

// GetUnusedAccounts returns all user not seen for the duration passed in arguments
func GetUnusedAccounts(db XODB, unseenFor time.Duration) ([]*User, error) {
	var err error

	limitTime := time.Now().Add(-unseenFor)

	// sql query
	const sqlstr = `SELECT ` +
		`id, created, modified, email, auth, next_external, parent_user_id, aws_customer_identifier, aws_customer_entitlement, next_update_entitlement, anomalies_filters, last_seen ` +
		`FROM trackit.user ` +
		`WHERE last_seen < ?`

	// run query
	XOLog(sqlstr, limitTime.Format(time.RFC3339))
	q, err := db.Query(sqlstr, limitTime.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*User{}
	for q.Next() {
		u := User{
			_exists: true,
		}

		// scan
		err = q.Scan(&u.ID, &u.Created, &u.Modified, &u.Email, &u.Auth, &u.NextExternal, &u.ParentUserID, &u.AwsCustomerIdentifier, &u.AwsCustomerEntitlement, &u.NextUpdateEntitlement, &u.AnomaliesFilters, &u.LastSeen)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}
