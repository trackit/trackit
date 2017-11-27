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
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/routes"
)

// WithRequestId is a decorator which adds a request ID to incoming requests.
// It stores it in the request context, sets the logger up to log it, and adds
// it to the response in the `X-Request-ID` HTTP header.
type WithRequestId struct{}

// WithRequestTime is a decorator which adds a request time to incoming
// requests. It stores it in the request context and sets the logger up to log
// it.
type WithRequestTime struct{}

// WithBackendId is a decorator which adds a backend ID to incoming requests.
// It adds it to tthe response in the `X-Backend-ID` HTTP header.
type WithBackendId struct {
	backendId string
}

// WithRouteLogging is a decorator which logs any calls to the route, with some
// data about the request.
type WithRouteLogging struct{}

func (d WithRequestId) Decorate(h routes.IntermediateHandler) routes.IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		requestId := uuid.NewV1().String()
		r = requestWithLoggedContextValue(r, contextKeyRequestId, "requestId", requestId)
		w.Header()["X-Request-ID"] = []string{requestId}
		return h(w, r, a)
	}
}

func (d WithBackendId) Decorate(h routes.IntermediateHandler) routes.IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		w.Header()["X-Backend-ID"] = []string{d.backendId}
		return h(w, r, a)
	}
}

func (d WithRouteLogging) Decorate(h routes.IntermediateHandler) routes.IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		l.Info("Received request.", getRequestLogData(r))
		return h(w, r, a)
	}
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

// requestLogData is the set of fields from an HTTP request that are logged
// when a request arrives on a route decorated by WithRouteLogging.
type requestLogData struct {
	Proto     string   `json:"protocol"`
	Method    string   `json:"method"`
	URL       string   `json:"url"`
	Address   string   `json:"address"`
	Host      string   `json:"host"`
	UserAgent []string `json:"userAgent"`
}

func (d WithRequestTime) Decorate(h routes.IntermediateHandler) routes.IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a routes.Arguments) (int, interface{}) {
		now := time.Now()
		r = requestWithLoggedContextValue(r, contextKeyRequestTime, "requestTime", now)
		return h(w, r, a)
	}
}

// requestWithLoggedContextValue returns a request with 'value' in its context
// at key 'contextKey', with a logger logging that value with key 'logKey'.
func requestWithLoggedContextValue(r *http.Request, contextKey interface{}, logKey string, value interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKey, value)
	log := jsonlog.LoggerFromContextOrDefault(ctx)
	log = log.WithContext(ctx)
	log = log.WithContextKey(contextKey, logKey)
	ctx = jsonlog.ContextWithLogger(ctx, log)
	return r.WithContext(ctx)
}
