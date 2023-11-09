package objectstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseobjectstorageCommand creates the base "objectstorage" command
func BaseobjectstorageCommand() commands.Command {
	return &objectstorageCommand{
		commands.New("objectstorage", "Manage Managed object storage services"),
	}
}

type objectstorageCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *objectstorageCommand) InitCommand() {
	c.Cobra().Aliases = []string{"object-storage", "objsto"}
}
