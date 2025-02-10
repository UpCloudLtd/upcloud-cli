package objectstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseobjectstorageCommand creates the base "object-storage" command
func BaseobjectstorageCommand() commands.Command {
	return &objectstorageCommand{
		commands.New("object-storage", "Manage managed object storage services"),
	}
}

type objectstorageCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *objectstorageCommand) InitCommand() {
	c.Cobra().Aliases = []string{"obs", "objectstorage", "objsto"}

	// Deprecating objectstorage and objsto in favour of object-storage
	// TODO: Remove this in the future
	commands.SetDeprecationHelp(c.Cobra(), []string{"objectstorage", "objsto"})
}
