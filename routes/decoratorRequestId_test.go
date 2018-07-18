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
	"net/http/httptest"
	"testing"

	"github.com/trackit/trackit-server/config"
)

func init() {
	config.PrettyJsonResponses = true
}

func TestDecoratorRequestId(t *testing.T) {
	h := H(getFoo).With(
		RequestId{},
	)
	ids := make(map[string]bool)
	for i := 0; i < 20; i++ {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		s, r := h.Func(response, request, Arguments{})
		if s != http.StatusOK {
			t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
		}
		if rt, ok := r.(string); ok {
			if rt != getFooResponse {
				t.Errorf(
					"Response should be '%s', is '%s' instead.",
					getFooResponse,
					rt,
				)
			}
		} else {
			t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", getFooResponse, r)
		}
		requestIds := response.Header()["X-Request-ID"]
		if l := len(requestIds); l == 1 {
			requestId := requestIds[0]
			if ids[requestId] {
				t.Error("Duplicate request ID.")
			} else {
				ids[requestId] = true
			}
		} else {
			t.Errorf("Response should have a single requets ID, has %d instead.", l)
		}
	}
}
