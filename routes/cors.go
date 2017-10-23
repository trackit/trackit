package routes

import (
	"fmt"
	"net/http"
	"strings"
)

// WithCors is a route decorator which adds some CORS headers to all responses
// in a static way. This is not in my knowledge correct behavior but it will do
// for now.
type WithCors struct {
	AllowOrigin      []string
	AllowHeaders     []string
	AllowMethods     []string
	AllowCredentials bool
}

func (d WithCors) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		w.Header()["Access-Control-Allow-Origin"] = []string{strings.Join(d.AllowOrigin, ",")}
		w.Header()["Access-Control-Allow-Methods"] = []string{strings.Join(d.AllowMethods, ",")}
		w.Header()["Access-Control-Allow-Headers"] = []string{strings.Join(d.AllowHeaders, ",")}
		w.Header()["Access-Control-Allow-Credentials"] = []string{fmt.Sprintf("%v", d.AllowCredentials)}
		return h(w, r, a)
	}
}
