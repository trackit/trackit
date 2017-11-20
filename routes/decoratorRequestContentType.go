package routes

import (
	"net/http"
)

const (
	ErrUnsupportedContentType = constError("bad content type")
	ErrMultipleContentTypes   = constError("multiple content types")
	ErrMissingContentType     = constError("missing content type")
	TagRequiredContentType    = "require:contenttype"
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
			return http.StatusBadRequest, ErrMissingContentType
		} else if len(contentTypes) > 1 {
			return http.StatusBadRequest, ErrMultipleContentTypes
		} else if supportedContentTypes[contentTypes[0]] {
			return hf(w, r, a)
		} else {
			return http.StatusBadRequest, ErrUnsupportedContentType
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
	hd.Tags[TagRequiredContentType] = []string(rct)
	return hd
}
