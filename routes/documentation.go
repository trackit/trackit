package routes

import (
	"net/http"
)

const (
	ComponentMethod      = "method"
	ComponentRequirement = "requirement"
	ComponentRequest     = "request"
	ComponentResponse    = "response"
	ComponentBody        = "body"
)

type Tags map[string][]string

type HandlerDocumentationBody struct {
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	Tags        Tags   `json:"tags,omitempty"`
}

type HandlerDocumentation struct {
	HandlerDocumentationBody
	Components map[string]HandlerDocumentation `json:"components,omitempty"`
}

var documentationHandler = MethodMuxer{
	http.MethodGet: H(getDocumentation).With(Documentation{
		Summary: "get the api's documentation",
		Description: ("Get the api's documentation in structured (JSON) " +
			"format. This documentation is automatically " +
			"generated from the definition of the route handlers " +
			"and thus should always be up to date. The same " +
			"documentation can be obtained for specific routes " +
			"using the OPTIONS request on them."),
	}),
}.H()

func DocumentationHandler() Handler {
	return documentationHandler
}

func getDocumentation(_ *http.Request, _ Arguments) (int, interface{}) {
	routes := make(map[string]HandlerDocumentation)
	for _, rh := range RegisteredHandlers {
		routes[rh.Pattern] = rh.Handler.Documentation
	}
	return http.StatusOK, routes
}
