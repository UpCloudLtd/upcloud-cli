package output

import "io"

// Raw is a way to output raw data from a source reader.
type Raw struct {
	Source io.ReadCloser
}

// MarshalJSON implements output.Output
func (s Raw) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

// MarshalHuman implements output.Output
func (s Raw) MarshalHuman() ([]byte, error) {
	return []byte{}, nil
}

// MarshalRawMap implements output.Output
func (s Raw) MarshalRawMap() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// Read implements io.ReadCloser.
func (s Raw) Read(p []byte) (n int, err error) {
	return s.Source.Read(p)
}

// Close implements io.ReadCloser.
func (s Raw) Close() error {
	return s.Source.Close()
}
