package completion

import (
	"fmt"
	"strings"
)

// MatchStringPrefix returns a list of string in vals which have a prefix as specified in key. Quotes are removed from key and output strings are escaped according to completion rules
func MatchStringPrefix(vals []string, key string, caseSensitive bool) []string {
	var r []string
	key = strings.Trim(key, "'\"")
	if caseSensitive {
		key = strings.ToLower(key)
	}
	for _, v := range vals {
		if (caseSensitive && strings.HasPrefix(v, key)) ||
			(!caseSensitive && strings.HasPrefix(strings.ToLower(v), key)) ||
			key == "" {
			r = append(r, Escape(v))
		}
	}
	return r
}

// Escape escapes a string according to completion rules (?)
// in effect, this means that the string will be quoted with double quotes if it contains a space or parentheses.
func Escape(s string) string {
	if strings.ContainsAny(s, ` ()`) {
		return fmt.Sprintf(`"%s"`, s)
	}
	return s
}
