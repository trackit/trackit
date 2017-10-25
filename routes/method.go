package routes

import (
	"fmt"
	"net/http"
	"strings"
)

// RequireMethod is a route decorator to produce a 405 'Method Not Allowed'
// error if a request uses a method not present in its backing slice. It also
// responds correctly to OPTIONS requests.
type RequireMethod []string

// methodNotAllowed is the error body sent to the client in case the decorater
// rejects their request.
var methodNotAllowed = ErrorBody{"Method not allowed for requested URL."}

func (d RequireMethod) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		if r.Method == "OPTIONS" {
			return d.handleOption(w, r)
		} else {
			return d.handleAny(h, w, r, a)
		}
	}
}

func (d RequireMethod) handleOption(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	if len(d) > 0 {
		w.Header()["Allow"] = []string{fmt.Sprintf("%s,OPTIONS", strings.Join(d, ","))}
	} else {
		w.Header()["Allow"] = []string{"OPTIONS"}
	}
	return 200, nil
}

func (d RequireMethod) handleAny(h IntermediateHandler, w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
	for _, m := range d {
		if r.Method == m {
			return h(w, r, a)
		}
	}
	return 405, methodNotAllowed
}
