package root

import (
	"fmt"
	// "io"
	// "os"

	// "path"
	// "path/filepath"
	"runtime"
	// "strings"
	// "time"

	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"

	// "github.com/spf13/viper"
	// "gopkg.in/yaml.v3"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	// "github.com/UpCloudLtd/cli/internal/commands/all"
	// "github.com/UpCloudLtd/cli/internal/config"
	// "github.com/UpCloudLtd/cli/internal/terminal"
	// "github.com/UpCloudLtd/cli/internal/ui"
	// "github.com/UpCloudLtd/cli/internal/upapi"
	// "github.com/UpCloudLtd/cli/internal/validation"
)

func BuildAllCommands(conf *config.Config, rootCmd *cobra.Command) {
	nsCmds := GetNamespaceCommands(conf)
	for _, cmd := range nsCmds {
		c := commands.BuildCommand(conf, cmd)
		rootCmd.AddCommand(c)
	}
}

func GetNamespaceCommands(conf *config.Config) []*commands.BaseCommand {
	return commands.GenerateCmdList(
		VersionCmd(),
		ShellCompletionCmd(),
	)
}

func VersionCmd() *commands.BaseCommand {
	return &commands.BaseCommand{
		Name:  "version",
		Short: "Display software version",
		Long:  "Show the information about the software and build infos",
		Run: func(conf *config.Config, args interface{}) (i interface{}, e error) {
			return fmt.Sprintf(
				"Upctl %v\n\tBuild date: %v\n\tBuilt with: %v",
				config.Version, config.BuildDate, runtime.Version(),
			), nil
		},
	}
}

func ShellCompletionCmd() *commands.BaseCommand {
	return &commands.BaseCommand{}
}
