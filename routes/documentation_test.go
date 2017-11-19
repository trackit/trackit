package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testDocumentationExpected = `{"/doc":{"summary":"get the api's documentation","components":{"GET":{"summary":"get the api's documentation","description":"Get the api's documentation in structured (JSON) format. This documentation is automatically generated from the definition of the route handlers and thus should always be up to date. The same documentation can be obtained for specific routes using the OPTIONS request on them.","componentType":"method"}}}}`

func TestDocumentationHandler(t *testing.T) {
	h := DocumentationHandler()
	h.Register("/doc")
	status, response := h.Func(nil, httptest.NewRequest(http.MethodGet, "/doc", nil), nil)
	if status != http.StatusOK {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, status)
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Error should be nil, is '%s' instead.", err.Error())
	} else if string(bytes) != testDocumentationExpected {
		t.Errorf("Documentation should be\n%s\nis\n$s\ninstead.", testDocumentationExpected, string(bytes))
	}
	resetRegisteredHandlers()
}
