package routes

import (
	"fmt"
	"net/http"
)

const (
	ErrMethodNotAllowed = constError("Method is not allowed.")
)

// MethodMuxer multiplexes requests based on their method. If a request arrives
// with a method not in the map, MethodMuxer responds with
// http.StatusMethodNotAllowed.
type MethodMuxer map[string]Handler

// H builds a handler from the MethodMuxer. If the MethodMuxer has only one
// element, the summary of the single handler's documentation will be copied
// over to the MethodMuxer's.
func (mm MethodMuxer) H() Handler {
	return Handler{
		Func:          mm.handlerFunc(),
		Documentation: mm.documentation(),
		methods:       mm.methods(),
	}
}

// handlerFunc builds the handler function for MethodMuxer.H.
func (mm MethodMuxer) handlerFunc() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		if h, ok := mm[r.Method]; ok {
			return h.Func(w, r, a)
		} else {
			return http.StatusMethodNotAllowed, ErrMethodNotAllowed
		}
	}
}

// documentation builds the documentation for MethodMuxer.H.
func (mm MethodMuxer) documentation() HandlerDocumentation {
	var hd HandlerDocumentation
	hd.Components = make(map[string]HandlerDocumentation)
	for m, h := range mm {
		m = fmt.Sprintf("method:%s", m)
		hd.Components[m] = h.Documentation
		if len(mm) == 1 {
			hd.Summary = h.Documentation.Summary
		}
	}
	return hd
}

// methods compiles a map of methods supported by the MethodMuxer.
func (mm MethodMuxer) methods() map[string]bool {
	o := make(map[string]bool)
	for m := range mm {
		o[m] = true
	}
	return o
}
