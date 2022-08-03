package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/UpCloudLtd/upcloud-cli/internal/clierrors"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
)

func renderJSON(cfg *config.Config, commandOutputs ...Output) ([]byte, error) {
	var output []byte
	var jsonOutput []json.RawMessage
	for _, commandOutput := range commandOutputs {
		if _, ok := commandOutput.(None); !ok {
			// don't marshal none outputs
			outBytes, err := commandOutput.MarshalJSON()
			if err != nil {
				return nil, err
			}
			jsonOutput = append(jsonOutput, outBytes)
		}
	}

	var err error
	if len(jsonOutput) == 1 {
		output, err = json.MarshalIndent(jsonOutput[0], "", "  ")
	} else if len(jsonOutput) > 0 {
		output, err = json.MarshalIndent(jsonOutput, "", "  ")
	}
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Render renders commandOutput to writer using cfg to configure the output.
func Render(writer io.Writer, cfg *config.Config, commandOutputs ...Output) (err error) {
	var output []byte
	switch {
	case cfg.OutputHuman():
		buffer := new(bytes.Buffer)
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
	case cfg.IsSet(config.KeyOutput) && cfg.Output() == config.ValueOutputJSON:
		output, err = renderJSON(cfg, commandOutputs...)
		if err != nil {
			return err
		}
		output = append(output, '\n')
	case cfg.IsSet(config.KeyOutput) && cfg.Output() == config.ValueOutputYAML:
		jsonOutput, err := renderJSON(cfg, commandOutputs...)
		if err != nil {
			return err
		}

		output, err = JSONToYAML(jsonOutput)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("not sure what to output")
	}
	if _, err := writer.Write(output); err != nil {
		return err
	}

	// Count failed outputs
	failedCount := 0
	for _, commandOutput := range commandOutputs {
		if _, ok := commandOutput.(Error); ok {
			failedCount++
		}
	}

	if failedCount > 0 {
		return &clierrors.CommandFailedError{
			FailedCount: failedCount,
		}
	}
	return nil
}
