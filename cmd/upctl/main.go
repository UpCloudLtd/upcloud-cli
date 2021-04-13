package main

import (
	"os"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands/all"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/core"
)

// BootstrapCLI is the CLI entrypoint
func BootstrapCLI(args []string) error {

	conf := config.New()
	rootCmd := core.BuildRootCmd(args, conf)

	all.BuildCommands(&rootCmd, conf)
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func main() {

	if err := BootstrapCLI(os.Args); err != nil {
		os.Exit(1)
	}
}
