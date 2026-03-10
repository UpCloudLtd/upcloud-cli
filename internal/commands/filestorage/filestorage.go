package filestorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseFileStorageCommand creates the base "file-storage" command
func BaseFileStorageCommand() commands.Command {
	baseCmd := commands.New("file-storage", "Manage file storage services")

	return &filestorageCommand{
		BaseCommand: baseCmd,
	}
}

type filestorageCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *filestorageCommand) InitCommand() {
	c.Cobra().Aliases = []string{"nfs", "filestorage"}
}
