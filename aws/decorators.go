package aws

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// RequireAwsAccount decorates handler to require that an AwsAccount be
// selected using RequiredQueryArgs{AwsAccountIdQueryArg}. The decorator will
// panic if no AwsAccountIdQueryArg query argument is found.
type RequireAwsAccountId struct{}

type routeArgKey uint

const (
	AwsAccountSelection = routeArgKey(iota)
)

func (d RequireAwsAccountId) Decorate(h routes.Handler) routes.Handler {
	h.Func = d.getFunc(h.Func)
	return h
}

func (_ RequireAwsAccountId) getFunc(hf routes.HandlerFunc) routes.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		user, tx, err := getUserAndTransactionFromArguments(a)
		if err != nil {
			l.Error("missing transaction or user for handler with AWS account", err.Error())
			return http.StatusInternalServerError, nil
		}
		aaid := a[routes.AwsAccountIdQueryArg].(int)
		aa, err := GetAwsAccountWithIdFromUser(user, aaid, tx)
		if err != nil {
			return http.StatusNotFound, errors.New("AWS account not found")
		} else {
			a[AwsAccountSelection] = aa
			return hf(w, r, a)
		}

	}
}

func getUserAndTransactionFromArguments(a routes.Arguments) (users.User, *sql.Tx, error) {
	u := a[users.AuthenticatedUser]
	t := a[db.Transaction]
	if ut, ok := u.(users.User); ok {
		if tt, ok := t.(*sql.Tx); ok && t != nil {
			return ut, tt, nil
		} else {
			return users.User{}, nil, errors.New("found no transaction")
		}
	} else {
		return users.User{}, nil, errors.New("found no user")
	}
}
