package routes

import (
	"encoding/json"
	"net/http"
)

// RegisteredHandlers is the list of all route handlers that were registered.
// Modules providing route handlers are expected to run Register in order to
// populate that list, and the main package is expected to use this list to
// populate its HTTP server.
var RegisteredHandlers = make([]RegisteredHandler, 0, 0x40)

// RegisteredHandler is a registered handler with the pattern it will serve.
type RegisteredHandler struct {
	Pattern string
	Handler IntermediateHandler
}

// Arguments is a map used by decorators to supply the route handler (or later
// decorators) with additional values.
type Arguments map[interface{}]interface{}

// Handler is the type a route Handler must have.
type Handler func(*http.Request, Arguments) (int, interface{})

// IntermediateHandler represents a decorated handler. Decorators decorate
// IntermediateHandlers and return IntermediateHandlers.
type IntermediateHandler func(http.ResponseWriter, *http.Request, Arguments) (int, interface{})

// Decorator is an interface for any type that can decorate an
// IntermediateHandler.
type Decorator interface {
	Decorate(IntermediateHandler) IntermediateHandler
}

// ServeHTTP is a method on IntermediateHandler that lets it fulfill the
// http.Handler interface.
func (h IntermediateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arguments := make(Arguments)
	status, output := h(w, r, arguments)
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(output)
}

// ApplyDecorators applies the supplied decorators on an IntermediateHandler.
// The first decorator is the outermost one, while the last is the innermost
// one (i.e. the decorators are applied in reverse order for readability).
func ApplyDecorators(h IntermediateHandler, ds ...Decorator) IntermediateHandler {
	l := len(ds) - 1
	for i := range ds {
		h = ds[l-i].Decorate(h)
	}
	return h
}

// Register decorates a Handler and registers it into the RegisteredHandlers
// list.
func Register(pattern string, handler Handler, decorators ...Decorator) {
	stage := baseIntermediate(handler)
	stage = ApplyDecorators(stage, decorators...)
	RegisteredHandlers = append(RegisteredHandlers, RegisteredHandler{
		pattern,
		stage,
	})
}

// baseIntermediate wraps a Handler into an IntermediateHandler.
func baseIntermediate(handler Handler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		return handler(r, a)
	}
}
