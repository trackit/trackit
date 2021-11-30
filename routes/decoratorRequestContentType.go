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

const (
	errUnsupportedContentType = constError("bad content type")
	errMultipleContentTypes   = constError("multiple content types")
	errMissingContentType     = constError("missing content type")
	tagRequiredContentType    = "required:contenttype"
)

// RequestContentType decorates handlers and requires requests to be of a given
// content type. Requests with a content type not in the RequestContentType
// list will be rejected with a http.StatusBadRequest.
type RequestContentType []string

func (rct RequestContentType) Decorate(h Handler) Handler {
	h.Func = rct.getFunc(h.Func)
	h.Documentation = rct.getDocumentation(h.Documentation)
	return h
}

// getFunc builds the handler function for RequestContentType.Decorate
func (rct RequestContentType) getFunc(hf HandlerFunc) HandlerFunc {
	supportedContentTypes := rct.getCtMap()
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		if bodyIsIgnoredForMethod(r.Method) {
			return hf(w, r, a)
		} else if contentTypes := r.Header["Content-Type"]; len(contentTypes) == 0 {
			return http.StatusBadRequest, errMissingContentType
		} else if len(contentTypes) > 1 {
			return http.StatusBadRequest, errMultipleContentTypes
		} else if supportedContentTypes[contentTypes[0]] {
			return hf(w, r, a)
		} else {
			return http.StatusBadRequest, errUnsupportedContentType
		}
	}
}

func (rct RequestContentType) getCtMap() map[string]bool {
	ctMap := make(map[string]bool)
	for _, ct := range rct {
		ctMap[ct] = true
	}
	return ctMap
}

func bodyIsIgnoredForMethod(method string) bool {
	switch method {
	case http.MethodGet:
		fallthrough
	case http.MethodHead:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodOptions:
		return true
	default:
		return false
	}
}

func (rct RequestContentType) getDocumentation(hd HandlerDocumentation) HandlerDocumentation {
	if hd.Tags == nil {
		hd.Tags = make(Tags)
	}
	hd.Tags[tagRequiredContentType] = []string(rct)
	return hd
}
