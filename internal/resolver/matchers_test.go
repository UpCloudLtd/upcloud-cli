package resolver

import "testing"

func TestMatchers(t *testing.T) {
	cases := []struct {
		name     string
		execFn   func(string, string) bool
		arg      string
		value    string
		expected bool
	}{
		{
			name:     "Matcher no case",
			execFn:   MatchArgWithWhitespace,
			arg:      "test",
			value:    "test",
			expected: true,
		},
		{
			name:     "Matcher with case",
			execFn:   MatchArgWithWhitespace,
			arg:      "TeSt",
			value:    "test",
			expected: true,
		},
		{
			name:     "Matcher invalid",
			execFn:   MatchArgWithWhitespace,
			arg:      "test",
			value:    "invalid",
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.execFn(tt.arg, tt.value); result != tt.expected {
				t.Errorf("Matcher() failed %v, wanted %v", result, tt.expected)
			}
		})
	}
}
