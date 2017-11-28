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

const (
	documentationDescription = "Get the api's documentation in " +
		"structured (JSON) format. This documentation is " +
		"automatically generated from the definition of the route " +
		"handlers and thus should always be up to date. The same " +
		"documentation can be obtained for specific routes using the " +
		"OPTIONS request on them."
	documentationSummary = "get the api's documentation"
)

var documentationHandler = MethodMuxer{
	http.MethodGet: H(getDocumentation).With(Documentation{
		Summary:     documentationSummary,
		Description: documentationDescription,
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
