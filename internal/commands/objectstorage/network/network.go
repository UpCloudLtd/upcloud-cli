package network

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseNetworkCommand creates the base "object-storage network" command
func BaseNetworkCommand() commands.Command {
	return &networkCommand{
		BaseCommand: commands.New("network", "Manage networks in managed object storage services"),
	}
}

type networkCommand struct {
	*commands.BaseCommand
}
