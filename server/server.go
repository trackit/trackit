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
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/costs"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
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
	logger := jsonlog.DefaultLogger
	handleDecoratedFunc("/costs", costs.HandleRequest)
	handleDecoratedFunc("/login", users.LogIn)
	handleDecoratedFunc("/test", users.TestToken)
	for _, rh := range routes.RegisteredHandlers {
		handleDecoratedFunc(rh.Pattern, rh.Handler)
		logger.Info(fmt.Sprintf("Registered route %s.", rh.Pattern), nil)
	}
}

// handleDecoratedFunc decorates an HTTP handler function that accepts a logger
// in order to provide the logger and some values in the context.
func handleDecoratedFunc(pattern string, f http.HandlerFunc) {
	http.HandleFunc(pattern, decorateWithLogger(loggedHandler(f)))
}

// requestLogData is the set of fields from an HTTP request that are logged
// when a request arrives.
type requestLogData struct {
	Proto     string   `json:"protocol"`
	Method    string   `json:"method"`
	URL       string   `json:"url"`
	Address   string   `json:"address"`
	Host      string   `json:"host"`
	UserAgent []string `json:"userAgent"`
}

// getRequestLogData fills a requestLogData instance with information from an
// http.Request.
func getRequestLogData(r *http.Request) requestLogData {
	return requestLogData{
		Proto:     r.Proto,
		Method:    r.Method,
		Host:      r.Host,
		URL:       r.URL.String(),
		Address:   r.RemoteAddr,
		UserAgent: r.Header["User-Agent"],
	}
}

// loggedHandler is a decorator around an HTTP handler function with logger
// which logs the request before calling the handler function.
func loggedHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		l.Info("Received request.", getRequestLogData(r))
		f(w, r)
	}
}

// decorateWithLogger decorates an HTTP handler function with logger. It sets
// some values in the context and configures the logger to log them.
func decorateWithLogger(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyRequestId, rand.Int63())
		ctx = context.WithValue(ctx, contextKeyRequestTime, time.Now())
		logger := jsonlog.DefaultLogger
		logger = logger.WithContext(ctx)
		logger = logger.WithContextKey(contextKeyRequestId, "requestId")
		logger = logger.WithContextKey(contextKeyRequestTime, "requestTime")
		logger = logger.WithLogLevel(jsonlog.LogLevelDebug)
		ctx = jsonlog.ContextWithLogger(ctx, logger)
		f(w, r.WithContext(ctx))
	}
}
