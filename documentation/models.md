# Models

The `xo` command-line tool can be useful for creating code to interact with the SQL database. It generates Go code based on a database schema.

In order to modify the database architecture, you must add a migration file in the `db/migration` folder. Your migration should be added at the end of the `db/schema.sql` file. In a development environment, you have to manually apply the migration before you can continue.

A command such as this should be used to check that `db/schema.sql` corresponds to the files contained in `db/migration` (it prints the difference, i.e. outputs nothing if `db/schema.sql` is correct):

`diff -u db/schema.sql <(awk 'FNR==1{print ""}1' db/migration/*.sql | tail -n+2)`

Once the database is modified, you can use the `xo` tool to generate go models. A command line such as this should get generated code alike to that which is currently used:

`xo schema --exclude='*.created' --exclude='*.modified' --go-context=disable mysql://user:pass@host/dbname -o models`.

Be careful, you should only commit the files linked to your changes.

You can then use the generated code to interact with the database.

### Examples:

Inserting data:
```go
dbReportGeneration := models.AwsAccountMasterReportsJob{
    AwsAccountID: aa.Id,
    WorkerID:     backendId,
}
err := dbReportGeneration.Insert(db)
```

Selecting and updating data:
```go
dbAccountReports, err := models.AwsAccountMasterReportsJobByID(db, int(updateId))
dbAccountReports.Completed = time.Now().UTC()
err = dbAccountReports.Update(db)
```
