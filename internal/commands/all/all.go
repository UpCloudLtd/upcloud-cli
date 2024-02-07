package all

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseAllCommand creates the base "all" command
func BaseAllCommand() commands.Command {
	return &allCommand{
		commands.New("all", "Manage all UpCloud resources within the account"),
	}
}

type allCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *allCommand) InitCommand() {}
