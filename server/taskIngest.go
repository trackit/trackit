package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"math/rand"
	"strconv"
	"time"

	"github.com/trackit/trackit2/aws"
	"github.com/trackit/trackit2/aws/s3"
	"github.com/trackit/trackit2/db"
)

func taskIngest(ctx context.Context) error {
	args := flag.Args()
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

func ingestBillingDataForBillRepository(ctx context.Context, aaId, brId int) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	var br s3.BillRepository
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
				println("ROLLBACK!")
			} else {
				tx.Commit()
				println("COMMIT!")
			}
		}
	}()
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if br, err = s3.GetBillRepositoryForAwsAccountById(aa, brId, tx); err != nil {
	} else if latestManifest, err := s3.UpdateReport(ctx, aa, br); err != nil {
	} else {
		err = updateBillRepositoryForNextUpdate(ctx, tx, br, latestManifest)
	}
	if err != nil {
		println(err.Error())
	}
	return
}

const (
	UpdateIntervalMinutes = 6 * 60
	UpdateIntervalWindow  = 2 * 60
)

func updateBillRepositoryForNextUpdate(ctx context.Context, tx *sql.Tx, br s3.BillRepository, latestManifest time.Time) error {
	br.LastImportedManifest = latestManifest
	updateDeltaMinutes := time.Duration(UpdateIntervalMinutes-UpdateIntervalWindow/2+rand.Int63n(UpdateIntervalWindow)) * time.Minute
	br.NextUpdate = time.Now().Add(updateDeltaMinutes)
	return s3.UpdateBillRepository(br, tx)
}
