package all

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseAllCommand creates the base "all" command
func BaseAllCommand() commands.Command {
	return &allCommand{
		commands.New("all", "Manage all UpCloud resources within the account (EXPERIMENTAL)"),
	}
}

type allCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *allCommand) InitCommand() {
	c.Cobra().Hidden = true
	c.Cobra().Long = `Manage all resources within the account (EXPERIMENTAL).

These commands are under development and not all resources types are yet supported.`
}
