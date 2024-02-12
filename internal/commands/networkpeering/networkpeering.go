package networkpeering

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseNetworkPeeringCommand creates the base "networkpeering" command
func BaseNetworkPeeringCommand() commands.Command {
	return &networkpeeringCommand{
		commands.New("networkpeering", "Manage network peerings"),
	}
}

type networkpeeringCommand struct {
	*commands.BaseCommand
}
