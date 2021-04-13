package output_test

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

var detailTests = []struct {
	name                 string
	sections             []output.DetailSection
	expectedHumanResult  string
	expectedJSONResult   string
	expectedYAMLResult   string
	expectedErrorMessage string
}{
	{
		name:                "no content",
		sections:            []output.DetailSection{},
		expectedHumanResult: "\n",
		expectedJSONResult:  "{}",
		expectedYAMLResult:  "{}\n", // TODO: is this what we want?
	},
	{
		name: "single section",
		sections: []output.DetailSection{{Key: "foo", Title: "BAR", Rows: []output.DetailRow{
			{Key: "baz", Title: "boz", Value: "thisistest"},
		}}},
		expectedHumanResult: "  \n  BAR\n    boz thisistest \n",
		expectedJSONResult:  "{\n  \"foo\": {\n    \"baz\": \"thisistest\"\n  }\n}",
		expectedYAMLResult:  "foo:\n  baz: thisistest\n",
	},
}

func TestDetails(t *testing.T) {
	for _, test := range detailTests {
		t.Run(test.name, func(t *testing.T) {
			input := output.Details{Sections: test.sections}
			if test.expectedErrorMessage == "" {
				bytes, err := input.MarshalHuman()
				assert.NoError(t, err)
				assert.Equal(t, test.expectedHumanResult, string(bytes))
				bytes, err = input.MarshalJSON()
				assert.NoError(t, err)
				assert.Equal(t, test.expectedJSONResult, string(bytes))
				bytes, err = input.MarshalYAML()
				assert.NoError(t, err)
				assert.Equal(t, test.expectedYAMLResult, string(bytes))
			} else {
				bytes, err := input.MarshalHuman()
				assert.EqualError(t, err, test.expectedErrorMessage)
				assert.Len(t, bytes, 0)
				bytes, err = input.MarshalJSON()
				assert.EqualError(t, err, test.expectedErrorMessage)
				assert.Len(t, bytes, 0)
				bytes, err = input.MarshalYAML()
				assert.EqualError(t, err, test.expectedErrorMessage)
				assert.Len(t, bytes, 0)
			}
		})
	}
}
