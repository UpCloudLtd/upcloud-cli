package output

import (
	"errors"
)

// MarshaledWithHumanDetails implements output.Command for a return value that is displayed as raw marshaled in JSON and YAML or as human details, such as UUID or IP address, eg. most create commands.
type MarshaledWithHumanDetails struct {
	Value   interface{}
	Details []DetailRow
}

// MarshalJSON implements json.Marshaler and output.Output
func (d MarshaledWithHumanDetails) MarshalJSON() ([]byte, error) {
	return OnlyMarshaled{Value: d.Value}.MarshalJSON()
}

// MarshalHuman implements output.Output
// For MarshaledWithHumanDetails outputs, we return *only* the details part in humanized output
func (d MarshaledWithHumanDetails) MarshalHuman() ([]byte, error) {
	return Details{Sections: []DetailSection{
		{Rows: d.Details},
	}}.MarshalHuman()
}

// MarshalRawMap implements output.Output
func (d MarshaledWithHumanDetails) MarshalRawMap() (map[string]interface{}, error) {
	return nil, errors.New("MarshaledWithHumanDetails output should not be used as part of multiple output, raw output is undefined")
}
