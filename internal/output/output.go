package output

import (
	"encoding/json"
)

// Output represent output from a command that can be rendered as JSON, YAML or human-readable
type Output interface {
	json.Marshaler
	// yaml.Marshaler
	MarshalYAML() ([]byte, error)
	MarshalHuman() ([]byte, error)
}
