package core

import (
	"fmt"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/terminal"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
)

// BuildRootCmd builds the root command
func BuildRootCmd(_ []string, conf *config.Config) cobra.Command {
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

			if err := conf.Load(); err != nil {
				return fmt.Errorf("cannot load configuration: %w", err)
			}

			// Validate viper output binding too
			if conf.Output() != config.ValueOutputHuman &&
				conf.Output() != config.ValueOutputJSON &&
				conf.Output() != config.ValueOutputYAML {
				return fmt.Errorf("Output format '%v' not accepted", conf.Output())
			}

			return nil
		},
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
	flags.DurationVarP(
		&conf.GlobalFlags.ClientTimeout, "client-timeout", "t",
		time.Duration(60*time.Second),
		"CLI timeout when using interactive mode on some commands",
	)

	// XXX: Apply viper value to the help as default
	// Add flags
	flags.VisitAll(func(flag *pflag.Flag) {
		rootCmd.PersistentFlags().AddFlag(flag)
	})
	conf.ConfigBindFlagSet(flags)

	rootCmd.SetUsageTemplate(ui.CommandUsageTemplate())
	rootCmd.SetUsageFunc(ui.UsageFunc)

	return rootCmd
}
