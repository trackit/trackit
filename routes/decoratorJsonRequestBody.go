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

type JsonRequestBody struct {
	Example interface{}
}

func (jrb JsonRequestBody) Decorate(h Handler) Handler {
	h.Func = jrb.getFunc(h.Func)
	h.Documentation = jrb.getDocumentation(h.Documentation)
	return h
}

func (jrb JsonRequestBody) getFunc(hf HandlerFunc) HandlerFunc {
	validate, err := req.CreateValidator(jrb.Example)
	var handleWithValidation func(http.ResponseWriter, *http.Request, Arguments, reflect.Value) (int, interface{})
	if err != nil {
		logger := jsonlog.DefaultLogger
		logger.Error("Failed to build validator for type %T.", jrb.Example)
		os.Exit(1)
	}
	if validate == nil {
		handleWithValidation = func(w http.ResponseWriter, r *http.Request, a Arguments, body reflect.Value) (int, interface{}) {
			a[argumentKeyJsonBody] = reflect.Indirect(body).Interface()
			return hf(w, r, a)
		}
	} else {
		handleWithValidation = func(w http.ResponseWriter, r *http.Request, a Arguments, body reflect.Value) (int, interface{}) {
			logger := jsonlog.LoggerFromContextOrDefault(r.Context())
			err := validate(body.Interface())
			if err == nil {
				a[argumentKeyJsonBody] = reflect.Indirect(body).Interface()
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
		body := reflect.New(reflect.TypeOf(jrb.Example))
		if err := json.NewDecoder(r.Body).Decode(body.Interface()); err != nil {
			logger.Warning("Failed to parse request body.", err.Error())
			return http.StatusBadRequest, errors.New("failed to parse request body")
		} else {
			return handleWithValidation(w, r, a, body)
		}
	}
}

func (jrb JsonRequestBody) getDocumentation(hd HandlerDocumentation) HandlerDocumentation {
	if hd.Components == nil {
		hd.Components = make(map[string]HandlerDocumentation)
	}
	hd.Components["input:body:example"] = HandlerDocumentation{
		HandlerDocumentationBody: HandlerDocumentationBody{
			Summary:     "input body example",
			Description: jrb.getExampleString(),
		},
	}
	return hd
}

func (jrb JsonRequestBody) getExampleString() string {
	bytes, err := json.MarshalIndent(jrb.Example, "", "\t")
	if err != nil {
		return "FAIL"
	} else {
		return string(bytes)
	}
}
