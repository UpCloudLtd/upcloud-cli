package output

import (
	"encoding/json"
)

// Error implements output.Command for a return value that is an error
type Error struct {
	Value    error
	Resolved string
	Original string
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
	// Human readable errors go to stderr. This output package is responsible for stdout output only.
	return []byte(""), nil
}

// MarshalRawMap implements output.Output
func (e Error) MarshalRawMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if e.Resolved != "" {
		m["resource"] = e.Resolved
	}
	if e.Original != "" {
		m["argument"] = e.Original
	}

	if e.Value == nil {
		m["error"] = "Unspecified error"
	} else {
		m["error"] = e.Value.Error()
	}

	return m, nil
}
