//   Copyright 2021 MSolution.IO
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
	"database/sql"
	"fmt"

	"github.com/trackit/jsonlog"
)

// Finalizes the transaction using tx.Rollback/tx.Commit (depending on whether *errPtr/*txPtr are nil) and properly handles an error from either of those operations. The first three arguments are passed by pointer because when this is used in defer, the arguments are evaluated early and would not be passed properly otherwise
func utilsUsualTxFinalize(txPtr **sql.Tx, errPtr *error, logger *jsonlog.Logger, transactionName string) {
	tx := *txPtr
	if tx != nil {
		if *errPtr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Error(fmt.Sprintf("While handling error, failed to rollback %s transaction", transactionName), map[string]interface{}{
					"error": rollbackErr.Error(),
				})
			} else {
				logger.Debug("Rolled back transaction", nil)
			}
		} else {
			*errPtr = tx.Commit()
			if *errPtr == nil {
				logger.Debug("Commited transaction", nil)
			}
		}
	}
}
