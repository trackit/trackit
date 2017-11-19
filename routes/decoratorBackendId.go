package routes

import (
	"net/http"
)

// BackendId is a decorator which adds a backend ID to incoming requests. It
// adds it to tthe response in the `X-Backend-ID` HTTP header.
type BackendId struct {
	BackendId string
}

func (d BackendId) Decorate(h Handler) Handler {
	h.Func = d.getBackendIdFunc(h.Func)
	return h
}

func (d BackendId) getBackendIdFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		w.Header()["X-Backend-ID"] = []string{d.BackendId}
		return hf(w, r, a)
	}
}
