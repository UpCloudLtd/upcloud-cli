package output_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

func TestDetails(t *testing.T) {
	detailTests := []outputTestCase{
		{
			name:                "no content",
			input:               output.Details{Sections: []output.DetailSection{}},
			expectedHumanResult: "\n",
			expectedJSONResult:  "{}",
			expectedYAMLResult:  "{}\n", // TODO: is this what we want?
		},
		{
			name: "single section",
			input: output.Details{Sections: []output.DetailSection{{Key: "foo", Title: "BAR", Rows: []output.DetailRow{
				{Key: "baz", Title: "boz", Value: "thisistest"},
			}}}},
			expectedHumanResult: "  \n  BAR\n    boz thisistest \n",
			expectedJSONResult:  "{\n  \"foo\": {\n    \"baz\": \"thisistest\"\n  }\n}",
			expectedYAMLResult:  "foo:\n    baz: thisistest\n",
		},
	}

	for _, test := range detailTests {
		t.Run(test.name, test.Generate())
	}
}
