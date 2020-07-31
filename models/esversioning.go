// Package models contains the types for schema 'trackit'.
package models

// OutdatedEsMappings retrieves a row from 'trackit.es_versioning' as a EsVersioning which is outdated
func OutdatedEsMappings(db XODB, templateName string, latestVersion int) ([]*EsVersioning, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, current_version, template_name, index_name ` +
		`FROM trackit.es_versioning ` +
		`WHERE current_version < ? AND template_name = ?`

	// run query
	XOLog(sqlstr, latestVersion, templateName)
	q, err := db.Query(sqlstr, latestVersion, templateName)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*EsVersioning{}
	for q.Next() {
		ev := EsVersioning{
			_exists: true,
		}

		// scan
		err = q.Scan(&ev.ID, &ev.CurrentVersion, &ev.TemplateName, &ev.IndexName)
		if err != nil {
			return nil, err
		}

		res = append(res, &ev)
	}

	return res, nil
}
