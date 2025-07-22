package completion_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"

	"github.com/stretchr/testify/assert"
)

func TestMatchStringPrefix(t *testing.T) {
	for _, test := range []struct {
		name          string
		vals          []string
		key           string
		caseSensitive bool
		expected      []string
	}{
		{
			name:          "empty",
			vals:          []string{},
			key:           "",
			caseSensitive: true,
			expected:      []string{},
		},
		{
			name:          "normal",
			vals:          []string{"aba", "bba", "cba"},
			key:           "ab",
			caseSensitive: true,
			expected:      []string{"aba"},
		},
		{
			name:          "capitalized",
			vals:          []string{"Capitalized", "capitalized"},
			key:           "Cap",
			caseSensitive: true,
			expected:      []string{"Capitalized"},
		},
		{
			name:          "doublequotedkey",
			vals:          []string{"aba", "bba", "cba"},
			key:           "\"ab\"",
			caseSensitive: true,
			expected:      []string{"aba"},
		},
		{
			name:          "singlequotedkey",
			vals:          []string{"aba", "bba", "cba"},
			key:           "'ab'",
			caseSensitive: true,
			expected:      []string{"aba"},
		},
		{
			name:          "case sensitive",
			vals:          []string{"aba", "aBa", "Aba"},
			key:           "ab",
			caseSensitive: true,
			expected:      []string{"aba"},
		},
		{
			name:          "case insensitive (lowercase key)",
			vals:          []string{"aba", "aBa", "Aba", "aab"},
			key:           "ab",
			caseSensitive: false,
			expected:      []string{"aba", "aBa", "Aba"},
		},
		{
			name:          "case insensitive (uppercase key)",
			vals:          []string{"aba", "aBa", "Aba", "aab"},
			key:           "AB",
			caseSensitive: false,
			expected:      []string{"aba", "aBa", "Aba"},
		},
		{
			name:          "output with special characters",
			vals:          []string{"a a ", "a(0)", "aab", "a;<!`'", "bbb"},
			key:           "a",
			caseSensitive: false,
			expected:      []string{"a a ", "a(0)", "aab", "a;<!`'"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := completion.MatchStringPrefix(test.vals, test.key, test.caseSensitive)
			assert.Len(t, result, len(test.expected))
			if len(test.expected) > 0 {
				// do not compare values unless there's more than one item, to avoid nil != []string{}
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
