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

package db

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/routes"
)

// WithTransaction is a decorator which manages a transaction for an HTTP
// request. It will Commit the transaction iff the handler returns something
// other than an error and it did not panic; it Rollbacks otherwise.
type WithTransaction struct {
	Db *sql.DB
}

type withTransactionArgumentKey uint

const (
	Transaction = withTransactionArgumentKey(iota)
)

func (d WithTransaction) Decorate(h routes.IntermediateHandler) routes.IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (status int, output interface{}) {
		ctx := r.Context()
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		transaction, err := d.Db.BeginTx(ctx, nil)
		if err == nil {
			a[Transaction] = transaction
			defer func() {
				rec := recover()
				if rec != nil {
					transaction.Rollback()
					status = 500
					output = routes.ErrorBody{fmt.Sprintf("Server suffered irrecoverable error.")}
					logger.Error("Route handler panicked.", rec)
				} else if _, ok := output.(error); ok {
					transaction.Rollback()
				} else {
					transaction.Commit()
				}
			}()
			return h(w, r, a)
		} else {
			logger.Error("Failed to start Sql transaction.", err.Error())
			return 500, routes.ErrorBody{fmt.Sprintf("Server could not initiate transaction.")}
		}
	}
}
