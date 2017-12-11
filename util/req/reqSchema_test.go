package req

import (
	"bytes"
	"reflect"
	"testing"
)

type testWithSomeTags struct {
	Foo string `req:"nonzero" json:"foo"`
	Bar int    `              json:"bar"`
}

const testWithTagsExpectedSchema = `{
	"foo": ¡string!,
	"bar": ¿int?
}`

func TestSchemaWithTags(t *testing.T) {
	typ := reflect.TypeOf(testWithSomeTags{})
	buf := bytes.NewBuffer(make([]byte, 2048))
	buf.Reset()
	err := GetSchema(buf, typ)
	if err != nil {
		t.Error(err)
	}
	if sch := buf.String(); sch != testWithTagsExpectedSchema {
		t.Errorf("Schema should be %#q, is %#q instead.", testWithTagsExpectedSchema, sch)
	}
}
