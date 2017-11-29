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
		DefaultNames{"foo val", "bar val", "baz val"},
		DefaultNames{"foo \" quote", "bar \"quote", "baz \" quote"},
		DefaultNames{"1", "3", "2"},
		DefaultNames{",", "", "\nnewline\n"},
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
		TaggedNames{"foo val", "bar val", "baz val"},
		TaggedNames{"foo \" quote", "bar \"quote", "baz \" quote"},
		TaggedNames{"1", "3", "2"},
		TaggedNames{",", "", "\nnewline\n"},
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
		IgnoredNames{"foo val", "bar val", ""},
		IgnoredNames{"foo \" quote", "bar \"quote", ""},
		IgnoredNames{"1", "3", ""},
		IgnoredNames{",", "", ""},
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
		AnyNames{"foo val", map[string]string{"Bar": "bar val", "Baz": "baz val"}},
		AnyNames{"foo \" quote", map[string]string{"Bar": "bar \"quote", "Baz": "baz \" quote"}},
		AnyNames{"1", map[string]string{"Bar": "3", "Baz": "2"}},
		AnyNames{",", map[string]string{"Bar": "", "Baz": "\nnewline\n"}},
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
