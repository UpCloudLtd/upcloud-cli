package ui

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	// DefaultHeaderColours defines the default colors used for headers
	DefaultHeaderColours = text.Colors{text.Bold}
	// DefaultUUUIDColours defines the default colors used for UUIDs
	DefaultUUUIDColours = text.Colors{text.FgHiBlue}
	// DefaultErrorColours defines the default colors used for errors
	DefaultErrorColours = text.Colors{text.FgHiRed, text.Bold}
	// DefaultAddressColours defines the default colors used for addresses
	DefaultAddressColours = text.Colors{text.FgHiMagenta}
	// DefaultBooleanColoursTrue defines the default colors used for boolean true values
	DefaultBooleanColoursTrue = text.Colors{text.FgHiGreen}
	// DefaultBooleanColoursFalse defines the default colors used for boolean false values
	DefaultBooleanColoursFalse = text.Colors{text.FgHiBlack}
	// DefaultNoteColours defines the default colors used for notes
	DefaultNoteColours = text.Colors{text.FgHiBlack}
)

// FormatBool return v formatted (eg. colorized)
func FormatBool(v bool) string {
	if v {
		return DefaultBooleanColoursTrue.Sprint("yes")
	}
	return DefaultBooleanColoursFalse.Sprint("no")
}

// FormatRange takes start and end value and generates a ranged value
func FormatRange(start, end string) string {
	if start == end {
		if start == "" {
			return "Any"
		}

		return start
	}

	if end == "" {
		return start
	}

	return fmt.Sprintf("%s →\n%s", start, end)
}

// ConcatStrings like join but handles well the empty strings
func ConcatStrings(strs ...string) string {
	ret := fmt.Sprintf(strs[0])

	if len(strs) <= 1 {
		return ret
	}

	for _, str := range strs[1:] {
		if str != "" {
			ret = fmt.Sprintf("%s/%s", ret, str)
		}
	}

	return ret
}
