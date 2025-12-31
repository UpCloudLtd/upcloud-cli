package account

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// LoginCommand creates the "account login" command
func LoginCommand() commands.Command {
	return &loginCommand{
		BaseCommand: commands.New(
			"login",
			"Configure an authentication token to the system keyring (EXPERIMENTAL) ",
			"upctl account login --with-token",
		),
	}
}

type loginCommand struct {
	*commands.BaseCommand

	withToken config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (s *loginCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	config.AddToggleFlag(fs, &s.withToken, "with-token", false, "Read token from standard input.")
	s.AddFlags(fs)

	// Require the with-token flag until we support using browser to authenticate.
	commands.Must(s.Cobra().MarkFlagRequired("with-token"))
}

// DoesNotUseServices implements commands.OfflineCommand as this command does not use services
func (s *loginCommand) DoesNotUseServices() {}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *loginCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.withToken.Value() {
		return s.executeWithToken(exec)
	}

	return output.None{}, nil
}

func (s *loginCommand) executeWithToken(exec commands.Executor) (output.Output, error) {
	in := bufio.NewReader(s.Cobra().InOrStdin())
	token, err := in.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read token from standard input: %w", err)
	}

	token = strings.TrimSpace(token)

	// Check if keyring is disabled via global flag or config
	// Note: We check the flag directly because the config isn't loaded for offline commands
	noKeyring, _ := s.Cobra().Root().PersistentFlags().GetBool("no-keyring")
	// Also check if it was set in config via viper (if available)
	if !noKeyring && s.Cobra().Root().PersistentFlags().Changed("config") {
		// Config file was specified, check if no-keyring is set there
		v := viper.New()
		configFile, _ := s.Cobra().Root().PersistentFlags().GetString("config")
		if configFile != "" {
			v.SetConfigFile(configFile)
			v.SetConfigType("yaml")
			_ = v.ReadInConfig()
			noKeyring = v.GetBool("no-keyring")
		}
	}
	if noKeyring {
		// Get config file path from flag if specified
		configFile, _ := s.Cobra().Root().PersistentFlags().GetString("config")
		msg := "Saving token to configuration file (keyring disabled)."
		exec.PushProgressStarted(msg)
		err = config.SaveTokenToConfigFile(token, configFile)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		exec.PushProgressSuccess(msg)
	} else {
		msg := "Saving provided token to the system keyring."
		exec.PushProgressStarted(msg)
		err = config.SaveTokenToKeyring(token)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		exec.PushProgressSuccess(msg)
	}

	return output.None{}, nil
}
