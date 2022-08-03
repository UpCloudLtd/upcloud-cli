package output

import (
	"errors"
)

// MarshaledWithHumanOutput implements output.Command for a return value that passes through raw marshaled JSON or YAML and only affects human output. Like MarshaledWithHumanDetails, but allows more complex output.
type MarshaledWithHumanOutput struct {
	Value  interface{}
	Output Output
}

// MarshalJSON implements json.Marshaler and output.Output
func (d MarshaledWithHumanOutput) MarshalJSON() ([]byte, error) {
	return OnlyMarshaled{Value: d.Value}.MarshalJSON()
}

// MarshalHuman implements output.Output
// For MarshaledWithHumanDetails outputs, we return *only* the details part in humanized output
func (d MarshaledWithHumanOutput) MarshalHuman() ([]byte, error) {
	return d.Output.MarshalHuman()
}

// MarshalRawMap implements output.Output
func (d MarshaledWithHumanOutput) MarshalRawMap() (map[string]interface{}, error) {
	return nil, errors.New("MarshaledWithHumanOutput output should not be used as part of multiple output, raw output is undefined")
}
