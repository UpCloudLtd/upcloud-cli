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

func IsStdoutTerminal() bool {
	return isStdoutTerminal
}

func IsStderrTerminal() bool {
	return isStderrTerminal
}
