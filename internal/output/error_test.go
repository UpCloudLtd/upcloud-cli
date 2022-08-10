package output_test

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/output"
)

func TestError(t *testing.T) {
	errorTests := []outputTestCase{
		{
			name:                "nil",
			input:               output.Error{Value: nil},
			expectedHumanResult: "\nERROR: Unspecified error",
			expectedJSONResult:  "{\n  \"error\": \"Unspecified error\"\n}",
			expectedYAMLResult:  "error: Unspecified error\n",
		},
		{
			name:                "error",
			input:               output.Error{Value: fmt.Errorf("MOMO")},
			expectedHumanResult: "\nERROR: MOMO",
			expectedJSONResult:  "{\n  \"error\": \"MOMO\"\n}",
			expectedYAMLResult:  "error: MOMO\n",
		},
	}
	for _, test := range errorTests {
		t.Run(test.name, test.Generate())
	}
}
