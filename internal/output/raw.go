package output

import "fmt"

// Raw is a way to output raw data to the user. It is *only* supported in humanized output and used for generating shell completion scripts.
type Raw []byte

// MarshalJSON implements output.Output.
func (s Raw) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("json output not supported")
}

// MarshalYAML implements output.Output.
func (s Raw) MarshalYAML() ([]byte, error) {
	return nil, fmt.Errorf("yaml output not supported")
}

// MarshalHuman implements output.Output.
func (s Raw) MarshalHuman() ([]byte, error) {
	return s, nil
}

// MarshalRawMap implements output.Output.
func (s Raw) MarshalRawMap() (map[string]interface{}, error) {
	return nil, fmt.Errorf("raw mao output not supported")
}
