package terminal

import (
	"os"

	"github.com/mattn/go-isatty"
)

var (
	isStdoutTerminal bool
	isStderrTerminal bool
)

func init() {
	isStdoutTerminal = isatty.IsTerminal(os.Stdout.Fd())
	isStderrTerminal = isatty.IsTerminal(os.Stderr.Fd())
}

// IsStdoutTerminal returns true if the terminal is stdout
func IsStdoutTerminal() bool {
	return isStdoutTerminal
}

// IsStderrTerminal returns true if the terminal is stderr
func IsStderrTerminal() bool {
	return isStderrTerminal
}
