package config

import (
	"testing"
)

func TestIdentifierToEnvVarName(t *testing.T) {
	cases := [][2]string{
		[2]string{"HttpAddress", envVarPrefix + "_HTTP_ADDRESS"},
		[2]string{"SqlProtocol", envVarPrefix + "_SQL_PROTOCOL"},
		[2]string{"SqlAddress", envVarPrefix + "_SQL_ADDRESS"},
		[2]string{"HashDifficulty", envVarPrefix + "_HASH_DIFFICULTY"},
	}
	for _, c := range cases {
		if r := IdentifierToEnvVarName(c[0]); r != c[1] {
			t.Errorf("Testing with %s, should be %s, is %s.", c[0], c[1], r)
		}
	}
}
