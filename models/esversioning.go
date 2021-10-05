// Package models contains the types for schema 'trackit'.
package models

// OutdatedEsMappings retrieves a row from 'trackit.es_versioning' as a EsVersioning which is outdated
func OutdatedEsMappings(db DB, templateName string, latestVersion int) ([]*EsVersioning, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, current_version, template_name, index_name ` +
		`FROM trackit.es_versioning ` +
		`WHERE current_version < ? AND template_name = ?`

	// run query
	Logf(sqlstr, latestVersion, templateName)
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

// EsVersioningByIndexName retrieves a row from 'trackit.es_versioning' as a EsVersioning (no ID required)
func EsVersioningByIndexName(db DB, indexName string) (*EsVersioning, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, current_version, template_name, index_name ` +
		`FROM trackit.es_versioning ` +
		`WHERE index_name = ?`
	// run
	logf(sqlstr, indexName)
	ev := EsVersioning{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, indexName).Scan(&ev.ID, &ev.CurrentVersion, &ev.TemplateName, &ev.IndexName); err != nil {
		return nil, logerror(err)
	}
	return &ev, nil
}
