package core

import (
	"fmt"
	"os"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/clierrors"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/base"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/terminal"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	valid "github.com/asaskevich/govalidator"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const rootCmdLongDescription = `UpCloud command-line interface

` + "`" + `upctl` + "`" + ` provides a command-line interface to UpCloud services. It allows you to control your resources from the command line or any compatible interface.

To be able to manage your UpCloud resources, you need to configure credentials for ` + "`" + `upctl` + "`" + `. We recommend using API tokens. Configure credentials by using ` + "`" + `upctl account login --with-token` + "`" + ` command to save a token to system keyring, in ` + "`" + `~/.config/upctl.yaml` + "`" + ` or set the ` + "`" + `UPCLOUD_TOKEN` + "`" + ` environment variable. Alternatively, use ` + "`" + `UPCLOUD_USERNAME` + "`" + ` and ` + "`" + `UPCLOUD_PASSWORD` + "`" + ` environment variables. When using username and password, API access can be configured on the Account page of the UpCloud Hub. We recommend you to set up a sub-account specifically for the API usage with its own credentials, as it allows you to assign specific permissions for increased security.`

// BuildRootCmd builds the root command
func BuildRootCmd(conf *config.Config) cobra.Command {
	rootCmd := cobra.Command{
		Use:   "upctl",
		Short: "UpCloud command-line interface",
		Long:  commands.WrapLongDescription(rootCmdLongDescription),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			_, err := valid.ValidateStruct(conf.GlobalFlags)
			if err != nil {
				return err
			}

			// detect desired colour output
			switch {
			case conf.GlobalFlags.ForceColours == config.True:
				text.EnableColors()
			// Environment variable with US spelling to support https://no-color.org
			case conf.GlobalFlags.NoColours == config.True || os.Getenv("NO_COLOR") != "":
				text.DisableColors()
			case conf.GlobalFlags.OutputFormat != config.ValueOutputHuman:
				text.DisableColors()
			default:
				if terminal.IsStdoutTerminal() {
					text.EnableColors()
				} else {
					text.DisableColors()
				}
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

	flags := &pflag.FlagSet{}
	flags.StringVarP(
		&conf.GlobalFlags.ConfigFile, "config", "", "", "Configuration file path.",
	)
	outputFormats := []string{config.ValueOutputHuman, config.ValueOutputJSON, config.ValueOutputYAML}
	flags.StringVarP(
		&conf.GlobalFlags.OutputFormat, "output", "o", "human",
		"Output format. Valid values are "+namedargs.ValidValuesHelp(outputFormats...)+".",
	)
	config.AddToggleFlag(flags, &conf.GlobalFlags.ForceColours, "force-colours", false, "Force coloured output despite detected terminal support.")
	config.AddToggleFlag(flags, &conf.GlobalFlags.NoColours, "no-colours", false, "Disable coloured output despite detected terminal support. Colours can also be disabled by setting NO_COLOR environment variable.")
	flags.BoolVar(
		&conf.GlobalFlags.Debug, "debug", false,
		"Print out more verbose debug logs.",
	)
	flags.DurationVarP(
		&conf.GlobalFlags.ClientTimeout, "client-timeout", "t",
		0,
		"Client timeout to use in API calls.",
	)

	// XXX: Apply viper value to the help as default
	// Add flags
	flags.VisitAll(func(flag *pflag.Flag) {
		rootCmd.PersistentFlags().AddFlag(flag)
	})
	conf.ConfigBindFlagSet(flags)

	rootCmd.SetUsageTemplate(ui.CommandUsageTemplate())
	rootCmd.SetUsageFunc(ui.UsageFunc)

	commands.Must(rootCmd.RegisterFlagCompletionFunc("output", cobra.FixedCompletions(outputFormats, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(rootCmd.RegisterFlagCompletionFunc("client-timeout", cobra.NoFileCompletions))

	return rootCmd
}

// BuildCLI generates the CLI tree and returns the rootCmd
func BuildCLI() cobra.Command {
	conf := config.New()
	rootCmd := BuildRootCmd(conf)

	base.BuildCommands(&rootCmd, conf)

	return rootCmd
}

// Execute is the application entrypoint. It returns the exit code that should be forwarded to the shell.
func Execute() int {
	rootCmd := BuildCLI()
	err := rootCmd.Execute()
	if err != nil {
		if clierr, ok := err.(clierrors.ClientError); ok {
			return clierr.ErrorCode()
		}

		return clierrors.UnspecifiedErrorCode
	}

	return 0
}
