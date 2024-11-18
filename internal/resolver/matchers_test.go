package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchers(t *testing.T) {
	cases := []struct {
		name     string
		execFn   func(string, string) MatchType
		arg      string
		value    string
		expected MatchType
	}{
		{
			name:     "Exact match",
			execFn:   MatchArgWithWhitespace,
			arg:      "McDuck",
			value:    "McDuck",
			expected: MatchTypeExact,
		},
		{
			name:     "Case-insensitive match",
			execFn:   MatchArgWithWhitespace,
			arg:      "mcduck",
			value:    "McDuck",
			expected: MatchTypeCaseInsensitive,
		},
		{
			name:     "No match",
			execFn:   MatchArgWithWhitespace,
			arg:      "scrooge",
			value:    "McDuck",
			expected: MatchTypeNone,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.execFn(tt.arg, tt.value))
		})
	}
}
