package stack

import "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"

// DeployCommand creates the "stack deploy" command
func DeployCommand() commands.Command {
	return &deployCommand{
		BaseCommand: commands.New(
			"deploy",
			"Deploy a stack (EXPERIMENTAL)",
			"upctl stack deploy <stack-name>",
			"upctl stack deploy (supabase|dokku|starter-kit)",
		),
	}
}

type deployCommand struct {
	*commands.BaseCommand
}
