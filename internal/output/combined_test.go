package output_test

import (
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

var combinedTests = []struct {
	name                 string
	sections             []output.CombinedSection
	expectedHumanResult  string
	expectedJSONResult   string
	expectedYAMLResult   string
	expectedErrorMessage string
}{
	{
		name:               "no content",
		sections:           []output.CombinedSection{},
		expectedJSONResult: "{}",
		expectedYAMLResult: "{}\n", // TODO: is this what we want?
	},
	{
		name: "single section",
		sections: []output.CombinedSection{
			{Key: "test", Title: "MOCK", Contents: output.Details{
				Sections: []output.DetailSection{{Key: "foo", Title: "BAR", Rows: []output.DetailRow{
					{Key: "baz", Title: "boz", Value: "thisistest"},
				}}},
			}},
		},
		expectedHumanResult: "  \n  BAR\n    boz thisistest \n\n",
		expectedJSONResult:  "{\n  \"test\": {\n    \"baz\": \"thisistest\"\n  }\n}",
		expectedYAMLResult:  "test:\n  baz: thisistest\n",
	},
	{
		name: "two sections",
		sections: []output.CombinedSection{
			{Key: "test1", Title: "MOCK", Contents: output.Details{
				Sections: []output.DetailSection{{Key: "foo", Title: "BAR", Rows: []output.DetailRow{
					{Key: "baz", Title: "boz", Value: "thisistest1"},
				}}},
			}},
			{Key: "test2", Title: "MOCK", Contents: output.Details{
				Sections: []output.DetailSection{{Key: "foo", Title: "BAR", Rows: []output.DetailRow{
					{Key: "baz", Title: "boz", Value: "thisistest2"},
				}}},
			}},
		},
		expectedHumanResult: "  \n  BAR\n    boz thisistest1 \n\n  \n  BAR\n    boz thisistest2 \n\n",
		expectedJSONResult:  "{\n  \"test1\": {\n    \"baz\": \"thisistest1\"\n  },\n  \"test2\": {\n    \"baz\": \"thisistest2\"\n  }\n}",
		expectedYAMLResult:  "test1:\n  baz: thisistest1\ntest2:\n  baz: thisistest2\n",
	},
}

func TestCombined(t *testing.T) {
	for _, test := range combinedTests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedErrorMessage == "" {
				bytes, err := output.Combined(test.sections).MarshalHuman()
				assert.NoError(t, err)
				assert.Equal(t, test.expectedHumanResult, string(bytes))
				bytes, err = output.Combined(test.sections).MarshalJSON()
				assert.NoError(t, err)
				assert.Equal(t, test.expectedJSONResult, string(bytes))
				bytes, err = output.Combined(test.sections).MarshalYAML()
				assert.NoError(t, err)
				assert.Equal(t, test.expectedYAMLResult, string(bytes))
			} else {
				bytes, err := output.Combined(test.sections).MarshalHuman()
				assert.EqualError(t, err, test.expectedErrorMessage)
				assert.Len(t, bytes, 0)
				bytes, err = output.Combined(test.sections).MarshalJSON()
				assert.EqualError(t, err, test.expectedErrorMessage)
				assert.Len(t, bytes, 0)
				bytes, err = output.Combined(test.sections).MarshalYAML()
				assert.EqualError(t, err, test.expectedErrorMessage)
				assert.Len(t, bytes, 0)
			}
		})
	}
}
