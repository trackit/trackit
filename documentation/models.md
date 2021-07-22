# Models

In order to interact with the SQL database, why use `xo` which is a command-line tool used to generate Go code based on a database schema.

In order to modify the database architecture, you must add a migration file in the `db/migration` folder. Your migration should be added at the end of the 'db/schema.sql' file. In a development environment, you have to manually apply the migration before you can continue.

Once the database is modified, you can use the `xo` tool to generate go models:

`xo mysql://user:pass@host/dbname -o models`.

Be careful, you should only commit the files linked to your changes.

You can now use the generated code to interact with the database.

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
