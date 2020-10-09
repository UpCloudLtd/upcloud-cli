package ui

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

const truncatedSuffix = "..."

var (
	CommandUsageLineLength = 70
)

func TruncateText(s string, maxLen int) string {
	if text.RuneCount(s) > maxLen {
		s = text.Trim(s, maxLen-len(truncatedSuffix))
		s += truncatedSuffix
	}
	return s
}

func IndentText(s, prefix string, repeatedPrefixAsSpaces bool) string {
	b := []byte(s)
	r := make([]byte, 0, len(b))
	lb := true
	indented := false
	for _, c := range b {
		switch {
		case lb && c != '\n':
			if repeatedPrefixAsSpaces && indented {
				r = append(r, strings.Repeat(" ", len(prefix))...)
			} else {
				r = append(r, prefix...)
			}
			indented = true
			lb = false
		case c == '\n':
			lb = true
		}
		r = append(r, c)
	}
	return string(r)
}
