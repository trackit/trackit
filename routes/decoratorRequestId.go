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

	"github.com/satori/go.uuid"
)

// RequestId is a decorator which adds a request ID to incoming requests.  It
// stores it in the request context, sets the logger up to log it, and adds it
// to the response in the `X-Request-ID` HTTP header.
type RequestId struct{}

func (ri RequestId) Decorate(h Handler) Handler {
	h.Func = ri.getFunc(h.Func)
	return h
}

// getFunc builds the handler function for RequestId.Decorate
func (RequestId) getFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		requestId := uuid.NewV1().String()
		r = requestWithLoggedContextValue(r, contextKeyRequestId, "requestId", requestId)
		w.Header()["X-Request-ID"] = []string{requestId}
		return hf(w, r, a)
	}
}
