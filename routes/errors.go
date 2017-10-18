package routes

import (
	"net/http"
)

// WithErrorBody is a decorator for an HTTP handler. If that handler returns an
// error, it will wrap it in an ErrorBody structure so that it can correctly be
// returned to the user as JSON.
type WithErrorBody struct{}

type ErrorBody struct {
	Error string `json:"error"`
}

func (d WithErrorBody) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		status, output := h(w, r, a)
		if err, ok := output.(error); ok {
			return status, ErrorBody{err.Error()}
		} else {
			return status, output
		}
	}
}
