package output

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

// Error implements output.Command for a return value that is an error
type Error struct {
	Value error
}

// MarshalJSON implements json.Marshaler and output.Output
func (e Error) MarshalJSON() ([]byte, error) {
	marshaled, err := e.MarshalRawMap()
	if err != nil {
		return json.Marshal(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return json.Marshal(marshaled)

}

// MarshalYAML implements output.Output, it marshals the value and returns the YAML as []byte
// nb. does *not* implement yaml.Marshaler
func (e Error) MarshalYAML() ([]byte, error) {
	marshaled, err := e.MarshalRawMap()
	if err != nil {
		return yaml.Marshal(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return yaml.Marshal(marshaled)
}

// MarshalHuman implements output.Output
func (e Error) MarshalHuman() ([]byte, error) {
	return []byte(fmt.Sprintf("\nERROR: %v", e.Value.Error())), nil
}

// MarshalRawMap implements output.Output
func (e Error) MarshalRawMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"error": e.Value.Error(),
	}, nil
}
