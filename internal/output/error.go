package output

import (
	"encoding/json"
	"fmt"
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
	if e.Value == nil {
		return []byte("\nERROR: Unspecified error"), nil
	}
	s := fmt.Sprintf("\nERROR: %v", e.Value.Error())
	if e.Original != "" {
		s += fmt.Sprintf(" (%s)", e.Original)
	}

	return []byte(s), nil
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
