package account

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/account/tokenreceiver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
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
}

// DoesNotUseServices implements commands.OfflineCommand as this command does not use services
func (s *loginCommand) DoesNotUseServices() {}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *loginCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.withToken.Value() {
		return s.executeWithToken(exec)
	}

	return s.execute(exec)
}

func (s *loginCommand) execute(exec commands.Executor) (output.Output, error) {
	msg := "Waiting to receive token from browser."
	exec.PushProgressStarted(msg)

	receiver := tokenreceiver.New()
	err := receiver.Start()
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	err = receiver.OpenBrowser()
	if err != nil {
		url := receiver.GetLoginURL()
		exec.PushProgressUpdate(messages.Update{
			Message: "Failed to open browser.",
			Status:  messages.MessageStatusError,
			Details: fmt.Sprintf("Please open a browser and navigate to %s to continue with the login.", url),
		})
	}

	token, err := receiver.Wait(exec.Context())
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressUpdate(messages.Update{
		Key:     msg,
		Message: "Saving created token to the system keyring.",
	})

	err = config.SaveTokenToKeyring(token)
	if err != nil {
		return commands.HandleError(exec, msg, fmt.Errorf("failed to save token to keyring: %w", err))
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

func (s *loginCommand) executeWithToken(exec commands.Executor) (output.Output, error) {
	in := bufio.NewReader(s.Cobra().InOrStdin())
	token, err := in.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read token from standard input: %w", err)
	}

	msg := "Saving provided token to the system keyring."
	exec.PushProgressStarted(msg)
	err = config.SaveTokenToKeyring(strings.TrimSpace(token))
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
