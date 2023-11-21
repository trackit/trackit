package models

// Code generated by xo. DO NOT EDIT.

import (
	"time"
)

// AwsAccountReportsJob represents a row from 'trackit.aws_account_reports_job'.
type AwsAccountReportsJob struct {
	ID                          int       `json:"id"`                          // id
	AwsAccountID                int       `json:"aws_account_id"`              // aws_account_id
	Completed                   time.Time `json:"completed"`                   // completed
	WorkerID                    string    `json:"worker_id"`                   // worker_id
	JobError                    string    `json:"jobError"`                    // jobError
	SpreadsheetError            string    `json:"spreadsheetError"`            // spreadsheetError
	CostDiffError               string    `json:"costDiffError"`               // costDiffError
	Ec2usageReportError         string    `json:"ec2UsageReportError"`         // ec2UsageReportError
	RdsUsageReportError         string    `json:"rdsUsageReportError"`         // rdsUsageReportError
	EsUsageReportError          string    `json:"esUsageReportError"`          // esUsageReportError
	ElasticacheUsageReportError string    `json:"elasticacheUsageReportError"` // elasticacheUsageReportError
	LambdaUsageReportError      string    `json:"lambdaUsageReportError"`      // lambdaUsageReportError
	RiEc2reportError            string    `json:"riEc2ReportError"`            // riEc2ReportError
	RiRdsReportError            string    `json:"riRdsReportError"`            // riRdsReportError
	OdToRiEc2reportError        string    `json:"odToRiEc2ReportError"`        // odToRiEc2ReportError
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [AwsAccountReportsJob] exists in the database.
func (aarj *AwsAccountReportsJob) Exists() bool {
	return aarj._exists
}

// Deleted returns true when the [AwsAccountReportsJob] has been marked for deletion
// from the database.
func (aarj *AwsAccountReportsJob) Deleted() bool {
	return aarj._deleted
}

// Insert inserts the [AwsAccountReportsJob] to the database.
func (aarj *AwsAccountReportsJob) Insert(db DB) error {
	switch {
	case aarj._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case aarj._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO trackit.aws_account_reports_job (` +
		`aws_account_id, completed, worker_id, jobError, spreadsheetError, costDiffError, ec2UsageReportError, rdsUsageReportError, esUsageReportError, elasticacheUsageReportError, lambdaUsageReportError, riEc2ReportError, riRdsReportError, odToRiEc2ReportError` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, aarj.AwsAccountID, aarj.Completed, aarj.WorkerID, aarj.JobError, aarj.SpreadsheetError, aarj.CostDiffError, aarj.Ec2usageReportError, aarj.RdsUsageReportError, aarj.EsUsageReportError, aarj.ElasticacheUsageReportError, aarj.LambdaUsageReportError, aarj.RiEc2reportError, aarj.RiRdsReportError, aarj.OdToRiEc2reportError)
	res, err := db.Exec(sqlstr, aarj.AwsAccountID, aarj.Completed, aarj.WorkerID, aarj.JobError, aarj.SpreadsheetError, aarj.CostDiffError, aarj.Ec2usageReportError, aarj.RdsUsageReportError, aarj.EsUsageReportError, aarj.ElasticacheUsageReportError, aarj.LambdaUsageReportError, aarj.RiEc2reportError, aarj.RiRdsReportError, aarj.OdToRiEc2reportError)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	aarj.ID = int(id)
	// set exists
	aarj._exists = true
	return nil
}

// Update updates a [AwsAccountReportsJob] in the database.
func (aarj *AwsAccountReportsJob) Update(db DB) error {
	switch {
	case !aarj._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case aarj._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE trackit.aws_account_reports_job SET ` +
		`aws_account_id = ?, completed = ?, worker_id = ?, jobError = ?, spreadsheetError = ?, costDiffError = ?, ec2UsageReportError = ?, rdsUsageReportError = ?, esUsageReportError = ?, elasticacheUsageReportError = ?, lambdaUsageReportError = ?, riEc2ReportError = ?, riRdsReportError = ?, odToRiEc2ReportError = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, aarj.AwsAccountID, aarj.Completed, aarj.WorkerID, aarj.JobError, aarj.SpreadsheetError, aarj.CostDiffError, aarj.Ec2usageReportError, aarj.RdsUsageReportError, aarj.EsUsageReportError, aarj.ElasticacheUsageReportError, aarj.LambdaUsageReportError, aarj.RiEc2reportError, aarj.RiRdsReportError, aarj.OdToRiEc2reportError, aarj.ID)
	if _, err := db.Exec(sqlstr, aarj.AwsAccountID, aarj.Completed, aarj.WorkerID, aarj.JobError, aarj.SpreadsheetError, aarj.CostDiffError, aarj.Ec2usageReportError, aarj.RdsUsageReportError, aarj.EsUsageReportError, aarj.ElasticacheUsageReportError, aarj.LambdaUsageReportError, aarj.RiEc2reportError, aarj.RiRdsReportError, aarj.OdToRiEc2reportError, aarj.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [AwsAccountReportsJob] to the database.
func (aarj *AwsAccountReportsJob) Save(db DB) error {
	if aarj.Exists() {
		return aarj.Update(db)
	}
	return aarj.Insert(db)
}

// Upsert performs an upsert for [AwsAccountReportsJob].
func (aarj *AwsAccountReportsJob) Upsert(db DB) error {
	switch {
	case aarj._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO trackit.aws_account_reports_job (` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, costDiffError, ec2UsageReportError, rdsUsageReportError, esUsageReportError, elasticacheUsageReportError, lambdaUsageReportError, riEc2ReportError, riRdsReportError, odToRiEc2ReportError` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`aws_account_id = VALUES(aws_account_id), completed = VALUES(completed), worker_id = VALUES(worker_id), jobError = VALUES(jobError), spreadsheetError = VALUES(spreadsheetError), costDiffError = VALUES(costDiffError), ec2UsageReportError = VALUES(ec2UsageReportError), rdsUsageReportError = VALUES(rdsUsageReportError), esUsageReportError = VALUES(esUsageReportError), elasticacheUsageReportError = VALUES(elasticacheUsageReportError), lambdaUsageReportError = VALUES(lambdaUsageReportError), riEc2ReportError = VALUES(riEc2ReportError), riRdsReportError = VALUES(riRdsReportError), odToRiEc2ReportError = VALUES(odToRiEc2ReportError)`
	// run
	logf(sqlstr, aarj.ID, aarj.AwsAccountID, aarj.Completed, aarj.WorkerID, aarj.JobError, aarj.SpreadsheetError, aarj.CostDiffError, aarj.Ec2usageReportError, aarj.RdsUsageReportError, aarj.EsUsageReportError, aarj.ElasticacheUsageReportError, aarj.LambdaUsageReportError, aarj.RiEc2reportError, aarj.RiRdsReportError, aarj.OdToRiEc2reportError)
	if _, err := db.Exec(sqlstr, aarj.ID, aarj.AwsAccountID, aarj.Completed, aarj.WorkerID, aarj.JobError, aarj.SpreadsheetError, aarj.CostDiffError, aarj.Ec2usageReportError, aarj.RdsUsageReportError, aarj.EsUsageReportError, aarj.ElasticacheUsageReportError, aarj.LambdaUsageReportError, aarj.RiEc2reportError, aarj.RiRdsReportError, aarj.OdToRiEc2reportError); err != nil {
		return logerror(err)
	}
	// set exists
	aarj._exists = true
	return nil
}

// Delete deletes the [AwsAccountReportsJob] from the database.
func (aarj *AwsAccountReportsJob) Delete(db DB) error {
	switch {
	case !aarj._exists: // doesn't exist
		return nil
	case aarj._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM trackit.aws_account_reports_job ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, aarj.ID)
	if _, err := db.Exec(sqlstr, aarj.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	aarj._deleted = true
	return nil
}

// AwsAccountReportsJobByID retrieves a row from 'trackit.aws_account_reports_job' as a [AwsAccountReportsJob].
//
// Generated from index 'aws_account_reports_job_id_pkey'.
func AwsAccountReportsJobByID(db DB, id int) (*AwsAccountReportsJob, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, costDiffError, ec2UsageReportError, rdsUsageReportError, esUsageReportError, elasticacheUsageReportError, lambdaUsageReportError, riEc2ReportError, riRdsReportError, odToRiEc2ReportError ` +
		`FROM trackit.aws_account_reports_job ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	aarj := AwsAccountReportsJob{
		_exists: true,
	}
	if err := db.QueryRow(sqlstr, id).Scan(&aarj.ID, &aarj.AwsAccountID, &aarj.Completed, &aarj.WorkerID, &aarj.JobError, &aarj.SpreadsheetError, &aarj.CostDiffError, &aarj.Ec2usageReportError, &aarj.RdsUsageReportError, &aarj.EsUsageReportError, &aarj.ElasticacheUsageReportError, &aarj.LambdaUsageReportError, &aarj.RiEc2reportError, &aarj.RiRdsReportError, &aarj.OdToRiEc2reportError); err != nil {
		return nil, logerror(err)
	}
	return &aarj, nil
}

// AwsAccountReportsJobByAwsAccountID retrieves a row from 'trackit.aws_account_reports_job' as a [AwsAccountReportsJob].
//
// Generated from index 'foreign_aws_account'.
func AwsAccountReportsJobByAwsAccountID(db DB, awsAccountID int) ([]*AwsAccountReportsJob, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, aws_account_id, completed, worker_id, jobError, spreadsheetError, costDiffError, ec2UsageReportError, rdsUsageReportError, esUsageReportError, elasticacheUsageReportError, lambdaUsageReportError, riEc2ReportError, riRdsReportError, odToRiEc2ReportError ` +
		`FROM trackit.aws_account_reports_job ` +
		`WHERE aws_account_id = ?`
	// run
	logf(sqlstr, awsAccountID)
	rows, err := db.Query(sqlstr, awsAccountID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*AwsAccountReportsJob
	for rows.Next() {
		aarj := AwsAccountReportsJob{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&aarj.ID, &aarj.AwsAccountID, &aarj.Completed, &aarj.WorkerID, &aarj.JobError, &aarj.SpreadsheetError, &aarj.CostDiffError, &aarj.Ec2usageReportError, &aarj.RdsUsageReportError, &aarj.EsUsageReportError, &aarj.ElasticacheUsageReportError, &aarj.LambdaUsageReportError, &aarj.RiEc2reportError, &aarj.RiRdsReportError, &aarj.OdToRiEc2reportError); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &aarj)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// AwsAccount returns the AwsAccount associated with the [AwsAccountReportsJob]'s (AwsAccountID).
//
// Generated from foreign key 'aws_account_reports_job_ibfk_1'.
func (aarj *AwsAccountReportsJob) AwsAccount(db DB) (*AwsAccount, error) {
	return AwsAccountByID(db, aarj.AwsAccountID)
}
