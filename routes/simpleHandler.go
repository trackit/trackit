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

// SimpleHandler is the type a route handler must have.
type SimpleHandler func(*http.Request, Arguments) (int, interface{})

// H builds a Handler from a SimpleHandler. The produced handler has an empty
// documentation.
func H(h SimpleHandler) Handler {
	return Handler{
		Func: func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
			return h(r, a)
		},
	}
}
