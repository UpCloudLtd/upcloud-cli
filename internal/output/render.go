package output

import (
	"github.com/UpCloudLtd/cli/internal/config"
	"io"
)

// Render renders commandOutput to writer using cfg to configure the output.
func Render(writer io.Writer, cfg *config.Config, commandOutput Command) (err error) {
	var bytes []byte
	switch {
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
	// TODO: remove this, a bit ugly but for some reason otherwise we dont get the tailing newline
	bytes = append(bytes, '\n')
	if _, err := writer.Write(bytes); err != nil {
		return err
	}
	return nil
}
