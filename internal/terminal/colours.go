package terminal

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-isatty"
)

var forceColours *bool

func init() {
	if !Colours() {
		text.DisableColors()
	}
}

func ForceColours(v bool) {
	forceColours = &v
	text.EnableColors()
	if v == false {
		text.DisableColors()
	}
}

func Colours() bool {
	if forceColours != nil {
		return *forceColours
	}
	return isatty.IsTerminal(os.Stdout.Fd())
}
