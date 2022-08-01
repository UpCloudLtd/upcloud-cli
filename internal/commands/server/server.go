package server

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	serverfirewall "github.com/UpCloudLtd/upcloud-cli/internal/commands/server/firewall"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/server/networkinterface"
	serverstorage "github.com/UpCloudLtd/upcloud-cli/internal/commands/server/storage"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
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

func (s *serverCommand) BuildSubCommands(cfg *config.Config) {
	commands.BuildCommand(ListCommand(), s.Cobra(), cfg)
	commands.BuildCommand(PlanListCommand(), s.Cobra(), cfg)
	commands.BuildCommand(ShowCommand(), s.Cobra(), cfg)
	commands.BuildCommand(StartCommand(), s.Cobra(), cfg)
	commands.BuildCommand(RestartCommand(), s.Cobra(), cfg)
	commands.BuildCommand(StopCommand(), s.Cobra(), cfg)
	commands.BuildCommand(CreateCommand(), s.Cobra(), cfg)
	commands.BuildCommand(ModifyCommand(), s.Cobra(), cfg)
	commands.BuildCommand(LoadCommand(), s.Cobra(), cfg)
	commands.BuildCommand(EjectCommand(), s.Cobra(), cfg)
	commands.BuildCommand(DeleteCommand(), s.Cobra(), cfg)

	commands.BuildCommand(networkinterface.BaseNetworkInterfaceCommand(), s.Cobra(), cfg)
	commands.BuildCommand(serverstorage.BaseServerStorageCommand(), s.Cobra(), cfg)
	commands.BuildCommand(serverfirewall.BaseServerFirewallCommand(), s.Cobra(), cfg)
}

// waitForServerState waits for server to reach given state and updates given logline with wait progress. Finally, logline is updated with given msg and either done state or timeout warning.
func waitForServerState(uuid, state string, service service.Server, logline *ui.LogEntry, msg string) {
	logline.SetMessage(fmt.Sprintf("Waiting for server %s to be in %s state: polling", uuid, state))

	if _, err := service.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         uuid,
		DesiredState: state,
		Timeout:      5 * time.Minute,
	}); err != nil {
		logline.SetMessage(fmt.Sprintf("%s: partially done", msg))
		logline.SetDetails(err.Error(), "Error: ")
		logline.MarkWarning()

		return
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()
}
