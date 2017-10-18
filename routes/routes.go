package routes

import (
	"encoding/json"
	"net/http"
)

var RegisteredHandlers = make([]RegisteredHandler, 0, 0x40)

type RegisteredHandler struct {
	Pattern string
	Handler IntermediateHandler
}

type Arguments map[interface{}]interface{}

type Handler func(*http.Request, Arguments) (int, interface{})

type IntermediateHandler func(http.ResponseWriter, *http.Request, Arguments) (int, interface{})

type Decorator interface {
	Decorate(IntermediateHandler) IntermediateHandler
}

type ErrorBody struct {
	Error string `json:"error"`
}

func (h IntermediateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arguments := make(Arguments)
	status, output := h(w, r, arguments)
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(output)
}

func ApplyDecorators(h IntermediateHandler, ds ...Decorator) IntermediateHandler {
	l := len(ds) - 1
	for i := range ds {
		h = ds[l-i].Decorate(h)
	}
	return h
}

func Register(pattern string, handler Handler, decorators ...Decorator) {
	stage := baseIntermediate(handler)
	stage = ApplyDecorators(stage, decorators...)
	RegisteredHandlers = append(RegisteredHandlers, RegisteredHandler{
		pattern,
		stage,
	})
}

func baseIntermediate(handler Handler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		return handler(r, a)
	}
}
