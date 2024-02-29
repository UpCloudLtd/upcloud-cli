package format

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// Boolean returns val formatted as a boolean
func Boolean(val interface{}) (text.Colors, string, error) {
	if vb, ok := val.(bool); ok {
		if vb {
			return ui.DefaultBooleanColoursTrue, "yes", nil
		}

		return ui.DefaultBooleanColoursFalse, "no", nil
	} else if upb, ok := val.(upcloud.Boolean); ok {
		if upb == upcloud.True {
			return ui.DefaultBooleanColoursTrue, "yes", nil
		}

		return ui.DefaultBooleanColoursFalse, "no", nil
	}

	return nil, "", fmt.Errorf("cannot parse '%v' (%T) as boolean", val, val)
}

// Dereference returns "%v" Sprintf'ed value of a pointer
func Dereference[T any](val interface{}) (text.Colors, string, error) {
	ptr, ok := val.(*T)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse %T, expected pointer", val)
	}

	if ptr != nil {
		return nil, fmt.Sprintf("%v", *ptr), nil
	}

	return text.Colors{text.FgHiBlack}, "nil", nil
}

// PossiblyUnknownString outputs "Unknown" in light black if input value is an empty string, otherwise passesthrough the input value.
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

func usingColorFunction[T ~string](colorFunction func(T) text.Colors, val interface{}) (text.Colors, string, error) {
	typedVal, ok := val.(T)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse value from %T, expected %T", val, T(""))
	}

	if typedVal == "" {
		return PossiblyUnknownString(typedVal)
	}

	return nil, colorFunction(typedVal).Sprint(typedVal), nil
}
