//   Copyright 2019 MSolution.IO
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

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
)

func taskUpdateAwsIdentity(ctx context.Context) error {
	var tx *sql.Tx
	var err error
	var aas []*models.AwsAccount
	l := jsonlog.LoggerFromContextOrDefault(ctx)
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
	} else if aas, err = models.AwsAccounts(tx); err == nil {
		for _, dbAa := range aas {
			aa := aws.AwsAccountFromDbAwsAccount(*dbAa)
			err = aa.UpdateIdentityAwsAccount(ctx, tx)
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		l.Info("Aws identity updated for all accounts", nil)
	} else {
		l.Error("Error while updating aws identity", err.Error())
	}
	return err
}
