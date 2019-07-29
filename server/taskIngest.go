//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
)

// taskIngest ingests billing data for a given BillRepository and AwsAccount.
func taskIngest(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'ingest'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 2 {
		return errors.New("taskIngest requires two integer arguments")
	} else if aa, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else if br, err := strconv.Atoi(args[1]); err != nil {
		return err
	} else {
		return ingestBillingDataForBillRepository(ctx, aa, br)
	}
}

// ingestBillingDataForBillRepository ingests the billing data for a
// BillRepository.
func ingestBillingDataForBillRepository(ctx context.Context, aaId, brId int) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	var br s3.BillRepository
	var updateId int64
	var latestManifest time.Time
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if br, err = s3.GetBillRepositoryForAwsAccountById(aa, brId, tx); err != nil {
	} else if updateId, err = registerUpdate(db.Db, br); err != nil {
	} else if latestManifest, err = s3.UpdateReport(ctx, aa, br); err != nil {
		if billError, castok := err.(awserr.Error); castok {
			br.Error = billError.Message()
			s3.UpdateBillRepositoryWithoutContext(br, db.Db)
		}
	} else {
		br.Error = ""
		err = updateBillRepositoryForNextUpdate(ctx, tx, br, latestManifest)
	}
	if err != nil {
		logger.Error("Failed to ingest billing data.", map[string]interface{}{
			"awsAccountId":     aaId,
			"billRepositoryId": brId,
			"error":            err.Error(),
		})
	}
	updateCompletion(ctx, aaId, brId, db.Db, updateId, err)
	updateSubAccounts(ctx, aa)
	var affectedRoutes = []string{
		"/costs",
		"/costs/diff",
		"/costs/tags/keys",
		"/costs/tags/values",
		"/s3/costs",
	}
	_ = cache.RemoveMatchingCache(affectedRoutes, []string{aa.AwsIdentity}, logger)
	return
}

func updateSubAccounts(ctx context.Context, aa aws.AwsAccount) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var tx *sql.Tx
	var err error
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
		logger.Error("Failed to get DB Tx", err.Error())
	} else if err = aws.PutSubAccounts(ctx, aa, tx); err != nil {
		logger.Warning("Failed to update sub accounts.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	} else {
		logger.Info("Sub accounts updated.", map[string]interface{}{
			"awsAccountId": aa.Id,
		})
	}
}

func registerUpdate(db *sql.DB, br s3.BillRepository) (int64, error) {
	const sqlstr = `INSERT INTO aws_bill_update_job(
		aws_bill_repository_id,
		worker_id,
		error
	) VALUES (?, ?, "")`
	res, err := db.Exec(sqlstr, br.Id, backendId)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateCompletion(ctx context.Context, aaId, brId int, db *sql.DB, updateId int64, err error) {
	rErr := registerUpdateCompletion(db, updateId, err)
	if rErr != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to register ingestion completion.", map[string]interface{}{
			"awsAccountId":     aaId,
			"billRepositoryId": brId,
			"error":            rErr.Error(),
			"updateId":         updateId,
		})
	}
}

func registerUpdateCompletion(db *sql.DB, updateId int64, err error) error {
	const sqlstr = `UPDATE aws_bill_update_job SET
		completed=?,
		error=?
	WHERE id=?`
	var errorValue string
	var now = time.Now()
	if err != nil {
		errorValue = err.Error()
	}
	_, err = db.Exec(sqlstr, now, errorValue, updateId)
	return err
}

const (
	UpdateIntervalMinutes = 6 * 60
	UpdateIntervalWindow  = 2 * 60
)

// updateBillRepositoryForNextUpdate plans the next update for a
// BillRepository.
func updateBillRepositoryForNextUpdate(ctx context.Context, tx *sql.Tx, br s3.BillRepository, latestManifest time.Time) error {
	if latestManifest.After(br.LastImportedManifest) {
		br.LastImportedManifest = latestManifest
	}
	updateDeltaMinutes := time.Duration(UpdateIntervalMinutes-UpdateIntervalWindow/2+rand.Int63n(UpdateIntervalWindow)) * time.Minute
	br.NextUpdate = time.Now().Add(updateDeltaMinutes)
	return s3.UpdateBillRepository(br, tx)
}
