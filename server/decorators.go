package main

import (
	"context"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/routes"
)

type WithRequestId struct{}

type WithRequestTime struct{}

type WithBackendId struct {
	backendId string
}

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
// when a request arrives.
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
