package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/clierrors"
)

const (
	formatHuman string = "human"
	formatJSON  string = "json"
	formatYAML  string = "yaml"
)

func formats() []string {
	return []string{
		formatHuman,
		formatJSON,
		formatYAML,
	}
}

// Render renders commandOutput to writer using cfg to configure the output.
func Render(writer io.Writer, outputFormat string, commandOutputs ...Output) (err error) {
	var b []byte
	switch outputFormat {
	case formatHuman:
		if b, err = toHuman(commandOutputs...); err == nil && len(b) != 0 {
			b = append(b, '\n')
		}
	case formatJSON:
		if b, err = toJSON(commandOutputs...); err == nil && len(b) != 0 {
			b = append(b, '\n')
		}
	case formatYAML:
		b, err = toYAML(commandOutputs...)
	default:
		err = fmt.Errorf("output format not valid: %s, valid formats: %v", outputFormat, formats())
	}
	if err != nil {
		return err
	}

	_, err = writer.Write(b)
	if err != nil {
		return err
	}

	// Render streaming outputs and count failed ones
	failedCount := 0
	for _, commandOutput := range commandOutputs {
		if rawOutput, ok := commandOutput.(Raw); ok {
			_, cErr := io.Copy(writer, rawOutput)
			err = errors.Join(err, cErr, rawOutput.Close())
		} else if _, ok := commandOutput.(Error); ok {
			failedCount++
		}
	}

	if err != nil {
		return err
	}
	if failedCount > 0 {
		return &clierrors.CommandFailedError{
			FailedCount: failedCount,
		}
	}
	return nil
}

func toHuman(commandOutputs ...Output) ([]byte, error) {
	buffer := new(bytes.Buffer)
	for _, commandOutput := range commandOutputs {
		outBytes, err := commandOutput.MarshalHuman()
		if err != nil {
			return nil, err
		}
		buffer.Write(outBytes)
	}

	return buffer.Bytes(), nil
}

func toJSON(commandOutputs ...Output) ([]byte, error) {
	var output []byte
	var jsonOutput []json.RawMessage
	for _, commandOutput := range commandOutputs {
		if _, ok := commandOutput.(None); !ok {
			// don't marshal none outputs
			outBytes, err := commandOutput.MarshalJSON()
			if err != nil {
				return nil, err
			}
			if len(outBytes) > 0 {
				jsonOutput = append(jsonOutput, outBytes)
			}
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

func toYAML(commandOutputs ...Output) ([]byte, error) {
	b, err := toJSON(commandOutputs...)
	if err != nil {
		return nil, err
	}

	return JSONToYAML(b)
}
