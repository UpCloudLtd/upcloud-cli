package main

import (
	"os"

	"github.com/UpCloudLtd/up cloud-cli/v3/internal/core"
)

func main() {
	exitCode := core.Execute()
	os.Exit(exitCode)
}
