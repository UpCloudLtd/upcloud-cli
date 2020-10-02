package terminal

import (
	"os"

	"github.com/mattn/go-isatty"
)

var forceColours *bool

func ForceColours(v bool) {
	forceColours = &v
}

func Colours() bool {
	if forceColours != nil {
		return *forceColours
	}
	return isatty.IsTerminal(os.Stdout.Fd())
}
