package accesskey

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseAccessKeyCommand creates the base "object-storage access-key" command
func BaseAccessKeyCommand() commands.Command {
	return &accessKeyCommand{
		BaseCommand: commands.New("access-key", "Manage access keys in managed object storage services"),
	}
}

type accessKeyCommand struct {
	*commands.BaseCommand
}
