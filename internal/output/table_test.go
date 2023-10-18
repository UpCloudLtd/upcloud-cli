package output_test

import (
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

func TestTable(t *testing.T) {
	tableTests := []outputTestCase{
		{
			name: "simple",
			input: output.Table{
				Columns: []output.TableColumn{
					{Key: "aa", Header: "BB"},
					{Key: "cc", Header: "DD"},
				},
				Rows: []output.TableRow{
					{"1", "2"},
					{3, 4},
				},
			},
			expectedHumanResult: `
 BB   DD 
──── ────
 1    2  
 3    4  
`,
			expectedJSONResult: `[
  {
    "aa": "1",
    "cc": "2"
  },
  {
    "aa": 3,
    "cc": 4
  }
]`,
			expectedYAMLResult: `- aa: "1"
  cc: "2"
- aa: 3
  cc: 4
`,
		},
		{
			name: "formatting",
			input: output.Table{
				Columns: []output.TableColumn{
					{Key: "aa", Header: "BB"},
					{Key: "cc", Header: "DD"},
				},
				Rows: []output.TableRow{
					{4.12345677890123456778901234567789, time.Date(2001, 1, 1, 12, 0, 0, 0, time.UTC)},
				},
			},
			expectedHumanResult: `
                BB   DD                            
─────────────────── ───────────────────────────────
 4.123456778901234   2001-01-01 12:00:00 +0000 UTC 
`,
			expectedJSONResult: `[
  {
    "aa": 4.123456778901234,
    "cc": "2001-01-01T12:00:00Z"
  }
]`,
			expectedYAMLResult: `- aa: 4.123456778901234
  cc: "2001-01-01T12:00:00Z"
`,
		},
		{
			name: "no headers",
			input: output.Table{
				Columns: []output.TableColumn{
					{Key: "aa"},
					{Key: "cc"},
				},
				Rows: []output.TableRow{
					{"1", "2"},
				},
			},
			expectedHumanResult: `
 aa   cc 
──── ────
 1    2  
`,
			expectedJSONResult: `[
  {
    "aa": "1",
    "cc": "2"
  }
]`,
			expectedYAMLResult: `- aa: "1"
  cc: "2"
`,
		},
	}
	for _, test := range tableTests {
		t.Run(test.name, test.Generate())
	}
}
