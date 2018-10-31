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
	"strconv"

	"github.com/trackit/jsonlog"
	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/trackit/trackit-server/config"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
)

// taskProcessAccount processes an AwsAccount to retrieve data from the AWS api.
func taskTest(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'test'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("test requires an integer argument")
	} else if aaId, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else {
		return test(ctx, aaId)
	}
}

// ingestDataForAccount ingests the AWS api data for an AwsAccount.
func test(ctx context.Context, aaId int) (err error) {
	var tx *sql.Tx
	var aa taws.AwsAccount
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
	} else if aa, err = taws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else {
		logger.Info("ACCOUNT", aa)
		creds, err := taws.GetTemporaryCredentials(aa, "test list account parent")
		if err != nil {
			logger.Error("Error when getting temporary credentials", err.Error())
			return err
		}
		defaultSession := session.Must(session.NewSession(&aws.Config{
			Credentials: creds,
			Region:      aws.String(config.AwsRegion),
		}))
		orga := organizations.New(defaultSession)
		parentId := strings.Split(aa.RoleArn, ":")[4]
		logger.Info("OUI", parentId)
		res, err := orga.ListAccounts(nil)
		if err != nil {
			logger.Error("CA MARCHE PAS", err)
			logger.Error("CA MARCHE PAS2", err.Error())
			return err
		}
		logger.Info("TEST", res)
	}
	if err != nil {
		logger.Error("Failed to test.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	return
}