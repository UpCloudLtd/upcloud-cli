package server

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

const (
	minStorageSize   = 10
	maxServerActions = 10
	// Server state related consts
	defaultStopType             = request.ServerStopTypeSoft
	defaultRestartTimeout       = time.Duration(120) * time.Second
	defaultRestartTimeoutAction = request.RestartTimeoutActionIgnore
	customPlan                  = "custom"

	avoidHostDescription    = "Host to avoid when scheduling the server. Use this to make sure VMs do not reside on specific host. Refers to value from `host` attribute. Useful when building HA-environments."
	hostDescription         = "Schedule the server on a specific host. Refers to value from `host` attribute. Only available in private clouds."
	simpleBackupDescription = "Simple backup rule in `HHMM,{daily,dailies,weeklies,monthlies}` format or `no`. For example: `2300,dailies`."
)

var (
	remoteAccessTypes = []string{upcloud.RemoteAccessTypeVNC}
	stopTypes         = []string{request.ServerStopTypeSoft, request.ServerStopTypeHard}
	videoModels       = []string{upcloud.VideoModelVGA, upcloud.VideoModelCirrus}
)

// BaseServerCommand creates the base "server" command
func BaseServerCommand() commands.Command {
	return &serverCommand{
		commands.New("server", "Manage servers"),
	}
}

type serverCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (srv *serverCommand) InitCommand() {
	srv.Cobra().Aliases = []string{"srv"}
}

// waitForServerState waits for server to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func waitForServerState(uuid, state string, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for server %s to be in %s state", uuid, state))

	ctx, cancel := context.WithTimeout(exec.Context(), 15*time.Minute)
	defer cancel()

	if _, err := exec.All().WaitForServerState(ctx, &request.WaitForServerStateRequest{
		UUID:         uuid,
		DesiredState: state,
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
