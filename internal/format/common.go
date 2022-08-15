package format

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
)

// PossiblyUnkonwnString outputs "Unknown" in light black if input value is an empty string, otherwise passesthrough the input value.
func PossiblyUnknownString(val interface{}) (text.Colors, string, error) {
	str, ok := val.(string)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse %T, expected string", val)
	}

	if str == "" {
		return nil, text.FgHiBlack.Sprint("unknown"), nil
	}
	return nil, str, nil
}
