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

package csv

import (
	"encoding/csv"
	"io"
	"reflect"
	"strings"
)

type csvError string

func (err csvError) Error() string {
	return string(err)
}

const (
	ErrNotStruct               = csvError("Record is not a struct.")
	ErrNotSupported            = csvError("Operation is not supported.")
	ErrInvalidFieldType        = csvError("Field type not supported.")
	ErrInvalidDefaultFieldType = csvError("Default field type not supported.")
	StructFieldTag             = "csv"
)

type recordType struct {
	reflect.Type
	fieldByName  map[string]int
	nameByField  map[int]string
	defaultField int
}

var recordTypeCache map[reflect.Type]recordType = make(map[reflect.Type]recordType)

type Decoder struct {
	reader *csv.Reader
	header []string
}

func NewDecoder(r io.Reader) Decoder {
	return Decoder{
		reader: csv.NewReader(r),
	}
}

func (d *Decoder) SetHeader(h []string) {
	d.header = h
}

func (d *Decoder) ReadHeader() error {
	if record, err := d.reader.Read(); err == nil {
		d.SetHeader(record)
		return nil
	} else {
		return err
	}
}

func (d *Decoder) ReadRecord(v interface{}) error {
	if rt, err := getRecordType(v); err != nil {
		return err
	} else if record, err := d.reader.Read(); err != nil {
		return err
	} else {
		return d.storeRecord(rt, v, record)
	}
}

func (d *Decoder) storeRecord(rt recordType, vi interface{}, record []string) error {
	v := reflect.ValueOf(vi)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i, h := range d.header {
		if fi, ok := rt.fieldByName[h]; ok {
			v.FieldByIndex([]int{fi}).SetString(record[i])
		} else if rt.defaultField >= 0 {
			fv := v.FieldByIndex([]int{rt.defaultField})
			if fv.IsNil() {
				fv.Set(reflect.MakeMap(reflect.TypeOf(map[string]string(nil))))
			}
			fv.SetMapIndex(reflect.ValueOf(h), reflect.ValueOf(record[i]))
		}
	}
	return nil
}

func getRecordType(v interface{}) (recordType, error) {
	t := reflect.TypeOf(v)
	if rt, ok := recordTypeCache[t]; ok {
		return rt, nil
	} else {
		rt, err := buildRecordType(t)
		if err == nil {
			recordTypeCache[t] = rt
		}
		return rt, err
	}
}

func buildRecordType(t reflect.Type) (recordType, error) {
	var rt = recordType{defaultField: -1}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return rt, ErrNotStruct
	} else {
		rt.Type = t
		rt.fieldByName = make(map[string]int)
		rt.nameByField = make(map[int]string)
		numField := rt.NumField()
		for i := 0; i < numField; i++ {
			field := rt.FieldByIndex([]int{i})
			name := getFieldName(field)
			if name == "" {
				rt.defaultField = i
			} else if name != "-" {
				rt.fieldByName[name] = i
				rt.nameByField[i] = name
			}
		}
		return validateRecordType(rt)
	}
}

func getFieldName(field reflect.StructField) string {
	tag := strings.Split(field.Tag.Get(StructFieldTag), ",")
	if len(tag) >= 2 && tag[1] == "any" {
		return ""
	} else if len(tag) > 0 && tag[0] != "" {
		return tag[0]
	} else {
		return field.Name
	}
}

func validateRecordType(rt recordType) (recordType, error) {
	if df := rt.defaultField; df >= 0 {
		if rt.Type.FieldByIndex([]int{df}).Type != reflect.TypeOf(map[string]string{}) {
			return rt, ErrInvalidDefaultFieldType
		}
	}
	for i := range rt.nameByField {
		if ft := rt.Type.FieldByIndex([]int{i}); ft.Type != reflect.TypeOf(string("")) {
			return rt, ErrInvalidFieldType
		}
	}
	return rt, nil
}
