package routes

import (
	"net/http"
)

// CorsAllowOrigin is a route decorator which adds the
// `Access-Control-Allow-Origin` header to all response of the route it
// decorates.
type WithCorsAllowOrigin []string

func (d WithCorsAllowOrigin) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		w.Header()["Access-Control-Allow-Origin"] = d
		return h(w, r, a)
	}
}
