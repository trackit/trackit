package aws

import (
	"encoding/json"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

func init() {
	routes.Register(
		"/aws",
		routeAws,
		routes.RequireMethod{"POST", "GET"},
		routes.RequireContentType{"application/json"},
		db.WithTransaction{db.Db},
		users.WithAuthenticatedUser{},
	)
	routes.Register(
		"/aws/next",
		nextExternal,
		routes.RequireMethod{"GET"},
		db.WithTransaction{db.Db},
		users.WithAuthenticatedUser{},
	)
}

// routeAws is a route handler for /aws. It delegates the handling to
// postAwsAccount or getAwsAccount depending on the method from the HTTP
// request.
func routeAws(r *http.Request, a routes.Arguments) (int, interface{}) {
	switch r.Method {
	case "POST":
		return postAwsAccount(r, a)
	case "GET":
		return getAwsAccount(r, a)
	default:
		logger := jsonlog.LoggerFromContextOrDefault(r.Context())
		logger.Error("Bad method. Did 'RequireMethod' do its job?", r.Method)
		return 500, nil
	}
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}
