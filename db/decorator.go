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
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/routes"
)

// Transaction is a decorator which manages a transaction for an HTTP request.
// It will Commit the transaction iff the handler returns something other than
// an error and it did not panic; it Rollbacks otherwise.
type RequestTransaction struct {
	Db *sql.DB
}

type transactionArgumentKey uint

const (
	Transaction = transactionArgumentKey(iota)
)

func (d RequestTransaction) Decorate(h routes.Handler) routes.Handler {
	h.Func = d.getFunc(h.Func)
	return h
}

func (d RequestTransaction) getFunc(hf routes.HandlerFunc) routes.HandlerFunc {
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
					status = http.StatusInternalServerError
					output = errors.New("server suffered irrecoverable error")
					logger.Error("Route handler panicked.", rec)
				} else if _, ok := output.(error); ok {
					transaction.Rollback()
				} else {
					transaction.Commit()
				}
			}()
			return hf(w, r, a)
		} else {
			logger.Error("Failed to start SQL transaction.", err.Error())
			return http.StatusInternalServerError, errors.New("server could not initiate transaction")
		}
	}
}
