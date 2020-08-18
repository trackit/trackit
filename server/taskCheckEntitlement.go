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
	"errors"
	"flag"
	"strconv"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/entitlement"
)

// taskCheckEntitlement checks the user Entitlement for AWS Marketplace users
func taskCheckEntitlement(ctx context.Context) (err error) {
	var tx *sql.Tx
	args := flag.Args()
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
	logger.Debug("Running task 'checkuserentitlement'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("taskCheckEntitlement requires one integer argument")
	} else if userId, err := strconv.Atoi(args[0]); err != nil {
		tx, err = db.Db.BeginTx(ctx, nil)
		if err != nil {
			logger.Error("Could not connect to database.", err.Error())
		}
		return entitlement.CheckUserEntitlements(ctx, tx, userId)
	} else {
	}
	return nil
}
