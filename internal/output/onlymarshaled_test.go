package output_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

func TestOnlyMarshaled(t *testing.T) {
	marshaledOnlyTests := []outputTestCase{
		{
			name:               "nil",
			input:              output.OnlyMarshaled{Value: nil},
			expectedJSONResult: "null",
			expectedYAMLResult: "null\n",
		},
		{
			name:               "string",
			input:              output.OnlyMarshaled{Value: "hello"},
			expectedJSONResult: `"hello"`,
			expectedYAMLResult: "hello\n",
		},
		{
			name:               "int",
			input:              output.OnlyMarshaled{Value: 123},
			expectedJSONResult: `123`,
			expectedYAMLResult: "123\n",
		},
		{
			name:               "struct",
			input:              output.OnlyMarshaled{Value: struct{ String string }{"mock"}},
			expectedJSONResult: "{\n  \"String\": \"mock\"\n}",
			expectedYAMLResult: "String: mock\n",
		},
	}

	for _, test := range marshaledOnlyTests {
		t.Run(test.name, test.Generate())
	}
}
