package output_test

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

func TestMarshaledWithHumanDetails(t *testing.T) {
	marshaledWithHumanDetailsTests := []outputTestCase{
		{
			name:                "nil",
			input:               output.MarshaledWithHumanDetails{Value: nil, Details: nil},
			expectedHumanResult: "  \n  \n",
			expectedJSONResult:  "null",
			expectedYAMLResult:  "null\n",
		},
		{
			name: "string",
			input: output.MarshaledWithHumanDetails{
				Value: "hello",
				Details: []output.DetailRow{
					{Key: "mock", Title: "Fake", Value: "hiya"},
				},
			},
			expectedHumanResult: "  \n  Fake hiya \n",
			expectedJSONResult:  `"hello"`,
			expectedYAMLResult:  "hello\n",
		},
		{
			name: "int",
			input: output.MarshaledWithHumanDetails{
				Value: 123,
				Details: []output.DetailRow{
					{Key: "mock", Title: "Fake", Value: 456},
				},
			},
			expectedHumanResult: "  \n  Fake 456 \n",
			expectedJSONResult:  `123`,
			expectedYAMLResult:  "123\n",
		},
		{
			name: "error",
			input: output.MarshaledWithHumanDetails{
				Value: fmt.Errorf("MOCK"),
				Details: []output.DetailRow{
					{Key: "mock", Title: "Fake", Value: fmt.Errorf("HMOCK")},
				},
			},
			expectedHumanResult: "  \n  Fake HMOCK \n",
			expectedJSONResult:  "{\n  \"error\": \"MOCK\"\n}",
			expectedYAMLResult:  "error: MOCK\n",
		},
		{
			name: "struct",
			input: output.MarshaledWithHumanDetails{
				Value: struct{ String string }{"mock"},
				Details: []output.DetailRow{
					{Key: "mock", Title: "Fake", Value: struct{ String string }{"mock"}},
				},
			},
			expectedHumanResult: "  \n  Fake {mock} \n",
			expectedJSONResult:  "{\n  \"String\": \"mock\"\n}",
			expectedYAMLResult:  "String: mock\n",
		},
	}
	for _, test := range marshaledWithHumanDetailsTests {
		t.Run(test.name, test.Generate())
	}
}
