package objectstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseobjectstorageCommand creates the base "object-storage" command
func BaseobjectstorageCommand() commands.Command {
	baseCmd := commands.New("object-storage", "Manage managed object storage services")

	return &objectstorageCommand{
		BaseCommand: baseCmd,
	}
}

type objectstorageCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *objectstorageCommand) InitCommand() {
	c.Cobra().Aliases = []string{"obs", "objectstorage", "objsto"}
}
