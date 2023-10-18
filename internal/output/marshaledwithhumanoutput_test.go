package output_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

func TestMarshaledWithHumanOutput(t *testing.T) {
	marshaledWithHumanOutputTests := []outputTestCase{
		{
			name: "struct",
			input: output.MarshaledWithHumanOutput{
				Value: struct {
					KeyFromTag string `json:"key_from_tag"`
				}{"mock"},
				Output: output.None{},
			},
			expectedHumanResult: "",
			expectedJSONResult:  "{\n  \"key_from_tag\": \"mock\"\n}",
			expectedYAMLResult:  "key_from_tag: mock\n",
		},
	}
	for _, test := range marshaledWithHumanOutputTests {
		t.Run(test.name, test.Generate())
	}
}
