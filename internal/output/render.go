package output

import (
	"io"

	"github.com/UpCloudLtd/cli/internal/config"
)

// Render renders commandOutput to writer using cfg to configure the output.
func Render(writer io.Writer, cfg *config.Config, commandOutput Output) (err error) {
	var bytes []byte
	switch {
	case commandOutput == nil:
		return nil
	case cfg.OutputHuman():
		bytes, err = commandOutput.MarshalHuman()
		if err != nil {
			return err
		}
	case cfg.Top().IsSet(config.KeyOutput) && cfg.Output() == config.ValueOutputJSON:
		bytes, err = commandOutput.MarshalJSON()
		if err != nil {
			return err
		}
	case cfg.Top().IsSet(config.KeyOutput) && cfg.Output() == config.ValueOutputYAML:
		bytes, err = commandOutput.MarshalYAML()
		if err != nil {
			return err
		}
	}
	if _, err := writer.Write(bytes); err != nil {
		return err
	}
	return nil
}
