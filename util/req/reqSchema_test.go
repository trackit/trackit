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
