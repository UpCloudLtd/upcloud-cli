package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
)

const truncatedSuffix = "..."

// TruncateText truncates s to maxLen and adds a suffix ("...") if the string was truncated.
// nb. maxLen will include the suffix so maxLen is respected in all cases
func TruncateText(s string, maxLen int) string {
	if text.RuneWidthWithoutEscSequences(s) > maxLen {
		s = text.Trim(s, maxLen-len(truncatedSuffix))
		s += truncatedSuffix
	}
	return s
}

// IndentText indents s with prefix. if repeatedPrefixAsSpaces is true, only the first indented line will get the prefix.
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

// FormatTime returns a time.Time in RFC3339, adding information on how long ago it was if the time was under 8 hours ago.
func FormatTime(tv time.Time) string {
	if tv.IsZero() {
		return ""
	}
	dur := time.Since(tv)
	if dur.Hours() < 8 {
		return fmt.Sprintf("%s (%dh%dm%ds ago)",
			tv.Format(time.RFC3339),
			dur/time.Hour,
			dur%time.Hour/time.Minute,
			dur%time.Minute/time.Second)
	}
	return tv.Format(time.RFC3339)
}
