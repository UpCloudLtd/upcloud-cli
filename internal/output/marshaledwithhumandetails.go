package output

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
)

// MarshaledWithHumanDetails implements output.Command for a return value that is only displayed as raw marshaled in JSON and YAML
// eg. most 'state change' commands
type MarshaledWithHumanDetails struct {
	Value   interface{}
	Details []DetailRow
}

// MarshalJSON implements json.Marshaler and output.Output
func (d MarshaledWithHumanDetails) MarshalJSON() ([]byte, error) {
	if errValue, ok := d.Value.(error); ok {
		return json.MarshalIndent(map[string]interface{}{
			"error": errValue.Error(),
		}, "", "  ")
	}
	return json.MarshalIndent(d.Value, "", "  ")
}

// MarshalYAML implements output.Output, it marshals the value and returns the YAML as []byte
// nb. does *not* implement yaml.Marshaler
func (d MarshaledWithHumanDetails) MarshalYAML() ([]byte, error) {
	if errValue, ok := d.Value.(error); ok {
		return yaml.Marshal(map[string]interface{}{
			"error": errValue.Error(),
		})
	}
	return yaml.Marshal(d.Value)
}

// MarshalHuman implements output.Output
// For MarshaledWithHumanDetails outputs, we return *only* the details part in humanized output
func (d MarshaledWithHumanDetails) MarshalHuman() ([]byte, error) {
	return Details{Sections: []DetailSection{
		{Rows: d.Details},
	}}.MarshalHuman()
}

// MarshalRawMap implements output.Output
func (d MarshaledWithHumanDetails) MarshalRawMap() (map[string]interface{}, error) {
	return nil, errors.New("marshaled should not be used as part of multiple output, raw output is undefined")
}
