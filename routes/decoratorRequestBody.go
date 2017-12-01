package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/util/req"
)

type RequestBody struct {
	Example interface{}
}

func (rb RequestBody) Decorate(h Handler) Handler {
	h.Func = rb.getFunc(h.Func)
	h.Documentation = rb.getDocumentation(h.Documentation)
	return h
}

func (rb RequestBody) getFunc(hf HandlerFunc) HandlerFunc {
	validate, err := req.CreateValidator(rb.Example)
	var handleWithValidation func(http.ResponseWriter, *http.Request, Arguments, reflect.Value) (int, interface{})
	if err != nil {
		logger := jsonlog.DefaultLogger
		logger.Error("Failed to build validator for type %T.", rb.Example)
		os.Exit(1)
	}
	if validate == nil {
		handleWithValidation = func(w http.ResponseWriter, r *http.Request, a Arguments, body reflect.Value) (int, interface{}) {
			a[argumentKeyBody] = reflect.Indirect(body).Interface()
			return hf(w, r, a)
		}
	} else {
		handleWithValidation = func(w http.ResponseWriter, r *http.Request, a Arguments, body reflect.Value) (int, interface{}) {
			logger := jsonlog.LoggerFromContextOrDefault(r.Context())
			err := validate(body.Interface())
			if err == nil {
				a[argumentKeyBody] = reflect.Indirect(body)
				return hf(w, r, a)
			} else if verr, ok := err.(req.ValidationError); ok {
				return http.StatusBadRequest, verr
			} else {
				logger.Error("Abnormal validation failure.", err.Error())
				return http.StatusInternalServerError, errors.New("failed to parse request body")
			}
		}
	}
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		logger := jsonlog.LoggerFromContextOrDefault(r.Context())
		body := reflect.New(reflect.TypeOf(rb.Example))
		if err := json.NewDecoder(r.Body).Decode(body.Interface()); err != nil {
			logger.Warning("Failed to parse request body.", err.Error())
			return http.StatusBadRequest, errors.New("failed to parse request body")
		} else {
			return handleWithValidation(w, r, a, body)
		}
	}
}

func (rb RequestBody) getDocumentation(hd HandlerDocumentation) HandlerDocumentation {
	if hd.Components == nil {
		hd.Components = make(map[string]HandlerDocumentation)
	}
	hd.Components["input:body:example"] = HandlerDocumentation{
		HandlerDocumentationBody: HandlerDocumentationBody{
			Summary:     "input body example",
			Description: rb.getExampleString(),
		},
	}
	return hd
}

func (rb RequestBody) getExampleString() string {
	bytes, err := json.MarshalIndent(rb.Example, "", "\t")
	if err != nil {
		return "FAIL"
	} else {
		return string(bytes)
	}
}

func MustRequestBody(a Arguments, ptr interface{}) {
	err := GetRequestBody(a, ptr)
	if err != nil {
		panic(err)
	}
}

func GetRequestBody(a Arguments, ptr interface{}) error {
	body := a[argumentKeyBody]
	if body == nil {
		return errors.New("request body not found")
	} else {
		return copyBodyTo(body.(reflect.Value), ptr)
	}
}

func copyBodyTo(val reflect.Value, dst interface{}) error {
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return errors.New("destination is not pointer")
	} else if dstType.Elem() != val.Type() {
		return errors.New("incompatible types")
	} else {
		reflect.Indirect(reflect.ValueOf(dst)).Set(val)
		return nil
	}
}
