package server

import (
	"time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

const (
	minStorageSize   = 10
	maxServerActions = 10
	// Server state related consts
	defaultStopType             = request.ServerStopTypeSoft
	defaultRestartTimeout       = time.Duration(120) * time.Second
	defaultRestartTimeoutAction = request.RestartTimeoutActionIgnore
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

// TODO: re-add when wait flag is refactored
/*func serverStateWaiter(uuid, state, msg string, service service.Server, logline *ui.LogEntry) func() error {
	return func() error {
		for {
			time.Sleep(100 * time.Millisecond)
			details, err := service.GetServerDetails(&request.GetServerDetailsRequest{UUID: uuid})
			if err != nil {
				return err
			}
			if details.State == state {
				return nil
			}
			logline.SetMessage(fmt.Sprintf("%s: waiting to start (%v)", msg, details.State))
		}
	}
}
*/
