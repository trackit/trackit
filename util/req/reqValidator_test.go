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
	"testing"
	"testing/quick"
)

type testWithNoTags struct {
	Foo string
	Bar int
}

type testWithTags struct {
	Foo string `req:"nonzero" json:"foo"`
	Bar int    `req:"nonzero" json:"bar"`
}

func TestValidatorWithNoTags(t *testing.T) {
	f := func(e testWithNoTags) bool {
		vdr, err := CreateValidator(e)
		return err == nil && vdr == nil
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestValidatorWithTags(t *testing.T) {
	f := func(e testWithTags) bool {
		vdr, err := CreateValidator(e)
		return err == nil && vdr != nil
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestValidatorSuccess(t *testing.T) {
	vdr, err := CreateValidator(testWithTags{})
	if err != nil {
		t.Errorf("Creating validator: error should be nil, is %q instead.", err.Error())
	} else {
		err = vdr(testWithTags{
			Foo: "non-empty string",
			Bar: 9000,
		})
		if err != nil {
			t.Errorf("Validating structure: error should be nil, is %q instead.", err.Error())
		}
	}
}

func TestValidatorFailure(t *testing.T) {
	const expectedError = "foo: value is zero; bar: value is zero"
	vdr, err := CreateValidator(testWithTags{})
	if err != nil {
		t.Errorf("Creating validator: error should be nil, is %q instead.", err.Error())
	} else {
		err = vdr(testWithTags{
			Foo: "",
			Bar: 0,
		})
		if err == nil {
			t.Errorf("Validating structure: error should be %q, is nil instead.", expectedError)
		} else if err.Error() != expectedError {
			t.Errorf("Validating structure: error should be %q, is %q instead.", expectedError, err.Error())
		}
	}
}
