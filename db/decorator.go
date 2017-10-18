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
			logger.Error("Failed to start SQL transaction.", err.Error())
			return 500, routes.ErrorBody{fmt.Sprintf("Server could not initiate transaction.")}
		}
	}
}
