package output

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

// Marshaled implements output.Command for a return value that is only displayed as raw marshaled in JSON and YAML
// eg. most 'state change' commands
type Marshaled struct {
	Value interface{}
}

// MarshalJSON implements json.Marshaler and output.Output
func (d Marshaled) MarshalJSON() ([]byte, error) {
	if errValue, ok := d.Value.(error); ok {
		return json.MarshalIndent(map[string]interface{}{
			"error": errValue.Error(),
		}, "", "  ")
	}
	return json.MarshalIndent(d.Value, "", "  ")
}

// MarshalYAML implements output.Output, it marshals the value and returns the YAML as []byte
// nb. does *not* implement yaml.Marshaler
func (d Marshaled) MarshalYAML() ([]byte, error) {
	if errValue, ok := d.Value.(error); ok {
		return yaml.Marshal(map[string]interface{}{
			"error": errValue.Error(),
		})
	}
	return yaml.Marshal(d.Value)
}

// MarshalHuman implements output.Output
// For Marshaled outputs, we dont return anything in humanized output as it's assumed the log output is what the user
// wants and it is down to the command itself to provide that.
func (d Marshaled) MarshalHuman() ([]byte, error) {
	if errValue, ok := d.Value.(error); ok {
		return []byte(fmt.Sprintf("ERROR: %v", errValue)), nil
	}
	return []byte{}, nil
}
