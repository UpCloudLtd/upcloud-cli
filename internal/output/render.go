package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/UpCloudLtd/cli/internal/config"
)

// Render renders commandOutput to writer using cfg to configure the output.
func Render(writer io.Writer, cfg *config.Config, commandOutputs ...Output) (err error) {
	var output []byte
	switch {
	case cfg.OutputHuman():
		var buffer = new(bytes.Buffer)
		for _, commandOutput := range commandOutputs {
			outBytes, err := commandOutput.MarshalHuman()
			if err != nil {
				return err
			}
			buffer.Write(outBytes)
		}
		// add a final newline to the end. all sections should print just the top newline for themselves.
		buffer.Write([]byte{'\n'})
		output = buffer.Bytes()
	case cfg.Top().IsSet(config.KeyOutput) && cfg.Output() == config.ValueOutputJSON:
		var jsonOutput []json.RawMessage
		for _, commandOutput := range commandOutputs {
			if _, ok := commandOutput.(None); !ok {
				// don't marshal none outputs
				outBytes, err := commandOutput.MarshalJSON()
				if err != nil {
					return err
				}
				jsonOutput = append(jsonOutput, outBytes)
			}
		}
		if len(jsonOutput) == 1 {
			output, err = json.MarshalIndent(jsonOutput[0], "", "  ")
		} else if len(jsonOutput) > 0 {
			output, err = json.MarshalIndent(jsonOutput, "", "  ")
		}
		if err != nil {
			return err
		}
		output = append(output, '\n')
	case cfg.Top().IsSet(config.KeyOutput) && cfg.Output() == config.ValueOutputYAML:
		var yamlOutput []string
		for _, commandOutput := range commandOutputs {
			if _, ok := commandOutput.(None); !ok {
				// don't marshal none outputs
				rawmap, err := commandOutput.MarshalYAML()
				if err != nil {
					return err
				}
				yamlOutput = append(yamlOutput, string(rawmap))
			}
		}
		switch {
		case len(yamlOutput) > 1:
			output = append([]byte("---\n"), []byte(strings.Join(yamlOutput, "---\n"))...)
			output = append(output, []byte("...\n")...)
		case len(yamlOutput) == 1:
			output = []byte(yamlOutput[0])
		default:
			output = []byte{'\n'}
		}
	default:
		return fmt.Errorf("not sure what to output")
	}
	if _, err := writer.Write(output); err != nil {
		return err
	}
	return nil
}
