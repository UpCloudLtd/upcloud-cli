package router

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
)

const maxRouterActions = 10

// BaseRouterCommand creates the base "router" command
func BaseRouterCommand() commands.Command {
	return &routerCommand{
		commands.New("router", "Manage router"),
	}
}

type routerCommand struct {
	*commands.BaseCommand
}
