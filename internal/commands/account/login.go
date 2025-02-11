package account

import (
	"fmt"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/account/tokenreceiver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/zalando/go-keyring"
)

// LoginCommand creates the "account login" command
func LoginCommand() commands.Command {
	return &loginCommand{
		BaseCommand: commands.New(
			"login",
			"Configure a authentication token.",
			"upctl account login",
		),
	}
}

type loginCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *loginCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
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

	err = keyring.Set("UpCloud", "", token)
	if err != nil {
		return commands.HandleError(exec, msg, fmt.Errorf("failed to save token to keyring: %w", err))
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
