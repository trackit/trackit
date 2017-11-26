//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package routes

import (
	"encoding/json"
	"net/http"

	"github.com/trackit/trackit2/config"
)

// RegisteredHandlers is the list of all route handlers that were registered.
// Modules providing route handlers are expected to run Register in order to
// populate that list, and the main package is expected to use this list to
// populate its HTTP server.
var RegisteredHandlers = make([]RegisteredHandler, 0, 0x40)

// RegisteredHandler is a registered handler with the pattern it will serve.
type RegisteredHandler struct {
	Pattern string
	Handler
}

// type Handler interface {
// 	Handle(http.ResponseWriter, *http.Request, Arguments) (int, interface{})
// 	Documentation() HandlerDocumentation
// }

type HandlerFunc func(http.ResponseWriter, *http.Request, Arguments) (int, interface{})

type Handler struct {
	Func          HandlerFunc
	Documentation HandlerDocumentation
	methods       map[string]bool
}

// Arguments is a map used by decorators to supply the route handler (or later
// decorators) with additional values.
type Arguments map[interface{}]interface{}

// Decorator is an interface for any type that can decorate an
// IntermediateHandler.
type Decorator interface {
	Decorate(Handler) Handler
}

func resetRegisteredHandlers() {
	RegisteredHandlers = RegisteredHandlers[:0]
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arguments := make(Arguments)
	status, output := h.Func(w, r, arguments)
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.WriteHeader(status)
	e := json.NewEncoder(w)
	if config.PrettyJsonResponses {
		e.SetIndent("", "\t")
	}
	e.Encode(output)
}

func (h Handler) With(ds ...Decorator) Handler {
	l := len(ds) - 1
	for i := range ds {
		h = ds[l-i].Decorate(h)
	}
	return h
}

func (h Handler) Register(pattern string) Handler {
	RegisteredHandlers = append(RegisteredHandlers, RegisteredHandler{
		Pattern: pattern,
		Handler: h,
	})
	return h
}
