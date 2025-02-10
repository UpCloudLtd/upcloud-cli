package router

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

const maxRouterActions = 10

// BaseRouterCommand creates the base "router" command
func BaseRouterCommand() commands.Command {
	return &routerCommand{
		commands.New("router", "Manage routers"),
	}
}

type routerCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (rt *routerCommand) InitCommand() {
	rt.Cobra().Aliases = []string{"rt"}
}
