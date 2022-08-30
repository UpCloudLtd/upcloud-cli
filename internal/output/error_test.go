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
		{
			name:                "error with resource",
			input:               output.Error{Value: fmt.Errorf("MOMO"), Resolved: "resourceName"},
			expectedHumanResult: "\nERROR: MOMO",
			expectedJSONResult:  "{\n  \"error\": \"MOMO\",\n  \"resource\": \"resourceName\"\n}",
			expectedYAMLResult:  "error: MOMO\nresource: resourceName\n",
		},
		{
			name:                "error with argument",
			input:               output.Error{Value: fmt.Errorf("MOMO"), Original: "testArg"},
			expectedHumanResult: "\nERROR: MOMO (testArg)",
			expectedJSONResult:  "{\n  \"argument\": \"testArg\",\n  \"error\": \"MOMO\"\n}",
			expectedYAMLResult:  "argument: testArg\nerror: MOMO\n",
		},
		{
			name:                "error with argument and resource",
			input:               output.Error{Value: fmt.Errorf("MOMO"), Resolved: "resourceName", Original: "testArg"},
			expectedHumanResult: "\nERROR: MOMO (testArg)",
			expectedJSONResult:  "{\n  \"argument\": \"testArg\",\n  \"error\": \"MOMO\",\n  \"resource\": \"resourceName\"\n}",
			expectedYAMLResult:  "argument: testArg\nerror: MOMO\nresource: resourceName\n",
		},
	}
	for _, test := range errorTests {
		t.Run(test.name, test.Generate())
	}
}
