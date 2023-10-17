package testutils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertOutputHasType(t *testing.T, output string, expected interface{}) {
	t.Helper()

	err := json.Unmarshal([]byte(output), expected)
	if typeError, ok := err.(*json.UnmarshalTypeError); ok {
		assert.Nil(t, err, "Expected %t, got %s: %s", expected, typeError.Value, output)
	} else {
		assert.Nil(t, err)
	}
}

func AssertOutputIsList(t *testing.T, output string) {
	t.Helper()

	list := make([]interface{}, 0)
	AssertOutputHasType(t, output, &list)
}
