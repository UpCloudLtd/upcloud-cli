package server

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
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

// waitForServerState waits for server to reach given state and updates given logline with wait progress. Finally, logline is updated with given msg and either done state or timeout warning.
func waitForServerState(uuid, state string, service service.Server, logline *ui.LogEntry, msg string) {
	logline.SetMessage(fmt.Sprintf("Waiting for server %s to be in %s state: polling", uuid, state))

	if _, err := service.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         uuid,
		DesiredState: state,
		Timeout:      5 * time.Minute,
	}); err != nil {
		logline.SetMessage(ui.LiveLogEntryWarningColours.Sprintf("%s: partially done (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")

		return
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()
}
