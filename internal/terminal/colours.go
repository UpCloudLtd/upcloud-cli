package terminal

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-isatty"
)

var forceColours *bool

func init() {
	if !Colours() {
		// TODO: make color/colour consistent (from everywhere in the codebase)
		text.DisableColors()
	}
}

// ForceColours forces the color mode to match the value given in v
func ForceColours(v bool) {
	forceColours = &v
	text.EnableColors()
	if !v {
		text.DisableColors()
	}
}

// Colours returns true if the color mode is enabled
func Colours() bool {
	if forceColours != nil {
		return *forceColours
	}
	return isatty.IsTerminal(os.Stdout.Fd())
}
