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

// RequireContentType is a route decorator to produce a 400 'Bad Request'
//error if a request uses a content type not present in its backing slice.
type RequireContentType []string

// badContentType is the error body sent to the client in case the decorater
// rejects their request.
var badContentType = ErrorBody{"Bad content type."}

func (d RequireContentType) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		if bodyIsIgnoredForMethod(r.Method) {
			return h(w, r, a)
		} else if contentTypes := r.Header["Content-Type"]; len(contentTypes) == 1 {
			contentType := contentTypes[0]
			for _, ct := range d {
				if contentType == ct {
					return h(w, r, a)
				}
			}
		}
		return 400, badContentType
	}
}

func bodyIsIgnoredForMethod(method string) bool {
	switch method {
	case "GET":
		return true
	case "OPTIONS":
		return true
	default:
		return false
	}
}
