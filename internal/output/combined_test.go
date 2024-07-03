package output_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

func TestCombined(t *testing.T) {
	combinedTests := []outputTestCase{
		{
			name:               "no content",
			input:              output.Combined{},
			expectedJSONResult: "{}",
			expectedYAMLResult: "{}\n", // TODO: is this what we want?
		},
		{
			name: "single section",
			input: output.Combined{
				{Key: "test", Title: "MOCK", Contents: output.Details{
					Sections: []output.DetailSection{{Key: "foo", Title: "BAR", Rows: []output.DetailRow{
						{Key: "baz", Title: "boz", Value: "thisistest"},
					}}},
				}},
			},
			expectedHumanResult: "  \n  BAR\n    boz thisistest \n",
			expectedJSONResult:  "{\n  \"test\": {\n    \"baz\": \"thisistest\"\n  }\n}",
			expectedYAMLResult:  "test:\n    baz: thisistest\n",
		},
		{
			name: "two sections",
			input: output.Combined{
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
			expectedHumanResult: "  \n  BAR\n    boz thisistest1 \n\n  \n  BAR\n    boz thisistest2 \n",
			expectedJSONResult:  "{\n  \"test1\": {\n    \"baz\": \"thisistest1\"\n  },\n  \"test2\": {\n    \"baz\": \"thisistest2\"\n  }\n}",
			expectedYAMLResult:  "test1:\n    baz: thisistest1\ntest2:\n    baz: thisistest2\n",
		},
		{
			name: "two tables",
			input: output.Combined{
				{Key: "test1", Title: "MOCK", Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "a", Header: "B"},
						{Key: "c", Header: "D"},
					},
					Rows: []output.TableRow{
						{1, "2a"},
						{"3", 4},
					},
				}},
				{Key: "test2", Title: "MOCK", Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "aa", Header: "BB"},
						{Key: "cc", Header: "DD"},
					},
					Rows: []output.TableRow{
						{"1a", 2},
						{3, "4"},
					},
				}},
			},
			expectedHumanResult: `
  MOCK

     B   D  
    ─── ────
     1   2a 
     3   4  
    
  MOCK

     BB   DD 
    ──── ────
     1a   2  
     3    4  
    `,
			expectedJSONResult: `{
  "test1": [
    {
      "a": 1,
      "c": "2a"
    },
    {
      "a": "3",
      "c": 4
    }
  ],
  "test2": [
    {
      "aa": "1a",
      "cc": 2
    },
    {
      "aa": 3,
      "cc": "4"
    }
  ]
}`,
			expectedYAMLResult: `test1:
    - a: 1
      c: 2a
    - a: "3"
      c: 4
test2:
    - aa: 1a
      cc: 2
    - aa: 3
      cc: "4"
`,
		},
	}
	for _, test := range combinedTests {
		t.Run(test.name, test.Generate())
	}
}
