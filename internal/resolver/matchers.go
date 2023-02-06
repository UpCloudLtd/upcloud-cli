package resolver

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
)

// MatchStringWithWhitespace checks if arg that may include whitespace matches given value. This checks both quoted args and auto-completed args handled with completion.RemoveWordBreaks.
func MatchArgWithWhitespace(arg string, value string) bool {
	return completion.RemoveWordBreaks(value) == arg || value == arg
}
