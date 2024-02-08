package gateway

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BasegatewayCommand creates the base "gateway" command
func BasegatewayCommand() commands.Command {
	return &gatewayCommand{
		commands.New("gateway", "Manage gateways"),
	}
}

type gatewayCommand struct {
	*commands.BaseCommand
}
