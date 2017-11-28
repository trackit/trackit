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
	"testing/quick"
)

func bidIdentity(s string) []string { return []string{s} }

func TestDecoratorBackindId(t *testing.T) {
	f := func(backendId string) []string {
		h := H(getFoo).With(
			BackendId{backendId},
		)
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
		return response.Header()["X-Backend-ID"]
	}
	if err := quick.CheckEqual(bidIdentity, f, nil); err != nil {
		t.Error(err)
	}
}
