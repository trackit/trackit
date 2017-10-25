package aws

import (
	"math/rand"
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
