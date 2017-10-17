package routes

import (
	"net/http"
)

// RequireMethod is a route decorator to produce a 405 'Method Not Allowed'
// error if a request uses a method not present in its backing slice.
type RequireMethod []string

// methodNotAllowed is the error body sent to the client in case the decorater
// rejects their request.
var methodNotAllowed = ErrorBody{"Method not allowed for requested URL."}

func (d RequireMethod) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		for _, m := range d {
			if r.Method == m {
				return h(w, r, a)
			}
		}
		return 405, methodNotAllowed
	}
}
