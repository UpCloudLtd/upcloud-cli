package core

import (
	"fmt"
	"os"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/all"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/terminal"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	valid "github.com/asaskevich/govalidator"
	"github.com/gemalto/flume"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// BuildRootCmd builds the root command
func BuildRootCmd(conf *config.Config) cobra.Command {
	rootCmd := cobra.Command{
		Use:   "upctl",
		Short: "UpCloud CLI",
		Long:  "upctl a CLI tool for managing your UpCloud services.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := valid.ValidateStruct(conf.GlobalFlags)
			if err != nil {
				return err
			}

			// detect desired colour output
			switch {
			case conf.GlobalFlags.ForceColours == config.True:
				text.EnableColors()
			case conf.GlobalFlags.NoColours == config.True:
				text.DisableColors()
			default:
				if terminal.IsStdoutTerminal() {
					text.EnableColors()
				} else {
					text.DisableColors()
				}
			}

			// Set up flume
			flume.SetOut(os.Stderr)
			// TODO: should we make the level configurable?
			logLvl := flume.DebugLevel
			if !conf.GlobalFlags.Debug {
				// not debugging, no log output!
				logLvl = flume.OffLevel
			}
			addCaller := true
			if err := flume.Configure(flume.Config{
				AddCaller:    &addCaller,
				DefaultLevel: logLvl,
				// do not colour logs, as it doesn't really fit our use case and complicates the colour handling
				Encoding: "ltsv",
			}); err != nil {
				return fmt.Errorf("flume config error: %w", err)
			}

			if err := conf.Load(); err != nil {
				return fmt.Errorf("cannot load configuration: %w", err)
			}

			// Validate viper output binding too
			if conf.Output() != config.ValueOutputHuman &&
				conf.Output() != config.ValueOutputJSON &&
				conf.Output() != config.ValueOutputYAML {
				return fmt.Errorf("output format '%v' not accepted", conf.Output())
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
	config.AddToggleFlag(flags, &conf.GlobalFlags.ForceColours, "force-colours", false, "force coloured output despite detected terminal support")
	config.AddToggleFlag(flags, &conf.GlobalFlags.NoColours, "no-colours", false, "disable coloured output despite detected terminal support")
	flags.BoolVar(
		&conf.GlobalFlags.Debug, "debug", false,
		"Print out more verbose debug logs",
	)
	flags.DurationVarP(
		&conf.GlobalFlags.ClientTimeout, "client-timeout", "t",
		0,
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

// BuildCLI generates the CLI tree and returns the rootCmd
func BuildCLI() cobra.Command {
	conf := config.New()
	rootCmd := BuildRootCmd(conf)

	all.BuildCommands(&rootCmd, conf)

	return rootCmd
}

// BootstrapCLI is the CLI entrypoint
func BootstrapCLI(args []string) error {
	rootCmd := BuildCLI()
	return rootCmd.Execute()
}
