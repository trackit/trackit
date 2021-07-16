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
	"bytes"
	"reflect"
	"testing"
)

type DefaultNames struct {
	Foo string
	Bar string
	Baz string
}

type TaggedNames struct {
	One   string `csv:"Foo"`
	Two   string `csv:"Bar"`
	Three string `csv:"Baz"`
}

type IgnoredNames struct {
	Foo string
	Bar string
	Baz string `csv:"-"`
}

type AnyNames struct {
	Foo     string
	Default map[string]string `csv:",any"`
}

var buf = `Foo,Baz,Bar
foo val,baz val,bar val
"foo "" quote","baz "" quote","bar ""quote"
1,2,3
",","
newline
",""
`

func testExecutor(s interface{}, t *testing.T, testMessage func(interface{}, interface{}) bool) {
	array := reflect.ValueOf(s)
	if reflect.TypeOf(s).Kind() != reflect.Slice {
		t.Error("Argument are not an array")
	}
	r := reflect.New(reflect.TypeOf(array.Index(0).Interface()))
	d := NewDecoder(bytes.NewBufferString(buf))
	d.ReadHeader()
	for i := 0; i < array.Len(); i++ {
		v := r.Interface()
		err := d.ReadRecord(v)
		e := array.Index(i)
		if err != nil {
			t.Errorf("ReadRecord should succeed. Failed with %s.", err.Error())
		}
		if !testMessage(e, v) {
			t.Errorf("Record structure should be %#v, is %#v instead.", e, v)
		}
	}
}

func TestDefaultNames(t *testing.T) {
	var dne = []DefaultNames{
		{"foo val", "bar val", "baz val"},
		{"foo \" quote", "bar \"quote", "baz \" quote"},
		{"1", "3", "2"},
		{",", "", "\nnewline\n"},
	}
	testExecutor(dne, t, func(e, v interface{}) bool { return e != v })
}

func TestTaggedNames(t *testing.T) {
	var dne = []TaggedNames{
		{"foo val", "bar val", "baz val"},
		{"foo \" quote", "bar \"quote", "baz \" quote"},
		{"1", "3", "2"},
		{",", "", "\nnewline\n"},
	}
	testExecutor(dne, t, func(e, v interface{}) bool { return e != v })
}

func TestIgnoredNames(t *testing.T) {
	var dne = []IgnoredNames{
		{"foo val", "bar val", ""},
		{"foo \" quote", "bar \"quote", ""},
		{"1", "3", ""},
		{",", "", ""},
	}
	testExecutor(dne, t, func(e, v interface{}) bool { return e != v })
}

func TestAnyNames(t *testing.T) {
	var dne = []AnyNames{
		{"foo val", map[string]string{"Bar": "bar val", "Baz": "baz val"}},
		{"foo \" quote", map[string]string{"Bar": "bar \"quote", "Baz": "baz \" quote"}},
		{"1", map[string]string{"Bar": "3", "Baz": "2"}},
		{",", map[string]string{"Bar": "", "Baz": "\nnewline\n"}},
	}
	testExecutor(dne, t, func(e, v interface{}) bool {
		return e.(AnyNames).Foo != v.(AnyNames).Foo || e.(AnyNames).Default["Bar"] != v.(AnyNames).Default["Bar"] || e.(AnyNames).Default["Baz"] != v.(AnyNames).Default["Baz"]
	})
}
