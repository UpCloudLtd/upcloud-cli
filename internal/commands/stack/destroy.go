package stack

import "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"

// DestroyCommand creates the "stack destroy" command
func DestroyCommand() commands.Command {
	return &destroyCommand{
		BaseCommand: commands.New(
			"destroy",
			"Destroy a stack (EXPERIMENTAL)",
			"upctl stack destroy <stack-name>",
			"upctl stack destroy (supabase|dokku|starterkit)",
		),
	}
}

type destroyCommand struct {
	*commands.BaseCommand
}
