package ui

import "github.com/jedib0t/go-pretty/v6/text"

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
