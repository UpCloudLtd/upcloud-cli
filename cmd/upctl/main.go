package main

import (
	"os"

	"github.com/UpCloudLtd/upcloud-cli/internal/core"
)

func main() {
	exitCode := core.Execute()
	os.Exit(exitCode)
}
