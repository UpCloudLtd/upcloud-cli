package resolver

import (
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

func MatchUUID(arg, value string) MatchType {
	if value == arg {
		return MatchTypeExact
	}
	if strings.HasPrefix(value, arg) {
		return MatchTypePrefix
	}
	return MatchTypeNone
}
