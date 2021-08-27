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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
	// TagRequiredQueryArg is the tag used to document required query args.
	TagRequiredQueryArg = "required:allof:queryarg"
	// TagOptionalQueryArg is the tag used to document optional query args.
	TagOptionalQueryArg = "optional:allof:queryarg"
	// iso8601DateFormat is the time format for the ISO8601 format
	iso8601DateFormat = "2006-01-02"
)

type (
	// QueryArgBool denotes an int query argument. It fulfills the QueryParser
	// interface.
	QueryArgBool struct{}

	// QueryArgInt denotes an int query argument. It fulfills the QueryParser
	// interface.
	QueryArgInt struct{}

	// QueryArgUint denotes a uint query argument. It fulfills the QueryParser
	// interface.
	QueryArgUint struct{}

	// QueryArgString denotes a string query argument. It fulfills the QueryParser
	// interface.
	QueryArgString struct{}

	// QueryArgIntSlice denotes a []int query argument. It fulfills the QueryParser
	// interface.
	QueryArgIntSlice struct{}

	// QueryArgUintSlice denotes a []uint query argument. It fulfills the QueryParser
	// interface.
	QueryArgUintSlice struct{}

	// QueryArgStringSlice denotes a []string query argument. It fulfills the
	// QueryParser interface.
	QueryArgStringSlice struct{}

	// QueryArgDate denotes a time.Time query argument. The time format is
	// the ISO8601. It fulfills the QueryParser interface.
	QueryArgDate struct{}

	// QueryParser parses a string and returns a typed value. An error can
	// be returned if the value could not be parsed.
	QueryParser interface {
		// QueryParse parses a query string argument.
		QueryParse(string) (interface{}, error)
		// FormatName returns the name of the parameter type.
		FormatName() string
	}

	// QueryArg defines an argument by its name and its type. A description
	// can be used for documentation purposes. The argument can be optional
	// by changing the Optional value to true.
	QueryArg struct {
		Name        string
		Description string
		Type        QueryParser
		Optional    bool
	}

	// QueryArgs contains all the arguments to parse in the URL.
	// QueryArgs has a method Decorate called to apply the
	// decorators on an endpoint.
	QueryArgs []QueryArg
)

func (d QueryArgBool) FormatName() string        { return "bool" }
func (d QueryArgInt) FormatName() string         { return "int" }
func (d QueryArgUint) FormatName() string        { return "uint" }
func (d QueryArgString) FormatName() string      { return "string" }
func (d QueryArgIntSlice) FormatName() string    { return "[]int" }
func (d QueryArgUintSlice) FormatName() string   { return "[]uint" }
func (d QueryArgStringSlice) FormatName() string { return "[]string" }
func (d QueryArgDate) FormatName() string        { return "time.Time" }

// QueryParse parses an int. A nil error indicates a success. With this func,
// QueryArgBool fulfills QueryArgType.
func (QueryArgBool) QueryParse(val string) (interface{}, error) {
	if val == "" {
		return true, nil
	} else if res, err := strconv.ParseBool(val); err == nil {
		return res, nil
	}
	return nil, errors.New("must be a bool")
}

// QueryParse parses an int. A nil error indicates a success. With this func,
// QueryArgInt fulfills QueryArgType.
func (QueryArgInt) QueryParse(val string) (interface{}, error) {
	if i, err := strconv.ParseInt(val, 10, 64); err == nil &&
		i <= int64(maxInt) && i >= int64(minInt) {
		return int(i), nil
	}
	return nil, errors.New("must be an int")
}

// QueryParse parses a uint. A nil error indicates a success. With this func,
// QueryArgUint fulfills QueryArgType.
func (QueryArgUint) QueryParse(val string) (interface{}, error) {
	if i, err := strconv.ParseUint(val, 10, 64); err == nil &&
		i <= uint64(maxUint) {
		return uint(i), nil
	}
	return nil, errors.New("must be a uint")
}

// QueryParse parses a string. A nil error indicates a success. With this func,
// QueryArgString fulfills QueryArgType.
func (QueryArgString) QueryParse(val string) (interface{}, error) {
	return val, nil
}

// QueryParse parses an []int. A nil error indicates a success. With this func,
// QueryArgIntSlice fulfills QueryArgType.
func (QueryArgIntSlice) QueryParse(val string) (interface{}, error) {
	vals := strings.Split(val, ",")
	res := make([]int, 0, len(vals))
	for _, v := range vals {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil &&
			i <= int64(maxInt) && i >= int64(minInt) {
			res = append(res, int(i))
		} else {
			return nil, errors.New("must be a slice of int")
		}
	}
	return res, nil
}

// QueryParse parses an []uint. A nil error indicates a success. With this func,
// QueryArgUintSlice fulfills QueryArgType.
func (QueryArgUintSlice) QueryParse(val string) (interface{}, error) {
	vals := strings.Split(val, ",")
	res := make([]uint, 0, len(vals))
	for _, v := range vals {
		if i, err := strconv.ParseUint(v, 10, 64); err == nil &&
			i <= uint64(maxUint) {
			res = append(res, uint(i))
		} else {
			return nil, errors.New("must be a slice of uint")
		}
	}
	return res, nil
}

// QueryParse parses a []string. Since it does not do any type conversion,
// it cannot fail. With this func, QueryArgStringSlice fulfills QueryArgType.
func (QueryArgStringSlice) QueryParse(val string) (interface{}, error) {
	vals := strings.Split(val, ",")
	res := make([]string, 0, len(vals))
	for _, v := range vals {
		res = append(res, v)
	}
	return res, nil
}

// QueryParse parses a time.Time in the ISO8601 format. With this func,
// QueryArgDate fulfills QueryArgType.
func (QueryArgDate) QueryParse(val string) (interface{}, error) {
	parsedDate, err := time.Parse(iso8601DateFormat, val)
	if err != nil {
		return nil, errors.New("could not parse the date" + err.Error())
	}
	return parsedDate, nil
}

func parseArg(arg QueryArg, r *http.Request, a Arguments) (int, error) {
	if rawVal, ok := r.URL.Query()[arg.Name]; ok &&
		len(rawVal) > 0 &&
		(rawVal[0] != "" || arg.Type.FormatName() == "bool") {
		if val, err := arg.Type.QueryParse(rawVal[0]); err == nil {
			a[arg] = val
		} else {
			return http.StatusBadRequest, fmt.Errorf("query arg '%s': %s", arg.Name, err.Error())
		}
	} else if !arg.Optional {
		return http.StatusBadRequest, fmt.Errorf("query arg '%s': not found", arg.Name)
	}
	return http.StatusOK, nil
}

// Decorate is the function called to apply the decorators to an endpoint. It returns
// a function. This function produces a 400 Bad Request error code with a json error message or
// calls the next IntermediateHandler.
// The goal of this function is to get the URL parameters to store them in
// the Arguments.
func (qa QueryArgs) Decorate(h Handler) Handler {
	h.Func = qa.getFunc(h.Func)
	h.Documentation = qa.getDocumentation(h.Documentation)
	return h
}

// getFunc builds a handler function for QueryArgs.Decorate
func (qa QueryArgs) getFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		for _, arg := range qa {
			if code, err := parseArg(arg, r, a); code != http.StatusOK {
				return code, err
			}
		}
		return hf(w, r, a)
	}
}

// createDocumentationSlice creates the slice of string and fills it with documentation
func createDocumentationSlice(index string, optional bool, qa QueryArgs) []string {
	ts := make([]string, 0)
	for i := range qa {
		if qa[i].Optional == optional {
			ts = append(ts, fmt.Sprintf("%s:%s:%s", qa[i].Name, qa[i].Type.FormatName(), qa[i].Description))
		}
	}
	return ts
}

// getDocumentation builds the documentation for QueryArgs.Decorate
func (qa QueryArgs) getDocumentation(hd HandlerDocumentation) HandlerDocumentation {
	if hd.Tags == nil {
		hd.Tags = make(Tags)
	}
	hd.Tags[TagRequiredQueryArg] = append(hd.Tags[TagRequiredQueryArg], createDocumentationSlice(TagRequiredQueryArg, false, qa)...)
	hd.Tags[TagOptionalQueryArg] = append(hd.Tags[TagOptionalQueryArg], createDocumentationSlice(TagOptionalQueryArg, true, qa)...)
	return hd
}
