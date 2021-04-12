package output_test

import (
	"bytes"
	"errors"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

type failWriter struct{}

func (f failWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("MOCKERROR")
}

func TestRenderFailingWriter(t *testing.T) {
	fw := failWriter{}
	cfg := config.New()
	cfg.Top().Viper().Set(config.KeyOutput, "human")
	err := output.Render(fw, cfg, output.None{})
	assert.EqualError(t, err, "MOCKERROR")
}

func TestRender(t *testing.T) {
	var renderTests = []struct {
		name                 string
		output               output.Output
		expectedHumanResult  string
		expectedJSONResult   string
		expectedYAMLResult   string
		expectedErrorMessage string
	}{
		{
			name:                "none",
			output:              output.None{},
			expectedHumanResult: "\n",
			expectedJSONResult:  "\n",
			expectedYAMLResult:  "\n",
		},
		{
			name:                "marshaled",
			output:              output.Marshaled{Value: "hello"},
			expectedHumanResult: "\n", // marshaled should not output in human mode
			expectedJSONResult: `"hello"
`,
			expectedYAMLResult: "hello\n",
		},
	}
	for _, test := range renderTests {
		t.Run(test.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cfg := config.New()

			cfg.Top().Viper().Set(config.KeyOutput, "human")
			err := output.Render(out, cfg, test.output)
			validateOutput(t, test.expectedHumanResult, test.expectedErrorMessage, out.Bytes(), err)
			out.Truncate(0)

			cfg.Top().Viper().Set(config.KeyOutput, "json")
			err = output.Render(out, cfg, test.output)
			validateOutput(t, test.expectedJSONResult, test.expectedErrorMessage, out.Bytes(), err)
			out.Truncate(0)

			cfg.Top().Viper().Set(config.KeyOutput, "yaml")
			err = output.Render(out, cfg, test.output)
			validateOutput(t, test.expectedYAMLResult, test.expectedErrorMessage, out.Bytes(), err)
		})
	}
}

func validateOutput(t *testing.T, expectedOutput string, expectedErrorMessage string, output []byte, err error) {
	if expectedErrorMessage != "" {
		assert.EqualError(t, err, expectedErrorMessage)
		assert.Equal(t, 0, len(output))
	} else {
		assert.NoError(t, err)
		if expectedOutput == "" {
			// expecting an empty result, which could be nil or not so just check the length
			assert.Equal(t, 0, len(output))
		} else {
			assert.Equal(t, expectedOutput, string(output))
		}
	}

}
