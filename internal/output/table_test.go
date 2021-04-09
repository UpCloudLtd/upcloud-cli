package output_test

import (
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tableTests = []struct {
	name                 string
	columns              []output.TableColumn
	rows                 []output.TableRow
	expectedHumanResult  string
	expectedJSONResult   string
	expectedYAMLResult   string
	expectedErrorMessage string
}{
	{
		name: "simple",
		columns: []output.TableColumn{
			{Key: "aa", Header: "BB"},
			{Key: "cc", Header: "DD"},
		},
		rows: []output.TableRow{
			{"1", "2"},
			{3, 4},
		},
		expectedHumanResult: ` BB   DD 
──── ────
 1    2  
 3    4  `,
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
}

func TestTable(t *testing.T) {
	for _, test := range tableTests {
		t.Run(test.name, func(t *testing.T) {
			input := output.Table{
				Columns: test.columns,
				Rows:    test.rows,
			}
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
