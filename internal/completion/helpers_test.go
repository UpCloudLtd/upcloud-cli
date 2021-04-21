package completion_test

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/stretchr/testify/assert"
	"testing"
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
			name:          "case insensitive",
			vals:          []string{"aba", "aBa", "Aba", "aab"},
			key:           "ab",
			caseSensitive: false,
			expected:      []string{"aba", "aBa", "Aba"},
		},
		{
			name:          "escaped output",
			vals:          []string{"a a ", "a(0)", "aab", "a;<!`'", "bbb"},
			key:           "a",
			caseSensitive: false,
			expected:      []string{"\"a a \"", "\"a(0)\"", "aab", "a;<!`'"},
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

func TestEscape(t *testing.T) {
	for _, test := range []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "no escape",
			in:       "asdasdasd",
			expected: "asdasdasd",
		},
		{
			name:     "escape spaces",
			in:       "asdas dasd",
			expected: "\"asdas dasd\"",
		},
		{
			name:     "escape open parentheses",
			in:       "asdas(",
			expected: "\"asdas(\"",
		},
		{
			name:     "escape closed parentheses",
			in:       "asdas()",
			expected: "\"asdas()\"",
		},
		{
			name:     "special chars not escaped",
			in:       "a;<!`'",
			expected: "a;<!`'",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, completion.Escape(test.in))
		})
	}
}
