package user

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseUserCommand creates the base "object-storage user" command
func BaseUserCommand() commands.Command {
	return &userCommand{
		BaseCommand: commands.New("user", "Manage users in managed object storage services"),
	}
}

type userCommand struct {
	*commands.BaseCommand
}
