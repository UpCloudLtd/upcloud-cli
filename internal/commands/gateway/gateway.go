package gateway

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseGatewayCommand creates the base "gateway" command
func BaseGatewayCommand() commands.Command {
	return &gatewayCommand{
		commands.New("gateway", "Manage gateways"),
	}
}

type gatewayCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (gw *gatewayCommand) InitCommand() {
	gw.Cobra().Aliases = []string{"gw"}
}
