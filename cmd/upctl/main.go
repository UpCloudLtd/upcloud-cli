package main

import (
	"os"

	"github.com/UpCloudLtd/upcloud-cli/internal/core"
)

func main() {
	if err := core.BootstrapCLI(os.Args); err != nil {
		os.Exit(1)
	}
}
