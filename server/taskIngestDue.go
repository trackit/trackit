package main

import (
	"context"
	"database/sql"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws/s3"
	"github.com/trackit/trackit-server/db"
)

// taskIngestDue lists all BillRepositories with due updates and updates them.
func taskIngestDue(ctx context.Context) (err error) {
	var tx *sql.Tx
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
				logger.Debug("Rolled back transaction.", nil)
			} else {
				tx.Commit()
				logger.Debug("Commited transaction.", nil)
			}
		}
	}()
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else {
		logger.Debug("Started transaction.", nil)
		conclusion, err := s3.UpdateDueReports(ctx, tx)
		if err == nil {
			err = updateBillRepositoriesFromConclusion(ctx, tx, conclusion)
		}
	}
	return
}
