package main

import (
	"os"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/core"
)

func main() {
	exitCode := core.Execute()
	os.Exit(exitCode)
}
