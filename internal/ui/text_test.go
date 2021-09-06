package ui_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
)

func TestIndentText(t *testing.T) {
	t.Parallel()
	type args struct {
		s                      string
		prefix                 string
		repeatedPrefixAsSpaces bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "BaseCase",
			args: args{"foo\nbar", " ", false},
			want: " foo\n bar",
		},
		{
			name: "EmptyLinesIgnored",
			args: args{"foo\n\nbar\n", " ", false},
			want: " foo\n\n bar\n",
		},
		{
			name: "RepeatedPrefixAsSpaces",
			args: args{"foo\nbar", "err: ", true},
			want: "err: foo\n     bar",
		},
	}
	for _, tt := range tests {
		// grab local reference for parallel tests
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ui.IndentText(tt.args.s, tt.args.prefix, tt.args.repeatedPrefixAsSpaces); got != tt.want {
				t.Errorf("IndentText() = %q, want %q	", got, tt.want)
			}
		})
	}
}
