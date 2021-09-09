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
	"testing"
)

type DefaultNames struct {
	Foo string
	Bar string
	Baz string
}

func TestDefaultNames(t *testing.T) {
	buf := bytes.NewBufferString(`Foo,Baz,Bar
foo val,baz val,bar val
"foo "" quote","baz "" quote","bar ""quote"
1,2,3
",","
newline
",""
`)
	var dn DefaultNames
	var dne = []DefaultNames{
		{"foo val", "bar val", "baz val"},
		{"foo \" quote", "bar \"quote", "baz \" quote"},
		{"1", "3", "2"},
		{",", "", "\nnewline\n"},
	}
	d := NewDecoder(buf)
	d.ReadHeader()
	for _, e := range dne {
		err := d.ReadRecord(&dn)
		if err != nil {
			t.Errorf("ReadRecord should succeed. Failed with %s.", err.Error())
		}
		if e != dn {
			t.Errorf("Record structure should be %#v, is %#v instead.", e, dn)
		}
	}
}

type TaggedNames struct {
	One   string `csv:"Foo"`
	Two   string `csv:"Bar"`
	Three string `csv:"Baz"`
}

func TestTaggedNames(t *testing.T) {
	buf := bytes.NewBufferString(`Foo,Baz,Bar
foo val,baz val,bar val
"foo "" quote","baz "" quote","bar ""quote"
1,2,3
",","
newline
",""
`)
	var dn TaggedNames
	var dne = []TaggedNames{
		{"foo val", "bar val", "baz val"},
		{"foo \" quote", "bar \"quote", "baz \" quote"},
		{"1", "3", "2"},
		{",", "", "\nnewline\n"},
	}
	d := NewDecoder(buf)
	d.ReadHeader()
	for _, e := range dne {
		err := d.ReadRecord(&dn)
		if err != nil {
			t.Errorf("ReadRecord should succeed. Failed with %s.", err.Error())
		}
		if e != dn {
			t.Errorf("Record structure should be %#v, is %#v instead.", e, dn)
		}
	}
}

type IgnoredNames struct {
	Foo string
	Bar string
	Baz string `csv:"-"`
}

func TestIgnoredNames(t *testing.T) {
	buf := bytes.NewBufferString(`Foo,Baz,Bar
foo val,baz val,bar val
"foo "" quote","baz "" quote","bar ""quote"
1,2,3
",","
newline
",""
`)
	var dn IgnoredNames
	var dne = []IgnoredNames{
		{"foo val", "bar val", ""},
		{"foo \" quote", "bar \"quote", ""},
		{"1", "3", ""},
		{",", "", ""},
	}
	d := NewDecoder(buf)
	d.ReadHeader()
	for _, e := range dne {
		err := d.ReadRecord(&dn)
		if err != nil {
			t.Errorf("ReadRecord should succeed. Failed with %s.", err.Error())
		}
		if e != dn {
			t.Errorf("Record structure should be %#v, is %#v instead.", e, dn)
		}
	}
}

type AnyNames struct {
	Foo     string
	Default map[string]string `csv:",any"`
}

func TestAnyNames(t *testing.T) {
	buf := bytes.NewBufferString(`Foo,Baz,Bar
foo val,baz val,bar val
"foo "" quote","baz "" quote","bar ""quote"
1,2,3
",","
newline
",""
`)
	var dn AnyNames
	var dne = []AnyNames{
		{"foo val", map[string]string{"Bar": "bar val", "Baz": "baz val"}},
		{"foo \" quote", map[string]string{"Bar": "bar \"quote", "Baz": "baz \" quote"}},
		{"1", map[string]string{"Bar": "3", "Baz": "2"}},
		{",", map[string]string{"Bar": "", "Baz": "\nnewline\n"}},
	}
	d := NewDecoder(buf)
	d.ReadHeader()
	for _, e := range dne {
		err := d.ReadRecord(&dn)
		if err != nil {
			t.Errorf("ReadRecord should succeed. Failed with %s.", err.Error())
		}
		if e.Foo != dn.Foo || e.Default["Bar"] != dn.Default["Bar"] || e.Default["Baz"] != dn.Default["Baz"] {
			t.Errorf("Record structure should be %#v, is %#v instead.", e, dn)
		}
	}
}
