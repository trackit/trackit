package routes

import (
	"testing"
	"testing/quick"
)

func identityString(s string) string { return s }

func TestConstError(t *testing.T) {
	f := func(s string) string {
		return constError(s).Error()
	}
	if err := quick.CheckEqual(identityString, f, nil); err != nil {
		t.Error(err)
	}
}
