package ui

import "testing"

func TestIndentText(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := IndentText(tt.args.s, tt.args.prefix, tt.args.repeatedPrefixAsSpaces); got != tt.want {
				t.Errorf("IndentText() = %q, want %q	", got, tt.want)
			}
		})
	}
}
