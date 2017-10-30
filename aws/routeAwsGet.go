package aws

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

// getAwsAccount is a route handler which returns the caller's list of
// AwsAccounts.
func getAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	awsAccounts, err := GetAwsAccountsFromUser(u, tx)
	if err == nil {
		return 200, awsAccounts
	} else {
		l.Error("Failed to get user's AWS accounts.", err.Error())
		return 500, errors.New("Failed to retrieve AWS accounts.")
	}
}
