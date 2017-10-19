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
	default:
		return false
	}
}
