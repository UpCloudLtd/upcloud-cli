package ui

import "github.com/jedib0t/go-pretty/v6/text"

var (
	DefaultHeaderColours       = text.Colors{text.Bold}
	DefaultUuidColours         = text.Colors{text.FgHiBlue}
	DefaultErrorColours        = text.Colors{text.FgHiRed, text.Bold}
	DefaultAddressColours      = text.Colors{text.FgHiMagenta}
	DefaultBooleanColoursTrue  = text.Colors{text.FgHiGreen}
	DefaultBooleanColoursFalse = text.Colors{text.FgHiBlack}
	DefaultNoteColours         = text.Colors{text.FgHiBlack}
)

func FormatBool(v bool) string {
	if v {
		return DefaultBooleanColoursTrue.Sprint("yes")
	}
	return DefaultBooleanColoursFalse.Sprint("no")
}
