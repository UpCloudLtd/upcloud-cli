package output

import (
	"encoding/json"
)

// Command represent output from a command that can be rendered as JSON, YAML or human-readable
type Command interface {
	json.Marshaler
	// yaml.Marshaler
	MarshalYAML() ([]byte, error)
	MarshalHuman() ([]byte, error)
}
