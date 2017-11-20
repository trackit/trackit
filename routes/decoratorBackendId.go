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

// BackendId is a decorator which adds a backend ID to incoming requests. It
// adds it to tthe response in the `X-Backend-ID` HTTP header.
type BackendId struct {
	BackendId string
}

func (d BackendId) Decorate(h Handler) Handler {
	h.Func = d.getFunc(h.Func)
	return h
}

// getFunc returns a decorated handler function for BackendId.
func (d BackendId) getFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		w.Header()["X-Backend-ID"] = []string{d.BackendId}
		return hf(w, r, a)
	}
}
