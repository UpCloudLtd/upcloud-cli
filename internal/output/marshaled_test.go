package output_test

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

var marshaledTests = []struct {
	name                 string
	marshalValue         interface{}
	expectedHumanResult  string
	expectedJSONResult   string
	expectedYAMLResult   string
	expectedErrorMessage string
}{
	{
		name:               "nil",
		marshalValue:       nil,
		expectedJSONResult: "null",
		expectedYAMLResult: "null\n",
	},
	{
		name:               "string",
		marshalValue:       "hello",
		expectedJSONResult: `"hello"`,
		expectedYAMLResult: "hello\n",
	},
	{
		name:               "int",
		marshalValue:       123,
		expectedJSONResult: `123`,
		expectedYAMLResult: "123\n",
	},
	{
		name:               "struct",
		marshalValue:       struct{ String string }{"mock"},
		expectedJSONResult: "{\n  \"String\": \"mock\"\n}",
		expectedYAMLResult: "string: mock\n",
	},
}

func TestMarshaled(t *testing.T) {
	for _, test := range marshaledTests {
		t.Run(test.name, func(t *testing.T) {
			input := output.OnlyMarshaled{Value: test.marshalValue}
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
