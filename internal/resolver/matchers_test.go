package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchers(t *testing.T) {
	cases := []struct {
		name     string
		execFn   func(string, ...string) MatchType
		arg      string
		value    string
		expected MatchType
	}{
		{
			name:     "Exact match",
			execFn:   MatchTitle,
			arg:      "McDuck",
			value:    "McDuck",
			expected: MatchTypeExact,
		},
		{
			name:     "Case-insensitive match",
			execFn:   MatchTitle,
			arg:      "mcduck",
			value:    "McDuck",
			expected: MatchTypeCaseInsensitive,
		},
		{
			name:     "Glob match",
			execFn:   MatchTitle,
			arg:      "McDuck-*",
			value:    "McDuck-1",
			expected: MatchTypeGlobPattern,
		},
		{
			name:     "No case-insensitive glob match",
			execFn:   MatchTitle,
			arg:      "mcduck-*",
			value:    "McDuck-1",
			expected: MatchTypeNone,
		},
		{
			name:     "No match",
			execFn:   MatchTitle,
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
