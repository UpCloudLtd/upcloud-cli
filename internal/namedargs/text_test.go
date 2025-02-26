package namedargs_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/stretchr/testify/assert"
)

func TestValidValuesHelp(t *testing.T) {
	for _, test := range []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "no values",
			values:   []string{},
			expected: "",
		},
		{
			name:     "single values",
			values:   []string{"foo"},
			expected: "`foo`",
		},
		{
			name:     "two values",
			values:   []string{"foo", "bar"},
			expected: "`foo` and `bar`",
		},
		{
			name:     "multiple values",
			values:   []string{"foo", "bar", "baz"},
			expected: "`foo`, `bar` and `baz`",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual := namedargs.ValidValuesHelp(test.values...)
			assert.Equal(t, test.expected, actual)
		})
	}
}
