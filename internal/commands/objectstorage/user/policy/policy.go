package policy

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseUserPolicyCommand creates the base "object-storage user policy" command
func BaseUserPolicyCommand() commands.Command {
	return &userPolicyCommand{
		BaseCommand: commands.New("policy", "Manage policies attached to a managed object storage user"),
	}
}

type userPolicyCommand struct {
	*commands.BaseCommand
}
