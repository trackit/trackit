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

package main

import (
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/routes"
	_ "github.com/trackit/trackit2/users"
)

// contextKey represents a key in a context. Using an unexported type in this
// fashion ensures there can be no collision with a key from some other
// package.
type contextKey int

const (
	// contextKeyRequestId is the key for a request's random ID stored in
	// its context.
	contextKeyRequestId = contextKey(iota)
	// contextKeyRequestTime is the key for the time a request was
	// received, which is stored in its context.
	contextKeyRequestTime
)

func main() {
	logger := jsonlog.DefaultLogger
	configuration := config.LoadConfiguration()
	initializeHandlers()
	logger.Info(fmt.Sprintf("Listening on %s.", configuration.HTTPAddress), nil)
	err := http.ListenAndServe(configuration.HTTPAddress, nil)
	logger.Error("Server stopped.", err.Error())
}

// initializeHandlers sets the HTTP server up with handler functions.
func initializeHandlers() {
	c := config.LoadConfiguration()
	globalDecorators := []routes.Decorator{
		WithRequestTime{},
		WithRequestId{},
		WithBackendId{c.BackendId},
		WithRouteLogging{},
		routes.WithErrorBody{},
	}
	logger := jsonlog.DefaultLogger
	for _, rh := range routes.RegisteredHandlers {
		applyDecoratorsAndHandle(rh.Pattern, rh.Handler, globalDecorators)
		logger.Info(fmt.Sprintf("Registered route %s.", rh.Pattern), nil)
	}
}

func applyDecoratorsAndHandle(p string, h routes.IntermediateHandler, ds []routes.Decorator) {
	h = routes.ApplyDecorators(h, ds...)
	http.Handle(p, h)
}
