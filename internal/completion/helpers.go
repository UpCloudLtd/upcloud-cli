package completion

import (
	"strings"
)

// MatchStringPrefix returns a list of string in vals which have a prefix as specified in key. Quotes are removed from key and output strings are escaped according to completion rules
func MatchStringPrefix(vals []string, key string, caseSensitive bool) []string {
	var r []string
	key = strings.Trim(key, "'\"")
	for _, v := range vals {
		if (caseSensitive && strings.HasPrefix(v, key)) ||
			(!caseSensitive && strings.HasPrefix(strings.ToLower(v), strings.ToLower(key))) ||
			key == "" {
			r = append(r, v)
		}
	}
	return r
}
