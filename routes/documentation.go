package routes

import (
	"net/http"
)

// Tags is a map of tags on a documentation.
type Tags map[string][]string

// HandlerDocumentationBody represents the values forming the body of a
// documentation component.
type HandlerDocumentationBody struct {
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	Tags        Tags   `json:"tags,omitempty"`
}

// HandlerDocumentation represent a handler's documentation. It can have a tree
// of component subdocumentations.
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

// DocumentationHandler returns a Handler which responds to http.MehodGet
// requests with a JSON representation of all registered routes.
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
