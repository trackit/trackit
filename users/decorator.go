package users

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
)

type WithAuthenticatedUser struct{}

type withAuthenticatedUserArgumentKey uint

const (
	AuthenticatedUser = withAuthenticatedUserArgumentKey(iota)
)

func (d WithAuthenticatedUser) Decorate(h routes.IntermediateHandler) routes.IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		auth := r.Header["Authorization"]
		tx := a[db.Transaction].(*sql.Tx)
		if auth != nil && len(auth) == 1 {
			tokenString := auth[0]
			if user, err := testToken(tx, tokenString); err == nil {
				a[AuthenticatedUser] = user
				return h(w, r, a)
			}
		}
		return 401, errors.New("Invalid or missing token.")
	}
}
