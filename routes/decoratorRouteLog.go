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
	"context"
	"net/http"
	"time"

	"github.com/trackit/jsonlog"
)

// requestLogData is the set of fields from an HTTP request that are logged
// when a request arrives on a route decorated by WithRouteLogging.
type requestLogData struct {
	Proto     string    `json:"protocol"`
	Method    string    `json:"method"`
	URL       string    `json:"url"`
	Address   string    `json:"address"`
	Host      string    `json:"host"`
	UserAgent []string  `json:"userAgent"`
	Time      time.Time `json:"time"`
}

// responseLogData is the set of values to be logged once the response to a
// request has been produced.
type responseLogData struct {
	Status      int           `json:"status"`
	Nanoseconds time.Duration `json:"nanoseconds"`
}

// RouteLog is a decorator which logs any calls to the route, with some data
// about the request.
type RouteLog struct{}

func (rl RouteLog) Decorate(h Handler) Handler {
	h.Func = rl.getFunc(h.Func)
	return h
}

// getFunc builds the route handler function for RouteLog.Decorate.
func (_ RouteLog) getFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, response interface{}) {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		rqd := getRequestLogData(r)
		l.Info("Received request.", rqd)
		status, response = hf(w, r, a)
		rsd := getResponseLogData(rqd, status, response)
		l.Info("Produced response to request.", rsd)
		return
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
		Time:      time.Now(),
	}
}

// getResponseLogData fills a responseLogData instance with information about
// how a request was handled.
func getResponseLogData(rqd requestLogData, status int, response interface{}) (rsd responseLogData) {
	rsd.Nanoseconds = time.Now().Sub(rqd.Time)
	rsd.Status = status
	return
}
