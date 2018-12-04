//   Copyright 2018 MSolution.IO
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
	"strconv"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/anomaliesDetection"
	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
)

// taskAnomaliesDetection processes an AwsAccount to email
// the user if anomalies are detected.
func taskAnomaliesDetection(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'anomalies-detection'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("taskAnomaliesDetection requires an integer argument")
	} else if aaId, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else {
		return processAnomaliesForAccount(ctx, aaId)
	}
}

func processAnomaliesForAccount(ctx context.Context, aaId int) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
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
	} else if err = anomalies.RunAnomaliesDetection(aa, ctx, tx); err != nil {
	}
	if err != nil {
		logger.Error("Failed to detect anomalies.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	return
}
