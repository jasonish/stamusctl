package compose

import (
	"testing"
)

type validateInputFlagTest struct {
	input  string
	errors bool
}

var validateInputFlagTests = []validateInputFlagTest{
	{"test", true},
	{"always", false},
}

func TestValidateInputFlag(t *testing.T) {

	for _, test := range validateInputFlagTests {
		input := Parameters{RestartMode: test.input, ElasticMemory: "1m", LogstashMemory: "1m"}

		output := ValidateInputFlag(input)
		if (output == nil) == test.errors {
			if test.errors {
				t.Fatalf(`ValidateInputFlag(Parameters{RestartMode: "%s"}), should get error`, test.input)
			} else {
				t.Fatalf(`ValidateInputFlag(Parameters{RestartMode: "%s"}), should not get error`, test.input)
			}
		}
	}
}
