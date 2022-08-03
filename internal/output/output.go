package output

import (
	"encoding/json"
)

// Output represent output from a command that can be rendered as JSON, YAML or human-readable
type Output interface {
	json.Marshaler
	MarshalHuman() ([]byte, error)
	MarshalRawMap() (map[string]interface{}, error)
}
