package output

import (
	"encoding/json"
	"errors"
	"fmt"
)

// OnlyMarshaled implements output.Command for a return value that is only displayed as raw marshaled in JSON and YAML
// eg. most 'state change' commands
type OnlyMarshaled struct {
	Value interface{}
}

// MarshalJSON implements json.Marshaler and output.Output
func (d OnlyMarshaled) MarshalJSON() ([]byte, error) {
	if errValue, ok := d.Value.(error); ok {
		return json.MarshalIndent(map[string]interface{}{
			"error": errValue.Error(),
		}, "", "  ")
	}
	return json.MarshalIndent(d.Value, "", "  ")
}

// MarshalHuman implements output.Output
// For OnlyMarshaled outputs, we dont return anything in humanized output as it's assumed the log output is what the user
// wants and it is down to the command itself to provide that.
func (d OnlyMarshaled) MarshalHuman() ([]byte, error) {
	if errValue, ok := d.Value.(error); ok {
		return []byte(fmt.Sprintf("ERROR: %v", errValue)), nil
	}
	return []byte{}, nil
}

// MarshalRawMap implements output.Output
func (d OnlyMarshaled) MarshalRawMap() (map[string]interface{}, error) {
	return nil, errors.New("onlymarshaled output should not be used as part of multiple output, raw output is undefined")
}
