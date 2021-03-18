package output

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

// Marshaled implements output.Command for a return value that is only displayed as raw marshaled in JSON and YAML
// eg. most 'state change' commands
type Marshaled struct {
	Value interface{}
}

// MarshalJSON implements json.Marshaler and output.Command
func (d Marshaled) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

// MarshalYAML implements output.Command, it marshals the value and returns the YAML as []byte
// nb. does *not* implement yaml.Marshaler
func (d Marshaled) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(d.Value)
}

// MarshalHuman implements output.Command
// For Marshaled outputs, we dont return anything as it's assumed the log output is what the user wants
func (d Marshaled) MarshalHuman() ([]byte, error) {
	return []byte{}, nil
}
