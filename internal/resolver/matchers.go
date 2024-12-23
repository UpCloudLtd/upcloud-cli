package resolver

import (
	"path/filepath"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
)

// MatchStringWithWhitespace checks if arg that may include whitespace matches given value. This checks both quoted args and auto-completed args handled with completion.RemoveWordBreaks.
func MatchArgWithWhitespace(arg, value string) MatchType {
	if completion.RemoveWordBreaks(value) == arg || value == arg {
		return MatchTypeExact
	}
	if strings.EqualFold(completion.RemoveWordBreaks(value), arg) || strings.EqualFold(value, arg) {
		return MatchTypeCaseInsensitive
	}
	return MatchTypeNone
}

// MatchStringWithWhitespace checks if arg matches given value as an unix style glob pattern.
func MatchArgWithGlobPattern(arg, value string) MatchType {
	if matched, _ := filepath.Match(arg, value); matched {
		return MatchTypeGlobPattern
	}
	return MatchTypeNone
}

// MatchTitle checks if arg matches any of the given values by using MatchArgWithWhitespace and MatchArgWithGlobPattern matchers.
func MatchTitle(arg string, values ...string) MatchType {
	match := MatchTypeNone
	for _, value := range values {
		match = max(match, MatchArgWithWhitespace(arg, value))
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
