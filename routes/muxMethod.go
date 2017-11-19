package routes

import (
	"fmt"
	"net/http"
)

const (
	ErrMethodNotAllowed = constError("Method is not allowed.")
)

type MethodMuxer map[string]Handler

func (mm MethodMuxer) H() Handler {
	return Handler{
		Func:          mm.handlerFunc(),
		Documentation: mm.documentation(),
		methods:       mm.methods(),
	}
}

func (mm MethodMuxer) handlerFunc() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		if h, ok := mm[r.Method]; ok {
			return h.Func(w, r, a)
		} else {
			return http.StatusMethodNotAllowed, ErrMethodNotAllowed
		}
	}
}

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

func (mm MethodMuxer) methods() map[string]bool {
	o := make(map[string]bool)
	for m := range mm {
		o[m] = true
	}
	return o
}
