package main

import (
	"context"
	"database/sql"

	"github.com/trackit/jsonlog"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/db"
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

// updateBillRepositoriesFromConclusion updates bill repositories in the
// database using the conclusion of an update task.
func updateBillRepositoriesFromConclusion(ctx context.Context, tx *sql.Tx, ruccs []s3.ReportUpdateConclusion) error {
	for _, r := range ruccs {
		if r.Error != nil {
			if billError, castok := r.Error.(awserr.Error); castok {
				r.BillRepository.Error = billError.Message()
				if err := s3.UpdateBillRepository(r.BillRepository, tx); err != nil {
					return err
				}
			}
		} else {
			r.BillRepository.Error = ""
			if err := updateBillRepositoryForNextUpdate(ctx, tx, r.BillRepository, r.LastImportedManifest); err != nil {
				return err
			}
		}
	}
	return nil
}
