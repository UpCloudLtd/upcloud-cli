package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
)

const truncatedSuffix = "..."

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

func FormatTime(tv time.Time) string {
	if tv.IsZero() {
		return ""
	}
	dur := time.Now().Sub(tv)
	if dur.Hours() < 8 {
		return fmt.Sprintf("%s (%dh%dm%ds ago)",
			tv.Format(time.RFC3339),
			dur/time.Hour,
			dur%time.Hour/time.Minute,
			dur%time.Minute/time.Second)
	}
	return tv.Format(time.RFC3339)
}
