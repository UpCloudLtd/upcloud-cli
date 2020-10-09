package wrapper

import (
	"encoding/json"
)

// Wrap wraps an error that can satisfies Stringer and json.Marshaller
func Error(err error) error {
	return &errorWrapper{err: err}
}

type errorWrapper struct {
	err error
}

func (s *errorWrapper) String() string {
	return s.err.Error()
}

func (s *errorWrapper) MarshalJSON() ([]byte, error) {
	if v, ok := s.err.(json.Marshaler); ok {
		return v.MarshalJSON()
	}
	return json.Marshal(s.err.Error())
}

func (s *errorWrapper) Error() string {
	return s.err.Error()
}
