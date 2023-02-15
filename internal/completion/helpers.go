package completion

import (
	"regexp"
	"strings"
)

var oneOrMoreSpace = regexp.MustCompile(` +`)

// RemoveWordBreaks replaces all whitespaces in input strings with non-breaking spaces to prevent bash from splitting completion with whitespace into multiple completions.
//
// This hack allows us to use cobras built-in completion logic and can be removed once cobra supports whitespace in bash completions (See https://github.com/spf13/cobra/issues/1740).
func RemoveWordBreaks(input string) string {
	return oneOrMoreSpace.ReplaceAllString(input, "\u00A0")
}

// MatchStringPrefix returns a list of string in vals which have a prefix as specified in key. Quotes are removed from key and output strings are escaped according to completion rules
func MatchStringPrefix(vals []string, key string, caseSensitive bool) []string {
	var r []string
	key = strings.Trim(key, "'\"")
	for _, v := range vals {
		if (caseSensitive && strings.HasPrefix(v, key)) ||
			(!caseSensitive && strings.HasPrefix(strings.ToLower(v), strings.ToLower(key))) ||
			key == "" {
			r = append(r, RemoveWordBreaks(v))
		}
	}
	return r
}
