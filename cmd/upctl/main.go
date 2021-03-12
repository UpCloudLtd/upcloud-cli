package main

import (
	"fmt"
	"os"

	"github.com/UpCloudLtd/cli/internal/core"
)

// var (
// 	mainConfig = config.New(viper.New())
// 	mc         = commands.BuildCommand(
// 		&mainCommand{BaseCommand: commands.New("upctl", "UpCloud command line client")},
// 		nil, mainConfig)
// )

func main() {

	if err := core.BootstrapCLI(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

}
