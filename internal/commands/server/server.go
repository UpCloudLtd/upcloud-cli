package server

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

const (
	minStorageSize   = 10
	maxServerActions = 10
	// Server state related consts
	defaultStopType             = request.ServerStopTypeSoft
	defaultRestartTimeout       = time.Duration(120) * time.Second
	defaultRestartTimeoutAction = request.RestartTimeoutActionIgnore
	customPlan                  = "custom"
)

// BaseServerCommand crestes the base "server" command
func BaseServerCommand() commands.Command {
	return &serverCommand{
		commands.New("server", "Manage servers"),
	}
}

type serverCommand struct {
	*commands.BaseCommand
}

// waitForServerState waits for server to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func waitForServerState(uuid, state string, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for server %s to be in %s state", uuid, state))

	if _, err := exec.All().WaitForServerState(exec.Context(), &request.WaitForServerStateRequest{
		UUID:         uuid,
		DesiredState: state,
		Timeout:      5 * time.Minute,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, msg)
	exec.PushProgressSuccess(msg)
}
