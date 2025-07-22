package resolver

import (
	"path/filepath"
	"strings"
)

// MatchArgWithEqualFold checks if arg is a exact or case-insensitive match.
func MatchArgWithEqualFold(arg, value string) MatchType {
	if value == arg {
		return MatchTypeExact
	}
	if strings.EqualFold(value, arg) {
		return MatchTypeCaseInsensitive
	}
	return MatchTypeNone
}

// MatchArgWithGlobPattern checks if arg matches given value as a unix style glob pattern.
func MatchArgWithGlobPattern(arg, value string) MatchType {
	if matched, _ := filepath.Match(arg, value); matched {
		return MatchTypeGlobPattern
	}
	return MatchTypeNone
}

// MatchTitle checks if arg matches any of the given values by using MatchArgWithEqualFold and MatchArgWithGlobPattern matchers.
func MatchTitle(arg string, values ...string) MatchType {
	match := MatchTypeNone
	for _, value := range values {
		match = max(match, MatchArgWithEqualFold(arg, value))
		match = max(match, MatchArgWithGlobPattern(arg, value))
	}
	return match
}

func MatchUUID(arg, value string) MatchType {
	if value == arg {
		return MatchTypeExact
	}
	if strings.HasPrefix(value, arg) {
		return MatchTypePrefix
	}
	return MatchTypeNone
}
