package output

import (
	"errors"

	"gopkg.in/yaml.v2"
)

// MarshaledWithHumanOutput implements output.Command for a return value that passes through raw marshaled JSON or YAML and only affects human output. Like MarshaledWithHumanDetails, but allows more complex output.
type MarshaledWithHumanOutput struct {
	Value  interface{}
	Output Output
}

// MarshalJSON implements json.Marshaler and output.Output
func (d MarshaledWithHumanOutput) MarshalJSON() ([]byte, error) {
	return OnlyMarshaled{Value: d.Value}.MarshalJSON()
}

// MarshalYAML implements output.Output, it marshals the value and returns the YAML as []byte
// nb. does *not* implement yaml.Marshaler
func (d MarshaledWithHumanOutput) MarshalYAML() ([]byte, error) {
	// Marshal to JSON and convert to YAML. This is to use key names from json field tags.
	jsonBytes, err := OnlyMarshaled{Value: d.Value}.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return jsonToYaml(jsonBytes)
}

// MarshalHuman implements output.Output
// For MarshaledWithHumanDetails outputs, we return *only* the details part in humanized output
func (d MarshaledWithHumanOutput) MarshalHuman() ([]byte, error) {
	return d.Output.MarshalHuman()
}

// MarshalRawMap implements output.Output
func (d MarshaledWithHumanOutput) MarshalRawMap() (map[string]interface{}, error) {
	return nil, errors.New("MarshaledWithHumanOutput output should not be used as part of multiple output, raw output is undefined")
}

func jsonToYaml(jsonIn []byte) ([]byte, error) {
	var yamlObj interface{}
	if err := yaml.Unmarshal(jsonIn, &yamlObj); err != nil {
		return nil, err
	}
	return yaml.Marshal(yamlObj)
}
