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
