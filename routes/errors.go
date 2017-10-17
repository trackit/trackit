package routes

import (
	"net/http"
)

type WithErrorBody struct{}

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
