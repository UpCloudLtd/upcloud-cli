package core

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/terminal"
	"github.com/UpCloudLtd/cli/internal/ui"
)

// BuildRootCmd()
func BuildRootCmd(args []string, conf *config.Config) cobra.Command {
	rootCmd := cobra.Command{
		Use:   "upctl",
		Short: "UpCloud CLI",
		Long:  "upctl a CLI tool for managing your UpCloud services.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := valid.ValidateStruct(conf.GlobalFlags)
			if err != nil {
				return err
			}

			terminal.ForceColours(conf.GlobalFlags.Colors)

			// Load config
			if err := conf.InitConfig(); err != nil {
				return fmt.Errorf("Config load: %v", err)
			}

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {},
	}

	rootCmd.BashCompletionFunction = commands.CustomBashCompletionFunc(rootCmd.Use)

	flags := &pflag.FlagSet{}
	flags.StringVarP(
		&conf.GlobalFlags.ConfigFile, "config", "", "", "Config file",
	)
	flags.StringVarP(
		&conf.GlobalFlags.OutputFormat, "output", "o", "human",
		"Output format (supported: json, yaml and human)",
	)
	flags.BoolVar(
		&conf.GlobalFlags.Colors, "colours", true,
		"Use terminal colours (supported: auto, true, false)",
	)

	// Add flags
	flags.VisitAll(func(flag *pflag.Flag) {
		rootCmd.PersistentFlags().AddFlag(flag)
	})
	conf.ConfigBindFlagSet(flags)

	rootCmd.SetUsageTemplate(ui.CommandUsageTemplate())
	rootCmd.SetUsageFunc(ui.UsageFunc)

	return rootCmd
}
