package output

// None is used when the command returns no output
type None struct{}

// MarshalJSON implements output.Output
func (s None) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

// MarshalHuman implements output.Output
func (s None) MarshalHuman() ([]byte, error) {
	return []byte{}, nil
}

// MarshalRawMap implements output.Output
func (s None) MarshalRawMap() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
