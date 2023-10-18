package output_test

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

type outputTestCase struct {
	name                 string
	input                output.Output
	expectedHumanResult  string
	expectedJSONResult   string
	expectedYAMLResult   string
	expectedErrorMessage string
}

func (c outputTestCase) Generate() func(t *testing.T) {
	return func(t *testing.T) {
		text.DisableColors()
		if c.expectedErrorMessage == "" {
			bytes, err := c.input.MarshalHuman()
			assert.NoError(t, err)
			assert.Equal(t, c.expectedHumanResult, string(bytes))
			bytes, err = c.input.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, c.expectedJSONResult, string(bytes))
			bytes, err = output.JSONToYAML(bytes)
			assert.NoError(t, err)
			assert.Equal(t, c.expectedYAMLResult, string(bytes))
		} else {
			bytes, err := c.input.MarshalHuman()
			assert.EqualError(t, err, c.expectedErrorMessage)
			assert.Len(t, bytes, 0)
			bytes, err = c.input.MarshalJSON()
			assert.EqualError(t, err, c.expectedErrorMessage)
			assert.Len(t, bytes, 0)
			bytes, err = output.JSONToYAML(bytes)
			assert.EqualError(t, err, c.expectedErrorMessage)
			assert.Len(t, bytes, 0)
		}
	}
}
