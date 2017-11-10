package routes

import (
	"net/http"
	"strconv"
	"strings"
)

type (
	parseIndex uint

	ArgName string

	ParseInt struct {
		ArgName ArgName
	}
	ParseFloat32 struct {
		ArgName ArgName
	}
	ParseFloat64 struct {
		ArgName ArgName
	}
	ParseString struct {
		ArgName ArgName
	}
	ParseVoid struct {
		Required string
	}
	ParseEnd struct{}
)

const (
	index = parseIndex(iota)
)

var (
	badURLTypeMessage  = ErrorBody{"Bad URL type."}
	badURLTypeCode     = 400
	notFoundURLMessage = ErrorBody{"Not Found."}
	notFoundURLCode    = 404
)

func incrementIndex(a Arguments) int {
	if _, ok := a[index]; ok {
		a[index] = a[index].(int) + 1
	} else {
		a[index] = 1
	}
	return a[index].(int)
}

func (d ParseInt) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		splitURL := strings.Split(r.URL.Path, "/")
		index := incrementIndex(a) - 1
		if len(splitURL) > index {
			if i, err := strconv.Atoi(splitURL[index]); err == nil {
				a[d.ArgName] = i
				return h(w, r, a)
			}
			return badURLTypeCode, badURLTypeMessage
		}
		return notFoundURLCode, notFoundURLMessage
	}
}

func (d ParseFloat32) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		splitURL := strings.Split(r.URL.Path, "/")
		index := incrementIndex(a) - 1
		if len(splitURL) > index {
			if i, err := strconv.ParseFloat(splitURL[index], 32); err == nil {
				a[d.ArgName] = i
				return h(w, r, a)
			}
			return badURLTypeCode, badURLTypeMessage
		}
		return notFoundURLCode, notFoundURLMessage
	}
}

func (d ParseFloat64) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		splitURL := strings.Split(r.URL.Path, "/")
		index := incrementIndex(a) - 1
		if len(splitURL) > index {
			if i, err := strconv.ParseFloat(splitURL[index], 64); err == nil {
				a[d.ArgName] = i
				return h(w, r, a)
			}
			return badURLTypeCode, badURLTypeMessage
		}
		return notFoundURLCode, notFoundURLMessage
	}
}

func (d ParseString) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		splitURL := strings.Split(r.URL.Path, "/")
		index := incrementIndex(a) - 1
		if len(splitURL) > index {
			a[d.ArgName] = splitURL[index]
			return h(w, r, a)
		}
		return notFoundURLCode, notFoundURLMessage
	}
}

func (d ParseVoid) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		splitURL := strings.Split(r.URL.Path, "/")
		index := incrementIndex(a) - 1
		if (len(splitURL) == index && d.Required == "") ||
			(len(splitURL) > index && splitURL[index] == d.Required) {
			return h(w, r, a)
		}
		return notFoundURLCode, notFoundURLMessage
	}
}

func (d ParseEnd) Decorate(h IntermediateHandler) IntermediateHandler {
	return ParseVoid{""}.Decorate(h)
}
