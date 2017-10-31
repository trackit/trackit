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

package aws

import (
	"testing"
)

const (
	externalCount = 24
)

func TestGenerateExternal(t *testing.T) {
	var e [externalCount]string
	for i := range e {
		e[i] = generateExternal()
		t.Log(e[i])
		if len(e[i]) != externalLength {
			t.Errorf("Length should be %d, is %d.", externalLength, len(e[i]))
		}
		for j := range e[:i] {
			if e[i] == e[j] {
				t.Errorf("Externals should be unique.")
			}
		}
	}
}

func BenchmarkGenerateExternal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateExternal()
	}
}
