package output_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failWriter struct{}

func (f failWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("MOCKERROR")
}

func TestRenderFailingWriter(t *testing.T) {
	fw := failWriter{}
	cfg := config.New()
	cfg.Viper().Set(config.KeyOutput, "human")
	err := output.Render(fw, cfg.Output(), output.None{})
	assert.EqualError(t, err, "MOCKERROR")
}

func TestRender(t *testing.T) {
	rr := strings.NewReader("raw hello")
	renderTests := []outputTestCase{
		{
			name:                "none",
			input:               output.None{},
			expectedHumanResult: "",
			expectedJSONResult:  "",
			expectedYAMLResult:  "",
		},
		{
			name:                "marshaled",
			input:               output.OnlyMarshaled{Value: "hello"},
			expectedHumanResult: "", // marshaled should not output in human mode
			expectedJSONResult: `"hello"
`,
			expectedYAMLResult: "hello\n",
		},
		{
			name:                "raw",
			input:               output.Raw{Source: io.NopCloser(rr)},
			expectedHumanResult: "raw hello",
			expectedJSONResult:  "raw hello",
			expectedYAMLResult:  "raw hello",
		},
	}
	for _, test := range renderTests {
		t.Run(test.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cfg := config.New()

			cfg.Viper().Set(config.KeyOutput, "human")
			err := output.Render(out, cfg.Output(), test.input)
			validateOutput(t, test.expectedHumanResult, test.expectedErrorMessage, out.Bytes(), err)
			_, err = rr.Seek(0, io.SeekStart)
			require.NoError(t, err)
			out.Truncate(0)

			cfg.Viper().Set(config.KeyOutput, "json")
			err = output.Render(out, cfg.Output(), test.input)
			validateOutput(t, test.expectedJSONResult, test.expectedErrorMessage, out.Bytes(), err)
			_, err = rr.Seek(0, io.SeekStart)
			require.NoError(t, err)
			out.Truncate(0)

			cfg.Viper().Set(config.KeyOutput, "yaml")
			err = output.Render(out, cfg.Output(), test.input)
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
