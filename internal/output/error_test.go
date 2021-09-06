package output_test

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/output"
)

func TestError(t *testing.T) {
	t.Parallel()
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
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			test.Generate()(t)
		})
	}
}
