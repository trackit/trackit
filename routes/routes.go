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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/trackit/trackit-server/config"
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

// csvGenerator is an interface for any type that can generate a CSV file content
type csvGenerator interface {
	ToCSVable() [][]string
}

// xlsGenerator is an interface for any type that can generate an xls file content
type xlsGenerator interface {
	GetFileContent() []byte
	GetFileName() string
}

func resetRegisteredHandlers() {
	RegisteredHandlers = RegisteredHandlers[:0]
}

// BUG(zanchi-r): We do not support content type negociation
// Only a single type without quality factor should be specified in the Accept header
// Example : "Accept: text/csv"
// If an unsupported type is specified an empty body will be returned
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arguments := make(Arguments)
	status, output := h.Func(w, r, arguments)
	acceptType := "*/*"
	if len(r.Header["Accept"]) > 0 {
		acceptType = r.Header["Accept"][0]
	}
	// If no specific accept type is specified we use application/json by default
	if strings.HasPrefix(acceptType, "*/*") {
		acceptType = "application/json"
	}
	w.Header()["Content-Type"] = []string{fmt.Sprintf("%s; charset=utf-8", acceptType)}
	switch acceptType {
	case "application/json":
		w.WriteHeader(status)
		e := json.NewEncoder(w)
		if config.PrettyJsonResponses {
			e.SetIndent("", "\t")
		}
		e.Encode(output)
	case "text/csv":
		w.Header().Set("Content-Disposition", "attachment; filename=trackit.csv")
		w.WriteHeader(status)
		if outputGen, ok := output.(csvGenerator); ok {
			csvWriter := csv.NewWriter(w)
			csvWriter.WriteAll(outputGen.ToCSVable())
		} else {
			// TODO: if the data do not implement the csvGenerator interface, try to generate it by reflection
		}
	case "application/vnd.ms-excel":
		if outputGen, ok := output.(xlsGenerator); ok {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputGen.GetFileName()))
			w.WriteHeader(status)
			w.Write(outputGen.GetFileContent())
		} else {
			if status == http.StatusOK {
				w.WriteHeader(http.StatusNotImplemented)
			} else {
				w.WriteHeader(status)
			}
		}
	}
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
