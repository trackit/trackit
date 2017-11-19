package routes

import (
	"net/http"
)

// Handler is the type a route Handler must have.
type SimpleHandler func(*http.Request, Arguments) (int, interface{})

func H(h SimpleHandler) Handler {
	return Handler{
		Func: func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
			return h(r, a)
		},
	}
}
