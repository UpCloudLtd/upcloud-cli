package permissions

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BasePermissionsCommand creates the base 'permissions' command
func BasePermissionsCommand() commands.Command {
	return &permissionsCommand{commands.New("permissions", "Manage permissions")}
}

type permissionsCommand struct {
	*commands.BaseCommand
}
