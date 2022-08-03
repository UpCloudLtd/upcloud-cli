package output

import (
	"encoding/json"
	"fmt"
)

// Error implements output.Command for a return value that is an error
type Error struct {
	Value error
}

// MarshalJSON implements json.Marshaler and output.Output
func (e Error) MarshalJSON() ([]byte, error) {
	marshaled, err := e.MarshalRawMap()
	if err != nil {
		return json.MarshalIndent(map[string]interface{}{
			"error": err.Error(),
		}, "", "  ")
	}

	return json.MarshalIndent(marshaled, "", "  ")
}

// MarshalHuman implements output.Output
func (e Error) MarshalHuman() ([]byte, error) {
	if e.Value == nil {
		return []byte("\nERROR: Unspecified error"), nil
	}

	return []byte(fmt.Sprintf("\nERROR: %v", e.Value.Error())), nil
}

// MarshalRawMap implements output.Output
func (e Error) MarshalRawMap() (map[string]interface{}, error) {
	if e.Value == nil {
		return map[string]interface{}{
			"error": "Unspecified error",
		}, nil
	}

	return map[string]interface{}{
		"error": e.Value.Error(),
	}, nil
}
