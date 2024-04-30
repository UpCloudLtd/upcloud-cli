package host

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseHostCommand creates the base "host" command
func BaseHostCommand() commands.Command {
	return &hostCommand{
		commands.New("host", "Manage private cloud hosts"),
	}
}

type hostCommand struct {
	*commands.BaseCommand
}
