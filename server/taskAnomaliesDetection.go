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
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/anomaliesDetection"
	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
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
	var dbaa *models.AwsAccount
	var lastUpdate time.Time
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
	} else if dbaa, err = models.AwsAccountByID(tx, aaId); err != nil {
	} else if aa = aws.AwsAccountFromDbAwsAccount(*dbaa); err != nil {
	} else if lastUpdate, err = anomalies.RunAnomaliesDetection(aa, dbaa.LastAnomaliesUpdate, ctx); err == nil {
		err = registerAnomaliesUpdate(tx, lastUpdate, aa.Id)
	}
	if err != nil && !elastic.IsNotFound(err) {
		if strings.HasPrefix(err.Error(), "Data not available yet") {
			logger.Warning("Failed to detect anomalies.", map[string]interface{}{
				"awsAccountId": aaId,
				"error":        err.Error(),
			})
			return
		}
		logger.Error("Failed to detect anomalies.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	return
}

func registerAnomaliesUpdate(tx *sql.Tx, lastUpdate time.Time, aaId int) error {
	dbaa, err := models.AwsAccountByID(tx, aaId)
	if err != nil {
		return err
	}
	dbaa.LastAnomaliesUpdate = lastUpdate
	return dbaa.Update(tx)
}
